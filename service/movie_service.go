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

<<<<<<< HEAD
// var db *sql.DB

func CreateMovie(db *sql.DB, movie *models.Movie) (*response.MovieResponse, error) {
=======
func CreateMovie(db *gorm.DB, movie *models.Movie) (*response.MovieResponse, error) {
>>>>>>> 18ab6fb (useing_gorm)
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

<<<<<<< HEAD
func ListMovies(db *sql.DB, req request.Req) ([]response.MovieResponse, int, error) {
	
	log.Println("service reqlist--->")

	// log.Printf("pageNo: %d, pageSize: %d, OrderBy: %s, Order: %s", pageNo, pageSize, OrderBy, Order)

	query := "SELECT * FROM movies WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM movies WHERE 1=1"
	var args []interface{}

	if req.Filter != "" || req.Year != 0 {
		if req.Filter != "" {
			query += " AND (title LIKE ? OR genre LIKE ?)"
			countQuery += " AND (title LIKE ? OR genre LIKE ?)"
			args = append(args, "%"+req.Filter+"%", "%"+req.Filter+"%")
		}
		if req.Year != 0 {
			query += " AND year = ?"
			countQuery += " AND year = ?"
			args = append(args,req.Year)
		}
	}

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ?", req.OrderBy, req.Order)
	offset := (req.PageNo- 1) * req.PageSize
	args = append(args, req.PageSize, offset)

	log.Printf("Executing query: %s with args: %v\n", query, args)

	var total int
	if err := db.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count movies: %v", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch movies: %v", err)
	}
	defer rows.Close()

	var movies []response.MovieResponse
	for rows.Next() {
		var movie response.MovieResponse
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Rating, &movie.CreatedAt, &movie.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to parse movie row: %v", err)
		}
		movies = append(movies, movie)
	}

	return movies, total, nil
}

func GetMoviesById(db *sql.DB, id int) (response.MovieResponse, error) {
	selectStatement := `
        SELECT * 
        FROM movies 
        WHERE id = ?`
=======
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
>>>>>>> 18ab6fb (useing_gorm)

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
	log.Println("Entering FetchMovieAnalyticsData")

<<<<<<< HEAD
func FetchMovieAnalyticsData(db *sql.DB) (response.AnalyticsResponse, error) {
	log.Println("Entering FetchMovieAnalyticsData")

	if db == nil {
		err := fmt.Errorf("database connection is not initialized")
		log.Println(err)
		return response.AnalyticsResponse{}, err
	}

	log.Println("Database connection is initialized")

	genreCounts, err := fetchGenreCounts(db)
	if err != nil {
		log.Println("Error fetching genre counts:", err)
		return response.AnalyticsResponse{}, err
	}

	topRatedData, err := fetchTopRatedMoviesCount(db)
	if err != nil {
		log.Println("Error fetching top-rated movies:", err)
		return response.AnalyticsResponse{}, err
	}

	recentlyAddedCount, err := fetchRecentlyAddedMoviesCount(db)
	if err != nil {
		log.Println("Error fetching recently added movies count:", err)
		return response.AnalyticsResponse{}, err
	}

	analytics := response.AnalyticsResponse{
		CountByGenre:       genreCounts,
		TopRatedMoviesData: topRatedData,
		RecentlyAddedCount: recentlyAddedCount,
	}

	log.Println("Successfully fetched all movie analytics data")
	return analytics, nil
}



func fetchGenreCounts(db *sql.DB) (map[string]int, error) {
	if db == nil {
		err := fmt.Errorf("database connection is not initialized")
		log.Println(err)
		return nil, err
	}
	log.Println("Fetching genre counts from the database")

	rows, err := db.Query("SELECT genre, COUNT(*) AS count FROM movies GROUP BY genre")
	if err != nil {
		log.Println("Error executing query for genre counts:", err)
		return nil, err
	}
	defer rows.Close()

	genreCounts := make(map[string]int)
	for rows.Next() {
		var genre string
		var count int
		if err := rows.Scan(&genre, &count); err != nil {
			log.Println("Error scanning row for genre counts:", err)
			return nil, err
		}
		genreCounts[genre] = count
	}

	log.Println("Successfully fetched genre counts")
	return genreCounts, nil
}

func fetchTopRatedMoviesCount(db *sql.DB) (map[string]interface{}, error) {
	if db == nil {
		err := fmt.Errorf("database connection is not initialized")
		log.Println(err)
		return nil, err
	}

	var highestRating float64
	var count int

	err := db.QueryRow("SELECT MAX(rating) FROM movies").Scan(&highestRating)
	if err != nil {
		log.Println("Error fetching highest rating:", err)
		return nil, err
	}
	err = db.QueryRow("SELECT COUNT(*) FROM movies WHERE ABS(rating - ?) < 0.001", highestRating).Scan(&count)
	if err != nil {
		log.Println("Error fetching movie count:", err)
		return nil, err
	}

	log.Println("Successfully fetched top-rated movie data")
	return map[string]interface{}{
		"highestRating": highestRating,
		"moviesCount":   count,
	}, nil
}

func fetchRecentlyAddedMoviesCount(db *sql.DB) (int, error) {
	if db == nil {
		err := fmt.Errorf("database connection is not initialized")
		log.Println(err)
		return 0, err
	}

	log.Println("Fetching recently added movie count")
	var count int
	query := "SELECT COUNT(*) FROM movies WHERE created_at >= NOW() - INTERVAL 1 MINUTE"

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Println("Error fetching recently added movie count:", err)
		return 0, err
	}

	log.Println("Successfully fetched recently added movie count")
	return count, nil
=======
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
>>>>>>> 18ab6fb (useing_gorm)
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
