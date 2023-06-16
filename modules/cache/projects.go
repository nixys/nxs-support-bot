package cache

import (
	"fmt"

	"git.nixys.ru/apps/nxs-support-bot/ds/redis"
	"git.nixys.ru/apps/nxs-support-bot/misc"
)

type Project struct {
	ID       int64
	Name     string
	Trackers []misc.IDName
}

func (c *Cache) ProjectsFilterActive(projects []misc.IDName) ([]misc.IDName, error) {

	fps := []misc.IDName{}

	projs, err := c.rds.ProjectsGet()
	if err != nil {
		return nil, fmt.Errorf("cache projects fileter active: %w", err)
	}

	for _, p := range projects {
		if _, b := projs[p.ID]; b == true {
			fps = append(fps, p)
		}
	}

	return fps, nil
}

func (c *Cache) projectsUpdate() error {

	projs := make(map[int64]redis.Project)

	pp, err := c.r.ProjectsGet()
	if err != nil {
		return fmt.Errorf("cache projects update: %w", err)
	}

	for k, v := range pp {
		projs[k] = redis.Project{
			ID:   v.ID,
			Name: v.Name,
			Trackers: func() []misc.IDName {
				trackers := []misc.IDName{}
				for _, t := range v.Trackers {
					trackers = append(trackers, misc.IDName{
						ID:   t.ID,
						Name: t.Name,
					})
				}
				return trackers
			}(),
		}
	}

	if err := c.rds.ProjectsSave(projs); err != nil {
		return fmt.Errorf("cache projects update: %w", err)
	}

	return nil
}
