package cache

import (
	"fmt"

	"github.com/nixys/nxs-support-bot/ds/redis"
	"github.com/nixys/nxs-support-bot/misc"
)

type Priority struct {
	ID        int64
	Name      string
	IsDefault bool
}

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
