package cache

import (
	"fmt"

	"git.nixys.ru/apps/nxs-support-bot/ds/redis"
	"git.nixys.ru/apps/nxs-support-bot/ds/redmine"
)

// Cache it is a module context structure
type Cache struct {
	rds redis.Redis
	r   redmine.Redmine
}

// Settings contains settings for cache
type Settings struct {
	Redmine   redmine.Redmine
	RedisHost string
}

// Init settings up Cache
func Init(s Settings) (Cache, error) {

	rds, err := redis.Connect(s.RedisHost)
	if err != nil {
		return Cache{}, err
	}

	return Cache{
		rds: rds,
		r:   s.Redmine,
	}, nil
}

// Close closes cache endpoints
func (c *Cache) Close() error {
	return c.rds.Close()
}

// Update updates cache
func (c *Cache) Update() error {

	if err := c.prioritiesUpdate(); err != nil {
		return fmt.Errorf("cache update: %w", err)
	}

	if err := c.usersUpdate(); err != nil {
		return fmt.Errorf("cache update: %w", err)
	}

	if err := c.projectsUpdate(); err != nil {
		return fmt.Errorf("cache update: %w", err)
	}

	return nil
}
