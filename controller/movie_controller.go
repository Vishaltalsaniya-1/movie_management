package controller

import (
	"log"
	"movie_management/managers"
	"movie_management/request"
	"movie_management/response"
	"strconv"
	"net/http"
	"github.com/beego/beego/v2/client/orm"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()


func CreateMovie(c echo.Context) error {
	userEmail, ok := c.Get("user_email").(string)
	if !ok || userEmail == "" {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "Unauthorized"})
	}

	var req request.MovieRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid request"})
	}
	log.Println("req----->")

	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Validation failed"})
	}
	log.Println("req2----->")

	createdMovie, err := managers.CreateMovie(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Error creating movie"})
	}

	return c.JSON(http.StatusCreated, createdMovie)
}

func UpdateMovie(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid ID"})
	}

	var req request.MovieRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid request"})
	}
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Validation failed"})
	}

	updatedMovie, err := managers.UpdateMovie(id, req)
	if err != nil {

		return c.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Movie not found"})

	}

	return c.JSON(http.StatusOK, updatedMovie)
}
func DeleteMovie(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid ID"})
	}
	if err := managers.DeleteMovie(id); err != nil {

		return c.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Movie not found"})
 
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Movie successfully deleted"})
}

func ListMovies(c echo.Context) error {
	var req request.Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid request parameters"})
	}

	if req.PageNo <= 0 {
		req.PageNo = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	if c.QueryParam("per_page") != "" {
		pageSize, err := strconv.Atoi(c.QueryParam("per_page"))
		if err == nil && pageSize > 0 {
			req.PageSize = pageSize
		}
	}

	o := orm.NewOrm()

	movieListResponse, err := managers.ListMovies(o, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to fetch movies"})
	}

	return c.JSON(http.StatusOK, movieListResponse)
}

func GetMoviesById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid movie ID"})
	}

	movie, err := managers.GetMoviesById(id)
	if err != nil {

		return c.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Movie not found"})

	}

	return c.JSON(http.StatusOK, movie)
}

func GetMovieAnalytics(c echo.Context) error {
	analytics, err := managers.GetMovieAnalytics()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to fetch analytics"})
	}

	return c.JSON(http.StatusOK, analytics)
}
