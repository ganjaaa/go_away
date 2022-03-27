package main

import (
	"github.com/gomodule/redigo/redis"
)

func cmdFlush() bool {
	// get pool
	conn := pool.Get()
	defer conn.Close()

	// check exists
	err := redis.Int(conn.Do("FLUSHALL"))

	// exit on error
	if err != nil {
		return false
	}

	return true
}
