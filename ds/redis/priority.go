package redis

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

const prioritiesKey = "cache:priorities"

type Priority struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
}

func (r *Redis) PrioritiesSave(priorities []Priority) error {

	b, err := json.Marshal(priorities)
	if err != nil {
		return err
	}

	s := r.client.Set(prioritiesKey, b, 0)
	if s.Err() != nil {
		return s.Err()
	}

	return nil
}

func (r *Redis) PrioritiesGet() ([]Priority, error) {

	priorities := []Priority{}

	prios := r.client.Get(prioritiesKey)
	if prios.Err() != nil {
		if prios.Err() == redis.Nil {
			// Empty keys
			return priorities, nil
		}
		return nil, prios.Err()
	}

	if err := json.Unmarshal([]byte(prios.Val()), &priorities); err != nil {
		return nil, err
	}

	return priorities, nil
}
