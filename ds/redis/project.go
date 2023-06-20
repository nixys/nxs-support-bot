package redis

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/nixys/nxs-support-bot/misc"
)

const projectsKey = "cache:projects"

type Project struct {
	ID       int64
	Name     string
	Trackers []misc.IDName
}

func (r *Redis) ProjectsSave(projects map[int64]Project) error {

	b, err := json.Marshal(projects)
	if err != nil {
		return fmt.Errorf("save redis projects: %w", err)
	}

	s := r.client.Set(projectsKey, b, 0)
	if s.Err() != nil {
		return fmt.Errorf("save redis projects: %w", s.Err())
	}

	return nil
}

func (r *Redis) ProjectsGet() (map[int64]Project, error) {

	projs := make(map[int64]Project)

	pp := r.client.Get(projectsKey)
	if pp.Err() != nil {
		if pp.Err() == redis.Nil {
			// Empty keys
			return projs, nil
		}
		return projs, fmt.Errorf("get redis projects: %w", pp.Err())
	}

	if err := json.Unmarshal([]byte(pp.Val()), &projs); err != nil {
		return projs, fmt.Errorf("get redis projects: %w", err)
	}

	return projs, nil
}
