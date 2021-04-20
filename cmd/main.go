package main

import (
	"log"

	"github.com/redis-go/redis"
)

func main() {
	log.Println("Work in Progress version")
	log.Fatal(redis.Run(":6380"))
}
