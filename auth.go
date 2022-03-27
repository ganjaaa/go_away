package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

func auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")

	conn := pool.Get()
	defer conn.Close()

	// Request IP
	paramIP := getIP(r)
	if checkIPAddress(paramIP) == false {
		fmt.Fprintf(w, "")
		return
	}

	exists, err := redis.Int(conn.Do("EXISTS", paramIP))
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(403), 403)
		return
	} else if exists == 1 {
		fmt.Fprintf(w, "OK")
		return
	}

	http.Error(w, http.StatusText(403), 403)
	return
}
