package redmine

import (
	"fmt"

	rdmn "github.com/nixys/nxs-go-redmine/v4"
	"github.com/nixys/nxs-support-bot/misc"
)

type Project struct {
	ID       int64
	Name     string
	Trackers []misc.IDName
}

// ProjectsGet gets all active project
func (r *Redmine) ProjectsGet() (map[int64]Project, error) {

	projs := make(map[int64]Project)

	rps, _, err := r.c.ProjectAllGet(rdmn.ProjectAllGetRequest{
		Includes: []string{"trackers"},
		Filters: rdmn.ProjectGetRequestFilters{
			Status: rdmn.ProjectStatusActive,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get redmine projects: %w", err)
	}

	for _, p := range rps.Projects {
		projs[int64(p.ID)] = Project{
			ID:   int64(p.ID),
			Name: p.Name,
			Trackers: func() []misc.IDName {
				trackers := []misc.IDName{}
				for _, t := range p.Trackers {
					trackers = append(trackers, misc.IDName{
						ID:   int64(t.ID),
						Name: t.Name,
					})
				}
				return trackers
			}(),
		}
	}

	return projs, nil
}

func (r *Redmine) ProjectGetByIdentifier(identifier string) (Project, error) {

	p, _, err := r.c.ProjectSingleGet(identifier, rdmn.ProjectSingleGetRequest{
		Includes: []string{"trackers"},
	})
	if err != nil {
		return Project{}, fmt.Errorf("redmine get project by identifier: %w", err)
	}

	return Project{
		ID:   int64(p.ID),
		Name: p.Name,
		Trackers: func() []misc.IDName {
			trackers := []misc.IDName{}
			for _, t := range p.Trackers {
				trackers = append(trackers, misc.IDName{
					ID:   int64(t.ID),
					Name: t.Name,
				})
			}
			return trackers
		}(),
	}, nil
}
