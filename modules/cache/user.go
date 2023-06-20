package cache

import (
	"fmt"

	"github.com/nixys/nxs-support-bot/ds/redis"
	"github.com/nixys/nxs-support-bot/misc"
)

// UserGet gets user ID by username from cache (Redmine)
func (c *Cache) UserGet(id int64) (redis.User, error) {

	user, err := c.rds.UsersGet()
	if err != nil {
		return redis.User{}, err
	}

	u, b := user[id]
	if b == false {
		return redis.User{}, misc.ErrNotFound
	}

	return u, nil
}

// usersUpdate updates users cache with new data
func (c *Cache) usersUpdate() error {

	uu := make(map[int64]redis.User)

	a, err := c.r.UsersGet()
	if err != nil {
		return fmt.Errorf("cache users update: %w", err)
	}

	for k, v := range a {
		uu[k] = redis.User{
			ID:        v.ID,
			Login:     v.Login,
			FirstName: v.FirstName,
			LastName:  v.LastName,
		}
	}

	if err := c.rds.UsersSave(uu); err != nil {
		return fmt.Errorf("cache users update: %w", err)
	}

	return nil
}
