package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
)

func cmdReload() bool {
	// Load File
	file, err := os.Open("blocklist.urls")
	if err != nil {
		return false
	}
	defer file.Close()

	// foreach URL in File
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lockAndLoad(scanner.Text())
	}
	return true
}

func lockAndLoad(fileUrl string) bool {
	// Download File
	err := DownloadFile("tmp.list", fileUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Downloaded: " + fileUrl)

	// Redis
	conn := pool.Get()
	defer conn.Close()
	redis.Int(conn.Do("EXISTS", "1"))

	// Read File
	file, err := os.Open("tmp.list")
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer file.Close()

	err = conn.Send("MULTI")
	if err != nil {
		log.Fatal(err)
		return false
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if net.ParseIP(scanner.Text()) != nil {
			err = conn.Send("SET", scanner.Text(), 1)
			if err != nil {
				log.Fatal(err)
				return false
			}
		}
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Fatal(err)
		return false
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
