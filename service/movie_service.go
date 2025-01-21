package service

import (
	"fmt"
	"log"
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
	log.Println("timenow--->")
	var existingMovie models.Movie
	if err := db.Where("title = ?", movie.Title).First(&existingMovie).Error; err == nil {
		return nil, fmt.Errorf("a movie with the title '%s' already exists", movie.Title)
	}
	log.Println("existing---->")
	now := time.Now()
	movie.CreatedAt = now
	movie.UpdatedAt = now
	log.Println("creat")
	if err := db.Create(movie).Error; err != nil {
		return nil, fmt.Errorf("failed to create movie: %v", err)
	}

	createdMovie := &response.MovieResponse{
		ID:        uint(movie.ID),
		Title:     movie.Title,
		Genre:     movie.Genre,
		Year:      movie.Year,
		Rating:    movie.Rating,
		CreatedAt: movie.CreatedAt.Format(time.RFC3339),
		UpdatedAt: movie.UpdatedAt.Format(time.RFC3339),
	}

	return createdMovie, nil
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

	updatedMovie := &response.MovieResponse{
		ID:        uint(existingMovie.ID),
		Title:     existingMovie.Title,
		Genre:     existingMovie.Genre,
		Year:      existingMovie.Year,
		Rating:    existingMovie.Rating,
		CreatedAt: existingMovie.CreatedAt.Format(time.RFC3339),
		UpdatedAt: existingMovie.UpdatedAt.Format(time.RFC3339),
	}

	return updatedMovie, nil
}

func DeleteMovie(db *gorm.DB, id int) error {
	if err := db.Delete(&models.Movie{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete movie: %v", err)
	}
	return nil
}

func ListMovies(db *gorm.DB, req request.Req) ([]response.MovieResponse, int, error) {
	
	query := db.Model(&response.MovieResponse{}).Select("id", "title", "genre", "year", "rating", "created_at", "updated_at")

	if req.Filter != "" {
		query = query.Where("title LIKE ? OR genre LIKE ?", "%"+req.Filter+"%", "%"+req.Filter+"%")
	}
	if req.Year != 0 {
		query = query.Where("year = ?", req.Year)
	}
	
	var total int64
if err := query.Count(&total).Error; err != nil {
    return nil, 0, fmt.Errorf("failed to count movies: %v", err)
}

	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "id"
	}
	order := req.Order
	if order == "" {
		order = "asc"
	}
	
	log.Printf("Applying order by: %s %s", orderBy, order)
	query = query.Order(fmt.Sprintf("%s %s", orderBy, order))

	offset := (req.PageNo - 1) * req.PageSize
	log.Printf("Applying pagination with limit: %d and offset: %d", req.PageSize, offset)
	query = query.Offset(offset).Limit(req.PageSize)

	var movies []response.MovieResponse
	if err := query.Find(&movies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch movies: %v", err)
	}

	log.Println("Successfully fetched movies using GORM")
	return movies,int(total), nil
}

func GetMoviesById(db *gorm.DB, id int) (response.MovieResponse, error) {
	var movie response.MovieResponse

	if err := db.Model(&models.Movie{}).Where("id = ?", id).First(&movie).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.MovieResponse{}, fmt.Errorf("movie with id %d not found", id)
		}
		return response.MovieResponse{}, fmt.Errorf("failed to retrieve movie: %v", err)
	}

	return movie, nil
}
func FetchMovieAnalyticsData(db *gorm.DB) (map[string]interface{}, error) {
	if db == nil {
		err := fmt.Errorf("database connection is not initialized")
		log.Println(err)
		return nil, err
	}

	log.Println("Database connection is initialized")

	genreCounts, err := fetchGenreCounts(db)
	if err != nil {
		log.Println("Error fetching genre counts:", err)
		return nil, err
	}

	topRatedData, err := fetchTopRatedMovies(db)
	if err != nil {
		log.Println("Error fetching top-rated movies:", err)
		return nil, err
	}

	recentlyAddedMovies, err := fetchRecentlyAddedMovies(db)
	if err != nil {
		log.Println("Error fetching recently added movies:", err)
		return nil, err
	}

	analytics := map[string]interface{}{
		"genreCounts":         genreCounts,
		"topRatedMovies":      topRatedData,
		"recentlyAddedMovies": recentlyAddedMovies,
	}

	log.Println("Successfully fetched all movie analytics data")
	return analytics, nil
}


func fetchGenreCounts(db *gorm.DB) ([]response.GenreCount, error) {
	log.Println("Fetching genre counts from the database")

	var genreCounts []response.GenreCount
	err := db.Table("movies").
		Select("genre, COUNT(*) AS count").
		Group("genre").
		Scan(&genreCounts).Error

	if err != nil {
		log.Println("Error executing query for genre counts:", err)
		return nil, err
	}

	log.Println("Successfully fetched genre counts")
	return genreCounts, nil
}

func fetchTopRatedMovies(db *gorm.DB) ([]response.MovieResponse, error) {
	log.Println("Fetching top-rated movie data")

	var highestRating float64
	err := db.Table("movies").
		Select("MAX(rating)").
		Row().
		Scan(&highestRating)
	if err != nil {
		log.Println("Error fetching highest rating:", err)
		return nil, err
	}

	var topRatedMovies []response.MovieResponse
	err = db.Table("movies").
		Where("rating = ?", highestRating).
		Find(&topRatedMovies).Error
	if err != nil {
		log.Println("Error fetching top-rated movies:", err)
		return nil, err
	}

	log.Println("Successfully fetched top-rated movie data")
	return topRatedMovies, nil
}

func fetchRecentlyAddedMovies(db *gorm.DB) ([]response.MovieResponse, error) {
	log.Println("Fetching recently added movies")

	var recentlyAddedMovies []response.MovieResponse
	err := db.Table("movies").
		Where("created_at >= NOW() - INTERVAL 1 MINUTE").
		Find(&recentlyAddedMovies).Error

	if err != nil {
		log.Println("Error fetching recently added movies:", err)
		return nil, err
	}

	log.Println("Successfully fetched recently added movie data")
	return recentlyAddedMovies, nil
}
