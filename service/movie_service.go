package service

import (
	"database/sql"
	"fmt"
	"log"
	"movie_management/models"
	"movie_management/request"
	"movie_management/response"
	"time"
)

// var db *sql.DB

func CreateMovie(db *sql.DB, movie *models.Movie) (*response.MovieResponse, error) {
	currentYear := time.Now().Year()
	if movie.Year < 1900 || movie.Year > currentYear {
		return nil, fmt.Errorf("year should be between 1900 and %d", currentYear)
	}

	var existingMovie models.Movie
	query := "SELECT id FROM movies WHERE title = ?"
	err := db.QueryRow(query, movie.Title).Scan(&existingMovie.ID)
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("a movie with the title '%s' already exists", movie.Title)
	}
	now := time.Now()

	query = "INSERT INTO movies (title, genre, year, rating, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, movie.Title, movie.Genre, movie.Year, movie.Rating, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie: %v", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get the last inserted movie ID: %v", err)
	}

	createdMovie := &response.MovieResponse{
		ID:        uint(lastInsertID),
		Title:     movie.Title,
		Genre:     movie.Genre,
		Year:      movie.Year,
		Rating:    movie.Rating,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}

	return createdMovie, nil
}

func UpdateMovie(db *sql.DB, movie *models.Movie, id int) (response.MovieResponse, error) {
	sqlStatement := `
        UPDATE movies
        SET title = ?, genre = ?, year = ?, rating = ?,updated_at = CURRENT_TIMESTAMP
        WHERE id = ?`

	_, err := db.Exec(sqlStatement, movie.Title, movie.Genre, movie.Year, movie.Rating, id)
	if err != nil {
		return response.MovieResponse{}, fmt.Errorf("failed to update movie: %v", err)
	}

	selectStatement := `
	SELECT id, title, genre, year, rating, created_at, updated_at
	FROM movies
	WHERE id = ?`

	var updatedMovie response.MovieResponse
	err = db.QueryRow(selectStatement, id).Scan(
		&updatedMovie.ID,
		&updatedMovie.Title,
		&updatedMovie.Genre,
		&updatedMovie.Year,
		&updatedMovie.Rating,
		&updatedMovie.CreatedAt,
		&updatedMovie.UpdatedAt,
	)
	if err != nil {
		return response.MovieResponse{}, fmt.Errorf("failed to retrieve updated movie: %v", err)
	}

	return updatedMovie, nil
}

func DeleteMovie(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %v", err)
	}
	return nil
}

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

	var movie response.MovieResponse

	err := db.QueryRow(selectStatement, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Genre,
		&movie.Year,
		&movie.Rating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.MovieResponse{}, fmt.Errorf("movie with id %d not found", id)
		}
		return response.MovieResponse{}, fmt.Errorf("failed to retrieve movie: %v", err)
	}

	return movie, nil
}

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
}

// func GetMovieByID(db *sql.DB, id int) (*models.Movie, error) {
// 	var movie models.Movie
// 	err := db.QueryRow("SELECT id, title, genre, year, rating, created_at, updated_at FROM movies WHERE id = ?", id).
// 		Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Rating, &movie.CreatedAt, &movie.UpdatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("movie not found")
// 		}
// 		return nil, fmt.Errorf("failed to fetch movie by ID: %v", err)
// 	}
// 	return &movie, nil
// }

// func GetByAnalytics(db *sql.DB) (*response.AnalyticsResponse, error) {
// 	genreCountQuery := "SELECT genre, COUNT(*) FROM movies GROUP BY genre"
// 	rows, err := db.Query(genreCountQuery)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to count movies by genre: %v", err)
// 	}
// 	defer rows.Close()

// 	genreCount := make(map[string]int)
// 	for rows.Next() {
// 		var genre string
// 		var count int
// 		if err := rows.Scan(&genre, &count); err != nil {
// 			return nil, fmt.Errorf("failed to scan genre count: %v", err)
// 		}
// 		genreCount[genre] = count
// 	}

// 	topRatedMoviesQuery := "SELECT id, title, genre, year, rating, created_at, updated_at FROM movies ORDER BY rating DESC LIMIT 10"
// 	rows, err = db.Query(topRatedMoviesQuery)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get top-rated movies: %v", err)
// 	}
// 	defer rows.Close()

// 	var topRatedMovies []response.MovieResponse
// 	for rows.Next() {
// 		var movie response.MovieResponse
// 		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Rating, &movie.CreatedAt, &movie.UpdatedAt); err != nil {
// 			return nil, fmt.Errorf("failed to scan movie: %v", err)
// 		}
// 		topRatedMovies = append(topRatedMovies, movie)
// 	}

// 	recentlyAddedMoviesQuery := "SELECT id, title, genre, year, rating, created_at, updated_at FROM movies ORDER BY created_at DESC LIMIT 10"
// 	rows, err = db.Query(recentlyAddedMoviesQuery)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get recently added movies: %v", err)
// 	}
// 	defer rows.Close()

// 	var recentlyAddedMovies []response.MovieResponse
// 	for rows.Next() {
// 		var movie response.MovieResponse
// 		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Rating, &movie.CreatedAt, &movie.UpdatedAt); err != nil {
// 			return nil, fmt.Errorf("failed to scan movie: %v", err)
// 		}
// 		recentlyAddedMovies = append(recentlyAddedMovies, movie)
// 	}

// 	return &response.AnalyticsResponse{
// 		GenreCount:          genreCount,
// 		TopRatedMovies:      topRatedMovies,
// 		RecentlyAddedMovies: recentlyAddedMovies,
// 	}, nil
// }
