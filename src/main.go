package main

// https://www.alexedwards.net/blog/working-with-redis
import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool
var validCmds := [4]string{"reload", "block", "unblock", "flush"}

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
	mux.HandleFunc("/", auth)
	mux.HandleFunc("/cmd", cmd)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func auth(w http.ResponseWriter, r *http.Request) {
	// Basic Header
	w.Header().Add("Content-Type", "text/plain")

	// get pool
	conn := pool.Get()
	defer conn.Close()

	// request IP
	paramIP := getIP(r)
	if net.ParseIP(paramIP) == nil {
		http.Error(w, http.StatusText(403), 403)
		return
	}

	// check Exists
	exists, err := redis.Int(conn.Do("EXISTS", paramIP))
	if err != nil {
		http.Error(w, http.StatusText(403), 403)
		return
	}

	// if doesnt exists
	if exists == 0 {
		fmt.Fprintf(w, "OK")
		return
	}

	// block
	http.Error(w, http.StatusText(403), 403)
	return
}

func cmd(w http.ResponseWriter, r *http.Request) {
	// Basic Header
	w.Header().Add("Content-Type", "text/plain")

	// Validate command
	paramCmd := r.URL.Query().Get("cmd")
	if itemExists(validCmds, paramCmd) == false {
		fmt.Fprintf(w, "invalid cmd\nsupports only: %v", validCmds)
		return
	}
	// Validate IP if block or unblock
	paramIP := r.URL.Query().Get("ip")
	if paramCmd == validCmds[1] || paramCmd == validCmds[2] {
		if net.ParseIP(paramIP) == nil {
			fmt.Fprintf(w, "invalid ip")
			return
		}
	}

	// switch commands
	switch paramCmd {
	case validCmds[0]:
		cmdReload()
	case validCmds[1]:
		cmdBlock(paramIP)
	case validCmds[2]:
		cmdUnblock(paramIP)
	case validCmds[3]:
		cmdFlush()
	}
}

func getIP(r *http.Request) string {
	xforwarded := r.Header.Get("X-FORWARDED-FOR")
	if xforwarded != "" {
		return xforwarded
	}
	xreal := r.Header.Get("HTTP_X_REAL_IP")
	if xreal != "" {
		return xreal
	}

	return r.RemoteAddr
}

func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)
	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

