package controller

import (
	"database/sql"
	"log"
	"movie_management/managers"
	"movie_management/request"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}

func CreateMovie(c echo.Context) error {
	var movieReq request.MovieRequest
	if err := c.Bind(&movieReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	log.Println("Req---------> done")
	movieResponse, err := managers.CreateMovie(db, movieReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, movieResponse)
}

func UpdateMovie(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid movie ID"})
	}

	var movieReq request.MovieRequest
	if err := c.Bind(&movieReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	movieResponse, err := managers.UpdateMovie(db, id, movieReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, movieResponse)
}

func DeleteMovie(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid movie ID"})
	}

	err = managers.DeleteMovie(db, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Movie deleted successfully"})
}

func ListMovies(c echo.Context) error {
    title := c.QueryParam("title")
    genre := c.QueryParam("genre")
    year, _ := strconv.Atoi(c.QueryParam("year"))
    
	limit, err := strconv.Atoi(c.QueryParam("limit"))
    if err != nil || limit <= 0 {
        limit = 10 
    }
    offset, err := strconv.Atoi(c.QueryParam("offset"))
    if err != nil || offset < 0 {
        offset = 0 
    }

    sort := c.QueryParam("sort")
    log.Println("list_movies------------>")

    movieListResponse, err := managers.ListMovies(db, genre, title, year, limit, offset, sort)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, movieListResponse)
}


func GetAnalytics(c echo.Context) error {

	analyticsResponse, err := managers.GetAnalytics(db)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, analyticsResponse)
}
