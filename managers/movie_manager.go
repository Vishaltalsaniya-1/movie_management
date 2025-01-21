package managers

import (
	"fmt"
	"log"
	"movie_management/models"
	"movie_management/request"
	"movie_management/response"
	"movie_management/service"

	"gorm.io/gorm"
)

func CreateMovie(db *gorm.DB, req request.MovieRequest) (*response.MovieResponse, error) {

	movie := models.Movie{
		Title:  req.Title,
		Genre:  req.Genre,
		Year:   req.Year,
		Rating: req.Rating,
			// CreatedAt: req.CreatedAt,
			// UpdatedAt: req.UpdatedAt,
	}

	log.Println("Manager: Creating movie...")

	createdMovie, err := service.CreateMovie(db, &movie)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie: %v", err)
	}

	return createdMovie, nil
}

func UpdateMovie(db *gorm.DB, id int, req *request.MovieRequest) (response.MovieResponse, error) {
	movie := &models.Movie{
		Title:  req.Title,
		Genre:  req.Genre,
		Year:   req.Year,
		Rating: req.Rating,
			//UpdatedAt: &time.Time{},
	}

	updatedMovie, err := service.UpdateMovie(db, movie, id)
	if err != nil {
		return response.MovieResponse{}, fmt.Errorf("failed to update movie: %v", err)
	}
	return *updatedMovie, nil
}

func DeleteMovie(db *gorm.DB, id int) error {

	err := service.DeleteMovie(db, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %v", err)
	}
	return nil
}

func ListMovies(db *gorm.DB, req request.Req) (response.ListMoviesResponse, error) {
	log.Println("managers reqList---------->")

	movies, total, err := service.ListMovies(db, req)
	if err != nil {
		return response.ListMoviesResponse{}, fmt.Errorf("failed to retrieve movies: %v", err)
	}

	lastPage := (total + req.PageSize - 1) / req.PageSize
	if lastPage == 0 {
		lastPage = 1
	}

	response := response.ListMoviesResponse{
		Movies:      movies,
		PageNo:      req.PageNo,
		PageSize:    req.PageSize,
		TotalCount:  total,
		LastPage:   lastPage,
		CurrentPage: req.PageNo,
	}

	return response, nil
}

func GetMovieAnalytics(db *gorm.DB) (map[string]interface{}, error) {
	log.Println("analytics_managers------>")

	analytics, err := service.FetchMovieAnalyticsData(db)
	if err != nil {
		return nil, err
	}

	return analytics, nil
}

func GetMoviesById(db *gorm.DB, id int) (response.MovieResponse, error) {
	movie, err := service.GetMoviesById(db, id)
	if err != nil {
		return response.MovieResponse{}, err
	}
	return movie, nil
}
