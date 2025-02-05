package main

import (
	"log"
	"movie_management/consumer"
)

func main() {
	consumerInstance := consumer.NewConsumer()

	if err := consumerInstance.Initialize(); err != nil {
		log.Fatal("Failed to initialize consumer:", err)
	}

	// Keep consuming messages
	select {} // Keeps the process running

	// Destination
	

}


