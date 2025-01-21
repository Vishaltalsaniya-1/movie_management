package service

import (
	"fmt"
	"movie_management/models"
	"movie_management/request"
	"movie_management/response"
	"time"

	"gorm.io/gorm"
)

func CreateMovie(db *gorm.DB, movie *models.Movie) (*response.MovieResponse, error) {
	currentYear := time.Now().Year()
	if movie.Year < 1900 || movie.Year > currentYear {
		return nil, fmt.Errorf("year should be between 1900 and %d", currentYear)
	}

	var existingMovie models.Movie
	if err := db.Where("title = ?", movie.Title).First(&existingMovie).Error; err == nil {
		return nil, fmt.Errorf("a movie with the title '%s' already exists", movie.Title)
	}

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	if err := db.Create(movie).Error; err != nil {
		return nil, fmt.Errorf("failed to create movie: %v", err)
	}

	return &response.MovieResponse{
		ID:        uint(movie.ID),
		Title:     movie.Title,
		Genre:     movie.Genre,
		Year:      movie.Year,
		Rating:    movie.Rating,
		CreatedAt: movie.CreatedAt,
		UpdatedAt: movie.UpdatedAt,
	}, nil
}

func UpdateMovie(db *gorm.DB, movie *models.Movie, id int) (*response.MovieResponse, error) {
	var existingMovie models.Movie
	if err := db.First(&existingMovie, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find movie with id %d: %v", id, err)
	}

	existingMovie.Title = movie.Title
	existingMovie.Genre = movie.Genre
	existingMovie.Year = movie.Year
	existingMovie.Rating = movie.Rating
	existingMovie.UpdatedAt = time.Now()

	if err := db.Save(&existingMovie).Error; err != nil {
		return nil, fmt.Errorf("failed to update movie: %v", err)
	}

	return &response.MovieResponse{
		ID:        uint(existingMovie.ID),
		Title:     existingMovie.Title,
		Genre:     existingMovie.Genre,
		Year:      existingMovie.Year,
		Rating:    existingMovie.Rating,
		CreatedAt: existingMovie.CreatedAt,
		UpdatedAt: existingMovie.UpdatedAt,
	}, nil
}

func DeleteMovie(db *gorm.DB, id int) error {
	if result := db.Delete(&models.Movie{}, id); result.RowsAffected == 0 {
		return fmt.Errorf("movie with ID %d not found", id)
	}
	return nil
}

func ListMovies(db *gorm.DB, req request.Req) ([]response.MovieResponse, int, error) {
	query := db.Model(&models.Movie{})
	
	if req.Filter != "" {
		query = query.Where("title LIKE ? OR genre LIKE ?", "%"+req.Filter+"%", "%"+req.Filter+"%")
	}
	if req.Year != 0 {
		query = query.Where("year = ?", req.Year)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "id"
	}
	order := req.Order
	if order == "" {
		order = "asc"
	}

	query = query.Order(fmt.Sprintf("%s %s", orderBy, order)).
		Offset((req.PageNo - 1) * req.PageSize).Limit(req.PageSize)

	var movies []response.MovieResponse
	if err := query.Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, int(total), nil
}


func FetchMovieAnalyticsData(db *gorm.DB) (map[string]interface{}, error) {
	genreCounts, err := fetchGenreCounts(db)
	if err != nil {
		return nil, err
	}

	topRatedMovies, err := fetchTopRatedMovies(db)
	if err != nil {
		return nil, err
	}

	recentlyAddedMovies, err := fetchRecentlyAddedMovies(db)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"genreCounts":         genreCounts,
		"topRatedMovies":      topRatedMovies,
		"recentlyAddedMovies": recentlyAddedMovies,
	}, nil
}

func GetMoviesById(db *gorm.DB, id int) (response.MovieResponse, error) {
	var movie response.MovieResponse
	if err := db.First(&movie, id).Error; err != nil {
		return response.MovieResponse{}, err
	}
	return movie, nil
}

func fetchGenreCounts(db *gorm.DB) ([]response.GenreCount, error) {
	var genreCounts []response.GenreCount
	if err := db.Table("movies").Select("genre, COUNT(*) AS count").Group("genre").Scan(&genreCounts).Error; err != nil {
		return nil, err
	}
	return genreCounts, nil
}

func fetchTopRatedMovies(db *gorm.DB) ([]response.MovieResponse, error) {
	var topRatedMovies []response.MovieResponse
	if err := db.Where("rating = (SELECT MAX(rating) FROM movies)").Find(&topRatedMovies).Error; err != nil {
		return nil, err
	}
	return topRatedMovies, nil
}

func fetchRecentlyAddedMovies(db *gorm.DB) ([]response.MovieResponse, error) {
	var recentlyAddedMovies []response.MovieResponse
	if err := db.Where("created_at >= NOW() - INTERVAL 1 MINUTE").Find(&recentlyAddedMovies).Error; err != nil {
		return nil, err
	}
	return recentlyAddedMovies, nil
}
