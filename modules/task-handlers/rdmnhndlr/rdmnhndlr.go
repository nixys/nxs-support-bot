package rdmnhndlr

import (
	"errors"
	"fmt"

	"git.nixys.ru/apps/nxs-support-bot/misc"
	tgbot "git.nixys.ru/apps/nxs-support-bot/modules/bot"
	"git.nixys.ru/apps/nxs-support-bot/modules/issues"
	"git.nixys.ru/apps/nxs-support-bot/modules/localization"
	"git.nixys.ru/apps/nxs-support-bot/modules/users"
)

type Settings struct {
	Bot        *tgbot.Bot
	LangBundle localization.Bundle
	Users      users.Users
	Issues     issues.Issues
	Feedback   *FeedbackSettings
}

type FeedbackSettings struct {
	ProjectID int64
	UserID    int64
}

type RdmnHndlr struct {
	b        *tgbot.Bot
	lb       localization.Bundle
	usrs     users.Users
	iss      issues.Issues
	feedback *feedbackCtx
}

type PermissionsData struct {
	ViewCurrentIssue bool
	ViewPrivateNotes bool
}

type feedbackCtx struct {
	projectID int64
	userID    int64
}

type sendMessagesPrepData struct {
	users             []misc.IDName
	author            misc.IDName
	permissions       permissionsCheck
	formatter         func(*RdmnHndlr, string, any) (string, error)
	formatterFeedback func(*RdmnHndlr, string, any) (string, error)
	data              any
	pu                *feedbackUser
}

type permissionsCheck struct {
	members        map[int64]PermissionsData
	isPrivateNotes bool
}

type feedbackUser struct {
	tgID int64
	lang string
}

func Init(s Settings) RdmnHndlr {

	return RdmnHndlr{
		b:    s.Bot,
		lb:   s.LangBundle,
		usrs: s.Users,
		iss:  s.Issues,
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

// sendMessagesPrep prepares a sending message for each user to be sent
// including feedback user if necessary
func (rh *RdmnHndlr) sendMessagesPrep(d sendMessagesPrepData) ([]tgbot.SendRcptData, error) {

	sd := []tgbot.SendRcptData{}

	rcpts := rcptsGet(d.users, d.author, d.permissions)

	for _, r := range rcpts {

		// Get internal user by current Redmine ID
		u, err := rh.usrs.InternalGetByRdmnID(r)
		if err != nil {
			if errors.Is(err, misc.ErrNotFound) == true {
				// User with current Redmine ID not registered in bot
				// or user inactive
				continue
			}
			return nil, fmt.Errorf("send message prep: %w", err)
		}

		// Create message for user
		m := ""
		if d.formatter != nil {
			m, err = d.formatter(rh, u.Lang, d.data)
			if err != nil {
				return nil, fmt.Errorf("send message prep: %w", err)
			}
		}

		sd = append(
			sd,
			tgbot.SendRcptData{
				ChatID:  u.TgID,
				Message: m,
			},
		)
	}

	// Prepare message for feedback user (if necessary)
	if rh.feedback != nil && d.formatterFeedback != nil && d.pu != nil && d.author.ID != rh.feedback.userID {

		// Create message for user
		m, err := d.formatterFeedback(rh, d.pu.lang, d.data)
		if err != nil {
			return nil, fmt.Errorf("send message prep: %w", err)
		}

		sd = append(
			sd,
			tgbot.SendRcptData{
				ChatID:  d.pu.tgID,
				Message: m,
			},
		)
	}

	return sd, nil
}

func (rh *RdmnHndlr) send(issueID int64, sd []tgbot.SendRcptData, attachments []int64) error {

	// Prepare attachments
	atts, err := rh.iss.AttachmentsDownload(attachments)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	// Send messages
	sr, err := rh.b.SendMessage(
		tgbot.SendData{
			Rcpts: sd,
			Files: func() []tgbot.SendFileData {
				files := []tgbot.SendFileData{}
				for _, a := range atts {
					files = append(
						files,
						tgbot.SendFileData{
							Reader:      a.Reader,
							Name:        a.Name,
							Caption:     a.Caption,
							ContentType: a.ContentType,
						},
					)
				}
				return files
			}(),
		},
	)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	for _, s := range sr {
		for _, m := range s.MessageIDs {
			if err := rh.iss.BunchAdd(s.ChatID, m, issueID); err != nil {
				return fmt.Errorf("send: %w", err)
			}
		}
	}

	return nil
}

// feedbackUserGet gets feedback user data for specified feedback issue ID
// Returns nil pointer if feedback issue not exist or
// user has a not feedback type (e.g. authorized and become an internal user)
func (rh *RdmnHndlr) feedbackUserGet(issueID int64) (*feedbackUser, error) {

	if rh.feedback == nil {
		return nil, nil
	}

	// Get feedback users Tg ID for specified Redmine issue ID
	puTgID, err := rh.iss.IssueFeedbackUserGet(issueID)
	if err != nil {
		if errors.Is(err, misc.ErrNotFound) == true {
			return nil, nil
		}
		return nil, fmt.Errorf("feedback user for issue get: %w", err)
	}

	// Get info for specified user
	u, err := rh.usrs.Get(puTgID)
	if err != nil {
		return nil, fmt.Errorf("feedback user for issue get: %w", err)
	}

	// Get only feedback users
	if u.Type == users.UserTypeFeedback {
		return &feedbackUser{
			tgID: u.TgID,
			lang: u.Lang,
		}, nil
	}

	return nil, nil
}

// rcptsGet gets recipients for update
func rcptsGet(accs []misc.IDName, exclude misc.IDName, permissions permissionsCheck) []int64 {

	var rcpts []int64

	for _, a := range accs {

		// Check rcpt ID is real
		if a.ID == 0 {
			continue
		}

		// Check rcpt is not an update author
		if a.ID == exclude.ID {
			continue
		}

		// Check rcpt is a project member
		m, e := permissions.members[a.ID]
		if e == false {
			continue
		}

		// Check rcpt allows to view current issue
		if m.ViewCurrentIssue == false {
			continue
		}

		// Check rcpt allows to view private notes (if necessary)
		if permissions.isPrivateNotes == true && m.ViewPrivateNotes == false {
			continue
		}

		// If all checks are passed

		// Preserve recipients uniqueness
		rcpts = func() []int64 {
			for _, r := range rcpts {
				if r == a.ID {
					return rcpts
				}
			}
			return append(rcpts, a.ID)
		}()
	}

	return rcpts
}
