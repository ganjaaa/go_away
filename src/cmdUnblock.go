package main

import (
	"github.com/gomodule/redigo/redis"
)

func cmdUnblock(ip string) bool {
	// get pool
	conn := pool.Get()
	defer conn.Close()

	// check exists
	exists, err := redis.Int(conn.Do("EXISTS", paramIP))

	// exit on error
	if err != nil {
		return false
	}

	// exit if is dont exits
	if exists == 0 {
		return true
	}

	// remove ip
	err = conn.Send("DEL", ip)

	// exit on error
	if err != nil {
		return false
	}

	return true
}
