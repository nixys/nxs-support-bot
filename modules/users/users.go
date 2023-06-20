package users

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nixys/nxs-support-bot/ds/primedb"
	"github.com/nixys/nxs-support-bot/ds/redmine"
	"github.com/nixys/nxs-support-bot/misc"
	"github.com/nixys/nxs-support-bot/modules/cache"
)

type Settings struct {
	DB       primedb.DB
	Cache    cache.Cache
	Redmine  redmine.Redmine
	Feedback *FeedbackSettings
}

type FeedbackSettings struct {
	ProjectID int64
	UserID    int64
}

type Users struct {
	d        primedb.DB
	c        cache.Cache
	r        redmine.Redmine
	feedback *feedbackCtx
}

type User struct {
	TgID      int64
	RdmnID    int64 // TODO: ???
	Login     string
	FirstName string
	LastName  string
	Lang      string
	Type      userType
}

type UserUpdateData struct {
	RedmineKey *string
	Lang       *string
}

type feedbackCtx struct {
	projectID int64
	userID    int64
}

type userType int

const (
	UserTypeUnauthorized userType = iota
	UserTypeInternal
	UserTypeFeedback
)

func (t userType) String() string {
	return []string{
		"unauthorized",
		"internal",
		"feedback",
	}[t]
}

func Init(s Settings) Users {
	return Users{
		d: s.DB,
		c: s.Cache,
		r: s.Redmine,
		feedback: func() *feedbackCtx {
			if s.Feedback == nil {
				return nil
			}
			return &feedbackCtx{
				projectID: s.Feedback.ProjectID,
				userID:    s.Feedback.UserID,
			}
		}(),
	}
}

func (usrs *Users) Get(tgID int64) (User, error) {

	u, err := usrs.tgRdmnMap(tgID)
	if err != nil {
		return u, fmt.Errorf("user get: %w", err)
	}

	return u, nil
}

func (usrs *Users) InternalGetByRdmnID(rdmnID int64) (User, error) {

	u, err := usrs.d.UserGetByRdmnID(rdmnID)
	if err != nil {
		return User{}, fmt.Errorf("internal user get by Redmine ID: %w", err)
	}

	c, err := usrs.c.UserGet(rdmnID)
	if err != nil {
		return User{}, fmt.Errorf("internal user get by Redmine ID: %w", err)
	}

	return User{
		TgID:      u.TgID,
		RdmnID:    u.RdmnID,
		Login:     c.Login,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Lang:      u.Lang,
		Type:      UserTypeInternal,
	}, nil
}

func (usrs *Users) MembershipsGet(tID int64) ([]misc.IDName, error) {

	u, err := usrs.tgRdmnMap(tID)
	if err != nil {
		return nil, fmt.Errorf("user memberships get: %w", err)
	}

	if u.Type != UserTypeInternal {
		return nil, nil
	}

	// Get all projects user is a member
	mm, err := usrs.r.UserMembershipsGet(u.RdmnID)
	if err != nil {
		return nil, fmt.Errorf("user memberships get: %w", err)
	}

	// Filter active projects
	fps, err := usrs.c.ProjectsFilterActive(mm)
	if err != nil {
		return nil, fmt.Errorf("user memberships get: %w", err)
	}

	sort.Slice(
		fps,
		func(i, j int) bool {
			return strings.ToLower(fps[i].Name) < strings.ToLower(fps[j].Name)
		},
	)

	return fps, nil
}

func (usrs *Users) UserUpdate(tgID int64, d UserUpdateData) (User, error) {

	if d.RedmineKey != nil {

		// If need to update Redmine API key for user

		rID, err := usrs.r.UserAuthCheck(*d.RedmineKey)
		if err != nil {
			return User{}, fmt.Errorf("update user record: %w", err)
		}

		// Get info for specified user
		u, err := usrs.c.UserGet(rID)
		if err != nil {
			return User{}, fmt.Errorf("update user record: %w", err)
		}

		// Update user data in DB
		du, err := usrs.d.UserUpdate(primedb.UserUpdateData{
			TgID:   tgID,
			RdmnID: &u.ID,
			Lang:   d.Lang,
		})
		if err != nil {
			return User{}, fmt.Errorf("update user record: %w", err)
		}

		return User{
			TgID:      du.TgID,
			RdmnID:    du.RdmnID,
			Login:     u.Login,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Lang:      du.Lang,
			Type:      UserTypeInternal,
		}, nil
	}

	u, err := usrs.tgRdmnMap(tgID)
	if err != nil {
		return User{}, fmt.Errorf("update user record: %w", err)
	}

	du, err := usrs.d.UserUpdate(primedb.UserUpdateData{
		TgID: tgID,
		Lang: d.Lang,
	})
	if err != nil {
		return User{}, fmt.Errorf("update user record: %w", err)
	}

	return User{
		TgID:      du.TgID,
		RdmnID:    u.RdmnID,
		Login:     u.Login,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Lang:      du.Lang,
		Type:      u.Type,
	}, nil
}

// tgRdmnMap gets user data from database and cache by Telegram user ID
func (usrs *Users) tgRdmnMap(tgID int64) (User, error) {

	u, err := usrs.d.UserGet(tgID)
	if err != nil {

		if errors.Is(err, misc.ErrNotFound) == true {
			return User{
				TgID:      tgID,
				RdmnID:    0,
				Login:     "",
				FirstName: "",
				LastName:  "",
				Lang:      "",
				Type:      UserTypeUnauthorized,
			}, nil
		}

		return User{}, err
	}

	if u.RdmnID == 0 {
		return User{
			TgID: tgID,
			RdmnID: func() int64 {
				if usrs.feedback != nil {
					return usrs.feedback.userID
				}
				return 0
			}(),
			Login:     "",
			FirstName: "",
			LastName:  "",
			Lang:      u.Lang,
			Type:      UserTypeFeedback,
		}, nil
	}

	c, err := usrs.c.UserGet(u.RdmnID)
	if err != nil {

		// Check user exist and active in Redmine
		if errors.Is(err, misc.ErrNotFound) == true {

			// Set for user Redmine ID in DB to 0
			u, err := usrs.d.UserUpdate(
				primedb.UserUpdateData{
					TgID: tgID,
					RdmnID: func() *int64 {
						a := int64(0)
						return &a
					}(),
				},
			)
			if err != nil {
				return User{}, err
			}

			// Forbidden

			return User{
				TgID: tgID,
				RdmnID: func() int64 {
					if usrs.feedback != nil {
						return usrs.feedback.userID
					}
					return 0
				}(),
				Login:     "",
				FirstName: "",
				LastName:  "",
				Lang:      u.Lang,
				Type:      UserTypeFeedback,
			}, nil
		}

		return User{}, err
	}

	return User{
		TgID:      tgID,
		RdmnID:    c.ID,
		Login:     c.Login,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Lang:      u.Lang,
		Type:      UserTypeInternal,
	}, nil
}
