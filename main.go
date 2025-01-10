package main

import (
	"log"
	"movie_management/config"
	"movie_management/controller"
	"movie_management/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.LoadEnv()
	if err := db.Connect(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("Error getting database instance: %v", err)
	}
	controller.InitDB(database)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/movies", controller.CreateMovie)
	e.PUT("/movies/:id", controller.UpdateMovie)
	e.DELETE("/movies/:id", controller.DeleteMovie)
	e.GET("/movies", controller.ListMovies)
	e.GET("/analytics", controller.GetAnalytics)

	e.Logger.Fatal(e.Start(":8080"))
}
