package redmine

type Priority struct {
	ID        int64
	Name      string
	IsDefault bool
}

func (r *Redmine) PrioritiesGet() ([]Priority, error) {

	prios := []Priority{}

	priorities, _, err := r.c.EnumerationPrioritiesAllGet()
	if err != nil {
		return nil, err
	}

	for _, p := range priorities {
		prios = append(prios, Priority{
			ID:        int64(p.ID),
			Name:      p.Name,
			IsDefault: p.IsDefault,
		})
	}

	return prios, nil
}
