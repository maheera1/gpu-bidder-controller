package main

import (
	"log"
	"time"
)

func main() {
	log.Println("bidder2: started")
	for {
		log.Println("bidder2: heartbeat")
		time.Sleep(2 * time.Second)
	}
}
