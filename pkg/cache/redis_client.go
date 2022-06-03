package cache

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func init() {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost:6379"
	}
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", host) },
	}
}
func NewRedisPoolConnection() redis.Conn {
	return pool.Get()
}
