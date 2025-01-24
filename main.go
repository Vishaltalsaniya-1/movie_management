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
}
