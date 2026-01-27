package main

import (
	"log"
	"time"
)

func main() {
	log.Println("bidder1: started")
	for {
		log.Println("bidder1: heartbeat")
		time.Sleep(2 * time.Second)
	}
}
