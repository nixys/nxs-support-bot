package redmine

import (
	"fmt"
	"net/http"

	rdmn "github.com/nixys/nxs-go-redmine/v4"
	"github.com/nixys/nxs-support-bot/misc"
)

type User struct {
	ID        int64
	Login     string
	FirstName string
	LastName  string
}

// UserAuthCheck checks Redmine API key correct. Function returns User ID on success.
func (r *Redmine) UserAuthCheck(key string) (int64, error) {

	var c rdmn.Context

	c.SetEndpoint(r.host)
	c.SetAPIKey(key)

	u, s, err := c.UserCurrentGet(rdmn.UserCurrentGetRequest{})
	if err != nil {
		if s == http.StatusUnauthorized {
			return 0, fmt.Errorf("redmine auth check: %w", misc.ErrAPIKey)
		}
		return 0, fmt.Errorf("redmine auth check: %w", err)
	}

	return int64(u.ID), nil
}

// TODO: ???
/*
func (r *Redmine) UserCurrentGet() (User, error) {

	u, c, err := r.c.UserCurrentGet(rdmn.UserCurrentGetRequest{})
	if err != nil {
		if c == http.StatusUnauthorized {
			return User{}, fmt.Errorf("get redmine current user: %w", misc.ErrAPIKey)
		}
		return User{}, fmt.Errorf("get redmine current user: %w", err)
	}

	return User{
		ID:        int64(u.ID),
		Login:     u.Login,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}
*/

// TODO: ???
/*
func (r *Redmine) UserGet(id int64) (User, error) {

	u, _, err := r.c.UserSingleGet(
		int(id),
		rdmn.UserSingleGetRequest{},
	)
	if err != nil {
		return User{}, fmt.Errorf("get redmine user: %w", err)
	}

	return User{
		ID:        int64(u.ID),
		Login:     u.Login,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}
*/

func (r *Redmine) UsersGet() (map[int64]User, error) {

	users := make(map[int64]User)

	u, _, err := r.c.UserAllGet(rdmn.UserAllGetRequest{
		Filters: rdmn.UserGetRequestFilters{
			Status: rdmn.UserStatusActive,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get redmine users: %w", err)
	}

	for _, e := range u.Users {
		users[int64(e.ID)] = User{
			ID:        int64(e.ID),
			Login:     e.Login,
			FirstName: e.FirstName,
			LastName:  e.LastName,
		}
	}

	return users, nil
}

func (r *Redmine) UserMembershipsGet(id int64) ([]misc.IDName, error) {

	memberships := []misc.IDName{}

	u, _, err := r.c.UserSingleGet(
		int(id),
		rdmn.UserSingleGetRequest{
			Includes: []string{"memberships"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get redmine user memberships: %w", err)
	}

	for _, m := range u.Memberships {
		memberships = append(memberships, misc.IDName{
			ID:   int64(m.Project.ID),
			Name: m.Project.Name,
		})
	}

	return memberships, nil
}
