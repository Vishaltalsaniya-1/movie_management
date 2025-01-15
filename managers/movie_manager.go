package managers

import (
	"database/sql"
	"fmt"
	"log"
	"movie_management/models"
	"movie_management/request"
	"movie_management/response"
	"movie_management/service"
)

func CreateMovie(db *sql.DB, req request.MovieRequest) (*response.MovieResponse, error) {

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

func UpdateMovie(db *sql.DB, id int, req *request.MovieRequest) (response.MovieResponse, error) {
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

	return updatedMovie, nil
}

func DeleteMovie(db *sql.DB, id int) error {

	err := service.DeleteMovie(db, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %v", err)
	}
	return nil
}

func ListMovies(db *sql.DB, filters map[string]interface{}, pageSize, pageNo int, orderBy, order string) ([]response.MovieResponse, int, error) {
	genre := ""
	if val, ok := filters["genre"]; ok {
		genre = val.(string)
	}

	year := 0
	if val, ok := filters["year"]; ok {
		year = val.(int)
	}

	title := ""
	if val, ok := filters["title"]; ok {
		title = val.(string)
	}

	movies, total, err := service.ListMovies(db, genre, title, year, pageSize, pageNo, orderBy, order)
	fmt.Printf("Genre: %s, Title: %s, Year: %d\n", genre, title, year)
fmt.Println("Calling service.ListMovies...")

	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve movies: %v", err)
	}

	return movies, total, nil
}

func GetMovieAnalytics(db *sql.DB) (map[string]interface{}, error) {
	log.Println("analytics_managers------>")
	analytics, err := service.FetchMovieAnalyticsData(db)
	if err != nil {
		return nil, err
	}
	return analytics, nil
}
func GetMoviesById(db *sql.DB, id int) (response.MovieResponse, error) {
	movie, err := service.GetMoviesById(db, id)
	if err != nil {
		return response.MovieResponse{}, err
	}
	return movie, nil
}
