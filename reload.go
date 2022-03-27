package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
)

func reload(w http.ResponseWriter, r *http.Request) {
	// Download File
	fileUrl := "https://raw.githubusercontent.com/firehol/blocklist-ipsets/master/firehol_level2.netset"
	err := DownloadFile("ip.list", fileUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Downloaded: " + fileUrl)

	// Redis
	conn := pool.Get()
	defer conn.Close()
	redis.Int(conn.Do("EXISTS", "1"))

	// Read File
	file, err := os.Open("ip.list")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	err = conn.Send("MULTI")
	if err != nil {
		log.Fatal(err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if checkIPAddress(scanner.Text()) {
			err = conn.Send("SET", scanner.Text(), 1)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprintf(w, "OK")
	return
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
