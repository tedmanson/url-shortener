package server

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type cacher interface {
	get(key string) (*entity, error)
	set(key string, value *entity) error
	exists(key string) (bool, error)
}

type cache struct {
	pool *redis.Pool
}

func newCache() (*cache, error) {
	var pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}

	var client = pool.Get()
	_, err := client.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("error with redis connection: %s", err)
	}

	return &cache{
		pool: pool,
	}, nil
}

func (c cache) set(key string, value *entity) error {
	conn := c.pool.Get()
	defer func() (string, error) {
		err := conn.Close()
		return "", err
	}()

	p, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, p, "EX", 10)
	if err != nil {
		return err
	}

	return nil
}

func (c cache) get(key string) (*entity, error) {
	conn := c.pool.Get()
	defer func() (string, error) {
		err := conn.Close()
		return "", err
	}()
	v, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	var e entity
	err = json.Unmarshal(v, &e)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (c cache) exists(key string) (bool, error) {
	conn := c.pool.Get()
	defer func() (string, error) {
		err := conn.Close()
		return "", err
	}()

	return redis.Bool(conn.Do("EXISTS", key))
}
