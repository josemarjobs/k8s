package cache

import (
	"log"
	"os"

	"github.com/fzzy/radix/redis"
)

var client *redis.Client

func init() {
	server := os.Getenv("REDIS_SERVER_URL")
	if server == "" {
		server = "localhost"
	}
	log.Println("Connecting to redis on: " + server + ":6379")
	var err error
	client, err = redis.Dial("tcp", server+":6379")
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	log.Println("Closing redis connection.")
	client.Close()
}

func Set(key, value string) error {
	r := client.Cmd("SETEX", key, 10, value)
	return r.Err
}

func Get(key string) (string, error) {
	r := client.Cmd("GET", key)
	if r.Err != nil {
		return "", r.Err
	}
	return r.String(), nil
}
