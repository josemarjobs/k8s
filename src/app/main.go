package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fzzy/radix/redis"
	"github.com/julienschmidt/httprouter"
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r := httprouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		client.Cmd("INCR", "count")
		reply := client.Cmd("GET", "count")
		fmt.Fprintf(w, "I've been hit %s times.\n", reply.String())
	})

	r.GET("/reset", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		reply := client.Cmd("GET", "count")
		client.Cmd("SET", "count", 0)
		fmt.Fprintf(w, "RESET: Last count was %s.\n", reply.String())
	})

	log.Println("Server running on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
