package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"movie_management/managers"
	"movie_management/request"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}

func CreateMovie(c echo.Context) error {
	var req request.MovieRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if req.Title == "" || req.Genre == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title and Genre are required"})
	}
	if req.Year == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Year is required"})
	}
	if req.Rating == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Rating must be between 1 and 5"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	createdMovie, err := managers.CreateMovie(db, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, createdMovie)
}

func UpdateMovie(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid movie ID"})
	}

	db, ok := c.Get("db").(*sql.DB)
	if !ok || db == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection is not available"})
	}

	var req request.MovieRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	updatedMovie, err := managers.UpdateMovie(db, id, &req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Movie not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedMovie)
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
    pageNo, _ := strconv.Atoi(c.QueryParam("page_no"))
    pageSize, _ := strconv.Atoi(c.QueryParam("per_page")) 
    orderBy := c.QueryParam("order_by")
    order := c.QueryParam("order")
    genre := c.QueryParam("genre")
    year := c.QueryParam("year")
    title := c.QueryParam("title")

    if pageNo <= 0 {
        pageNo = 1
    }
    if pageSize <= 0 {
        pageSize = 10 
    }
    if order == "" {
        order = "asc"
    }

    validColumns := map[string]bool{"id": true, "title": true, "genre": true, "year": true, "rating": true}
    if !validColumns[orderBy] {
        orderBy = "id"
    }
    if order != "asc" && order != "desc" {
        order = "asc"
    }

    filters := map[string]interface{}{}
    if genre != "" {
        filters["genre"] = genre
    }
    if year != "" {
        yearInt, err := strconv.Atoi(year)
        if err == nil {
            filters["year"] = yearInt
        }
    }
    if title != "" {
        filters["title"] = title
    }

    movies, total, err := managers.ListMovies(db, filters, pageSize, pageNo, orderBy, order)
    fmt.Printf("Filters: %v, pageNo: %d, pageSize: %d, orderBy: %s, order: %s\n", filters, pageNo, pageSize, orderBy, order)
    fmt.Println("Calling managers.ListMovies...")

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch movies"})
    }

    lastPages := (total + pageSize - 1) / pageSize
    if lastPages == 0 {
        lastPages = 1
    }

    response := map[string]interface{}{
        "movies":       movies,
        "page_no":      pageNo,
        "page_size":    pageSize,
        "last_pages":   lastPages,
        "total_count":  total,
        "current_page": pageNo,
    }

    return c.JSON(http.StatusOK, response)
}



func GetMovieAnalytics(c echo.Context) error {
	db := c.Get("db").(*sql.DB)

	log.Println("Analytics_controller--------->")
	analytics, err := managers.GetMovieAnalytics(db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch movie analytics",
		})
	}

	return c.JSON(http.StatusOK, analytics)
}


func GetMoviesById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}

	movie, err := managers.GetMoviesById(db, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "movie not found"})
	}

	return c.JSON(http.StatusOK, movie)
}
