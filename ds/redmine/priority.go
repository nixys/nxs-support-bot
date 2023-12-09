package redmine

import (
	"net/http"
	"net/url"
)

type Priority struct {
	ID        int64             `json:"id"`
	Name      map[string]string `json:"name"`
	IsDefault bool              `json:"is_default"`
}

type issuePrioritiesAllResult struct {
	Priorities []priorityResult `json:"issue_priorities"`
}

type priorityResult struct {
	ID        int64             `json:"id"`
	Name      map[string]string `json:"name"`
	IsDefault bool              `json:"is_default"`
	Active    bool              `json:"active"`
}

// PrioritiesGet gets active issue priorities from
// nxs-chat-redmine plugin (additional API method)
func (r *Redmine) PrioritiesGet() ([]Priority, error) {

	var (
		e  issuePrioritiesAllResult
		rs []Priority
	)

	ur := url.URL{
		Path: "/localizations/issue_priorities.json",
	}

	_, err := r.c.Get(&e, ur, http.StatusOK)
	if err != nil {
		return nil, err
	}

	for _, p := range e.Priorities {
		if p.Active == true {
			rs = append(
				rs,
				Priority{
					ID:        p.ID,
					Name:      p.Name,
					IsDefault: p.IsDefault,
				},
			)
		}
	}

	return rs, nil
}
