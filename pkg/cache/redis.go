package cache

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type RedisKeyInformer struct {
	addr string
	conn redis.Conn
	psc  redis.PubSubConn
	ctx  context.Context
}

func NewRedisKeyInformer(ctx context.Context, addr string) (*RedisKeyInformer, error) {
	ri := &RedisKeyInformer{addr: addr, ctx: ctx}
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	ri.conn = conn

	ri.psc = redis.PubSubConn{Conn: conn}

	if err := ri.psc.PSubscribe(redis.Args{}.AddFlat([]string{"*"})...); err != nil {
		return nil, err
	}

	go ri.do()

	return ri, nil
}

func (r *RedisKeyInformer) do() {
	cMsg := make(chan redis.Message, 100)
	cErr := make(chan error)
	go func() {
		for {
			select {
			case <-r.ctx.Done():
				r.conn.Close()
				return
			case msg := <-cMsg:
				fmt.Println(msg)
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
}
