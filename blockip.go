package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

func blockIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")

	conn := pool.Get()
	defer conn.Close()

	// Request IP
	paramIP := r.URL.Query().Get("ip")
	if checkIPAddress(paramIP) == false {
		http.Error(w, http.StatusText(403), 403)
		return
	}

	// Check existing of IP
	exists, err := redis.Int(conn.Do("EXISTS", paramIP))
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(403), 403)
		return
	} else if exists == 1 {
		fmt.Fprintf(w, "OK")
		return
	}

	// Set IP
	err = conn.Send("SET", paramIP, 1)
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(403), 403)
		return
	}

	fmt.Fprintf(w, "OK")
}
