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
	mysqlCfg, err := config.Mysqlconfig()
	if err != nil {
		log.Fatalf("Failed to load MySQL configuration: %v", err)
	}

	database, err := db.Connect(mysqlCfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	controller.InitDB(database)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", database)
			return next(c)
		}
	})

	e.POST("/movies", controller.CreateMovie)
	e.PUT("/movies/:id", controller.UpdateMovie)
	e.DELETE("/movies/:id", controller.DeleteMovie)
	e.GET("/movies", controller.ListMovies)
	e.GET("/movies/analytics", controller.GetMovieAnalytics)
	e.GET("/movies/:id", controller.GetMoviesById)

	e.Logger.Fatal(e.Start(":8080"))
}
