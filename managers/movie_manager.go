package managers

import (
	"database/sql"
	"fmt"
	"log"
	"movie_management/models"
	"movie_management/request"
	"movie_management/response"
	"movie_management/service"

	"time"
)

func CreateMovie(db *sql.DB, movieRequest request.MovieRequest) (*response.MovieResponse, error) {
	if movieRequest.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if movieRequest.Year < 1900 || movieRequest.Year > time.Now().Year() {
		return nil, fmt.Errorf("year must be between 1900 and the current year")
	}
	log.Println("managers req ------->")
	newMovie, err := service.CreateMovie(db, &movieRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie: %v", err)
	}

	return mapToMovieResponse(newMovie), nil
}

func UpdateMovie(db *sql.DB, id int, movieRequest request.MovieRequest) (*response.MovieResponse, error) {
	if movieRequest.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if movieRequest.Year < 1900 || movieRequest.Year > time.Now().Year() {
		return nil, fmt.Errorf("year must be between 1900 and the current year")
	}

	updatedMovie, err := service.UpdateMovie(db, id, &movieRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to update movie: %v", err)
	}

	return mapToMovieResponse(updatedMovie), nil
}

func DeleteMovie(db *sql.DB, id int) error {

	err := service.DeleteMovie(db, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %v", err)
	}
	return nil
}

func ListMovies(db *sql.DB, genre, title string, year,  limit, offset int, sort string) (*response.MovieListResponse, error) {
    log.Println("list_managers---------->")

    movies, err := service.ListMovies(db, genre, year, title,  limit, offset, sort)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve movies: %v", err)
    }

    total := len(movies)  
    var movieResponses []response.MovieResponse
    for _, movie := range movies {
        movieResponses = append(movieResponses, *mapToMovieResponse(&movie))
    }

    return &response.MovieListResponse{
        Movies: movieResponses,
        Total:  total,
    }, nil
}


func GetAnalytics(db *sql.DB) (*response.AnalyticsResponse, error) {

	analytics, err := service.GetAnalytics(db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve analytics: %v", err)
	}

	return mapToAnalyticsResponse(analytics), nil
}

func mapToMovieResponse(movie *models.Movie) *response.MovieResponse {
	return &response.MovieResponse{
		ID:     movie.ID,
		Title:  movie.Title,
		Genre:  movie.Genre,
		Year:   movie.Year,
		Rating: movie.Rating,
	}
}

func mapToAnalyticsResponse(analytics *models.AnalyticsResponse) *response.AnalyticsResponse {
	return &response.AnalyticsResponse{
		TotalMovies:   analytics.TotalMovies,
		AverageRating: float64(analytics.AverageRating),
	}
}
