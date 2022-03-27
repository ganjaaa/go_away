package main

import (
	"github.com/gomodule/redigo/redis"
)

func cmdBlock(ip string) bool {
	// get pool
	conn := pool.Get()
	defer conn.Close()

	// check exists
	exists, err := redis.Int(conn.Do("EXISTS", paramIP))

	// exit on error
	if err != nil {
		return false
	}

	// exit if already exists
	if exists == 1 {
		return false
	}

	// set ip
	err = conn.Send("SET", ip, 1)

	// exit on error
	if err != nil {
		return false
	}

	return true
}
