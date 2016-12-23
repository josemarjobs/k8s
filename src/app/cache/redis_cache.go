package cache

import (
	"fmt"
	"log"

	"github.com/fzzy/radix/redis"
)

func main() {
	fmt.Println("Using redis")
	client, err := redis.Dial("tcp", "localhost:6379")
	defer client.Close()
	if err != nil {
		log.Fatal(err)

	}
	r := client.Cmd("SET", "foo", "bar")
	if r.Err != nil {
		log.Fatal(r.Err)
	}

	log.Println("SET r.String():", r.String())
	r = client.Cmd("GET", "foo")
	log.Println("GET Response: ", r.String())

	client.Append("SET", "name", "Peter")
	client.Append("GET", "name")

	r = client.GetReply()
	log.Println(r.String())

	r = client.GetReply()
	log.Println(r.String())

	client.Cmd("SET", "first_name", "Peter")
	client.Cmd("SET", "last_name", "Griffin")
	r = client.Cmd("MGET", "first_name", "last_name")
	list, _ := r.List()
	for _, n := range list {
		log.Println(n)
	}

}
