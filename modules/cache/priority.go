package cache

import (
	"fmt"

	"github.com/nixys/nxs-support-bot/ds/redis"
	"github.com/nixys/nxs-support-bot/misc"
)

type Priority struct {
	ID        int64
	Name      map[string]string
	IsDefault bool
}

type PriorityLocale struct {
	ID        int64
	Name      string
	IsDefault bool
}

const PriorityLangDefault = "default"

// PriorityGetDefault gets default priority from cache (Redmine)
func (c *Cache) PriorityGetDefault() (Priority, error) {

	priorities, err := c.rds.PrioritiesGet()
	if err != nil {
		return Priority{}, err
	}

	for _, p := range priorities {
		if p.IsDefault == true {
			return Priority(p), nil
		}
	}

	return Priority{}, misc.ErrNotFound
}

func (c *Cache) PriorityGetDefaultLocale(lang string) (PriorityLocale, error) {

	priorities, err := c.rds.PrioritiesGet()
	if err != nil {
		return PriorityLocale{}, err
	}

	for _, p := range priorities {
		if p.IsDefault == true {
			prio := Priority(p)
			return PriorityLocale{
				ID:        prio.ID,
				Name:      prio.NameLocale(lang),
				IsDefault: prio.IsDefault,
			}, nil
		}
	}

	return PriorityLocale{}, misc.ErrNotFound
}

// PriorityGetByID gets priority by specified ID from cache (Redmine)
func (c *Cache) PriorityGetByID(id int64) (Priority, error) {

	priorities, err := c.rds.PrioritiesGet()
	if err != nil {
		return Priority{}, err
	}

	for _, p := range priorities {
		if p.ID == id {
			return Priority(p), nil
		}
	}

	return Priority{}, misc.ErrNotFound
}

func (c *Cache) PriorityGetByIDLocale(id int64, lang string) (PriorityLocale, error) {

	priorities, err := c.rds.PrioritiesGet()
	if err != nil {
		return PriorityLocale{}, err
	}

	for _, p := range priorities {
		if p.ID == id {
			prio := Priority(p)
			return PriorityLocale{
				ID:        prio.ID,
				Name:      prio.NameLocale(lang),
				IsDefault: prio.IsDefault,
			}, nil
		}
	}

	return PriorityLocale{}, misc.ErrNotFound
}

func (c *Cache) PrioritiesGet() ([]Priority, error) {

	var priorities []Priority

	rp, err := c.rds.PrioritiesGet()
	if err != nil {
		return nil, err
	}

	for _, p := range rp {
		priorities = append(priorities, Priority(p))
	}

	return priorities, nil
}

func (c *Cache) PrioritiesGetLocale(lang string) ([]PriorityLocale, error) {

	var pls []PriorityLocale

	rp, err := c.rds.PrioritiesGet()
	if err != nil {
		return nil, err
	}

	for _, p := range rp {
		prio := Priority(p)
		pls = append(
			pls,
			PriorityLocale{
				ID:        prio.ID,
				Name:      prio.NameLocale(lang),
				IsDefault: prio.IsDefault,
			},
		)
	}

	return pls, nil
}

// prioritiesUpdate updates priorities cache with new data
func (c *Cache) prioritiesUpdate() error {

	// Get data from Redmine
	priorities, err := c.r.PrioritiesGet()
	if err != nil {
		return fmt.Errorf("cache priorities update: %w", err)
	}

	// Prepare data for cache
	prios := []redis.Priority{}
	for _, p := range priorities {
		prios = append(prios, redis.Priority(p))
	}

	// Save data into cache
	if err := c.rds.PrioritiesSave(prios); err != nil {
		return fmt.Errorf("cache priorities update: %w", err)
	}

	return nil
}

func (p Priority) NameLocale(lang string) string {
	n, b := p.Name[lang]
	if b == false {
		return PriorityLangDefault
	}
	return n
}
