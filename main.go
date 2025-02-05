package main

import (
	"log"
	"movie_management/controller"
	"movie_management/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/movies", controller.CreateMovie)
	e.PUT("/movies/:id", controller.UpdateMovie)
	e.DELETE("/movies/:id", controller.DeleteMovie)
	e.GET("/movies", controller.ListMovies)
	e.GET("/movies/analytics", controller.GetMovieAnalytics)
	e.GET("/movies/:id", controller.GetMoviesById)

	e.Logger.Fatal(e.Start(":8080"))

	// if managers.MovieProcessingConfig.EnableProducer {
	// 	// Producer mode: Initialize and start producer
	// 	log.Println("Producer is enabled. Initializing producer...")
	// 	rmp := producer.NewProducer()
	// 	producerService := producer.NewProducerService(rmp)
	// 	if err := producerService.Initialize(); err != nil {
	// 		log.Fatal("Failed to initialize producer service:", err)
	// 	}
	// 	// Producer is ready to send tasks
	// 	log.Println("Producer initialized successfully.")
	// } else {
	// 	// Consumer mode: Initialize and start consumer
	// 	log.Println("Producer is disabled. Initializing consumer...")
	// 	consumerInstance := consumer.NewConsumer()
	// 	if err := consumerInstance.Initialize(); err != nil {
	// 		log.Fatal("Failed to initialize consumer service:", err)
	// 	}
	// 	// Consumer is ready to consume tasks
	// 	log.Println("Consumer initialized successfully.")
	// }
}
