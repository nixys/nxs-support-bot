package redis

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

const usersKey = "cache:users"

type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (r *Redis) UsersSave(users map[int64]User) error {

	b, err := json.Marshal(users)
	if err != nil {
		return err
	}

	s := r.client.Set(usersKey, b, 0)
	if s.Err() != nil {
		return s.Err()
	}

	return nil
}

func (r *Redis) UsersGet() (map[int64]User, error) {

	uu := make(map[int64]User)

	p := r.client.Get(usersKey)
	if p.Err() != nil {
		if p.Err() == redis.Nil {
			// Empty keys
			return uu, nil
		}
		return uu, p.Err()
	}

	if err := json.Unmarshal([]byte(p.Val()), &uu); err != nil {
		return uu, err
	}

	return uu, nil
}
