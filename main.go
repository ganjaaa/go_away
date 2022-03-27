package main

// https://www.alexedwards.net/blog/working-with-redis
import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func main() {
	// Redis
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	// Webserver
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", auth)
	mux.HandleFunc("/block", blockIP)
	mux.HandleFunc("/unblock", unblockIP)
	mux.HandleFunc("/reload", reload)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func checkIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	}
	return true
}
