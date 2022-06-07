package cache

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
)

var pool *redis.Pool

type RedisManager struct {
}

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

func NewRedisManager() *RedisManager {
	return &RedisManager{}
}

func (r *RedisManager) JSONGetBytes(key string) ([]byte, error) {
	con := NewRedisPoolConnection()
	defer con.Close()

	rh := rejson.NewReJSONHandler()
	data, err := redis.Bytes(rh.JSONGet("all-pr-index", "."))
	if err != nil {
		return nil, fmt.Errorf("failed to redis.JSONGet key %s: %s", key, err)
	}
	return data, nil
}
