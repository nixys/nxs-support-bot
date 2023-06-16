package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func Connect(host string) (Redis, error) {

	var r Redis

	client := redis.NewClient(&redis.Options{
		Addr:         host,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	p := client.Ping()

	if p.Err() != nil {
		return r, p.Err()
	}

	r.client = client

	return r, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}
