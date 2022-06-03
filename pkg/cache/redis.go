package cache

import (
	"context"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
)

type RedisKeyInformer struct {
	addr string
	conn redis.Conn
	psc  redis.PubSubConn
	ctx  context.Context
	C    <-chan *RedisEvent
	_c   chan *RedisEvent
}

type EventType int

const (
	Set EventType = iota
	Del
)

func (e EventType) String() string {
	switch e {
	case Set:
		return "set"
	case Del:
		return "del"
	}
	return "unknown"
}

type RedisEvent struct {
	Key  string
	Type EventType
}

// watches key events
func NewRedisKeyInformer(ctx context.Context, addr string) (*RedisKeyInformer, error) {
	ri := &RedisKeyInformer{addr: addr, ctx: ctx, _c: make(chan *RedisEvent, 1000)}
	ri.C = ri._c
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	ri.conn = conn

	ri.psc = redis.PubSubConn{Conn: conn}

	if err := ri.psc.PSubscribe(redis.Args{}.AddFlat([]string{"__keyevent*"})...); err != nil {
		return nil, err
	}

	return ri, nil
}

func (r *RedisKeyInformer) Start() {
	go func() {
		defer r.conn.Close()
		cMsg := make(chan redis.Message, 100)
		cErr := make(chan error)
		go func() {
			for {
				select {
				case <-r.ctx.Done():
					return
				case msg := <-cMsg:
					if strings.HasSuffix(msg.Channel, ":json.set") {
						r._c <- &RedisEvent{Key: string(msg.Data), Type: Set}
					} else if strings.HasSuffix(msg.Channel, ":json.del") {
						r._c <- &RedisEvent{Key: string(msg.Data), Type: Del}
					}
				case err := <-cErr:
					fmt.Println(err)
				}
			}
		}()

		for {
			// psc.Receive is a blocking call
			switch n := r.psc.Receive().(type) {
			case error:
				cErr <- n
			case redis.Message:
				cMsg <- n
			default:
				// fmt.Printf("%v\n", n)
			}
		}
	}()
}
