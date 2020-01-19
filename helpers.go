package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

const base = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type redisDriver struct {
}

type dbGetter interface {
	Get(s string) (string, int, error)
}

func dbGet(db dbGetter, s string) (string, int, error) {
	url, status, err := db.Get(s)
	if err != nil {
		return "", status, err
	}
	return url, status, nil
}

func (r *redisDriver) Get(s string) (string, int, error) {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Int(conn.Do("EXISTS", s))
	if err != nil {
		return "", http.StatusInternalServerError, err
	} else if exists == 0 {
		return "", http.StatusNotFound, err
	}
	url, err := redis.String(conn.Do("GET", s))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return url, 0, nil
}

func strToInt(str string) int {
	res := 0
	for _, r := range str {
		res = (62 * res) + strings.Index(base, string(r))
	}
	if res < 0 {
		return -res
	}
	return res
}
