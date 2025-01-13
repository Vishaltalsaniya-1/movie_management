package service

import (
	"database/sql"
	"fmt"
	"movie_management/models"
	"movie_management/response"
	"time"
)

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

func ListMovies(db *sql.DB, genre, title string, year, pageSize, pageNo int, orderBy, order string) ([]response.MovieResponse, int, error) {
	query := "SELECT SQL_CALC_FOUND_ROWS id, title, genre, year, rating, created_at, updated_at FROM movies WHERE 1=1"
	args := []interface{}{}

	if genre != "" {
		query += " AND genre = ?"
		args = append(args, genre)
	}
	if year != 0 {
		query += " AND year = ?"
		args = append(args, year)
	}
	if title != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+title+"%")
	}

	if orderBy == "" {
		orderBy = "id" 
	}
	if order == "" {
		order = "ASC" 
	}
	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)

	offset := (pageNo - 1) * pageSize
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var movies []response.MovieResponse
	for rows.Next() {
		var movie response.MovieResponse
		var title, genre sql.NullString
		if err := rows.Scan(&movie.ID, &title, &genre, &movie.Year, &movie.Rating, &movie.CreatedAt, &movie.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan error: %v", err)
		}

		movie.Title = title.String
		movie.Genre = genre.String

		movies = append(movies, movie)
	}

	var total int
	totalQuery := "SELECT FOUND_ROWS()"
	if err := db.QueryRow(totalQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("total count error: %v", err)
	}

	return movies, total, nil
}

func GetMoviesById(db *sql.DB, id int) (response.MovieResponse, error) {
	selectStatement := `
        SELECT id, title, genre, year, rating, created_at, updated_at 
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

func GetMovieAnalytics(db *sql.DB, topRatedLimit, recentlyAddedLimit int) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	genreCountQuery := "SELECT genre, COUNT(*) FROM movies GROUP BY genre"
	genreRows, err := db.Query(genreCountQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to count movies by genre: %v", err)
	}
	defer genreRows.Close()

	genreCount := make(map[string]int)
	for genreRows.Next() {
		var genre string
		var count int
		if err := genreRows.Scan(&genre, &count); err != nil {
			return nil, fmt.Errorf("failed to scan genre count: %v", err)
		}
		genreCount[genre] = count
	}
	analytics["genreCount"] = genreCount

	topRatedQuery := "SELECT id, title, genre, year, rating, created_at FROM movies ORDER BY rating DESC LIMIT ?"
	topRatedRows, err := db.Query(topRatedQuery, topRatedLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top-rated movies: %v", err)
	}
	defer topRatedRows.Close()

	var topRatedMovies []response.MovieResponse
	for topRatedRows.Next() {
		var movie response.MovieResponse
		if err := topRatedRows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Rating, &movie.CreatedAt, &movie.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan top-rated movie: %v", err)
		}
		topRatedMovies = append(topRatedMovies, movie)
	}
	analytics["topRatedMovies"] = topRatedMovies

	recentlyAddedQuery := "SELECT id, title, genre, year, rating, created_at FROM movies ORDER BY created_at DESC LIMIT ?"
	recentlyAddedRows, err := db.Query(recentlyAddedQuery, recentlyAddedLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recently added movies: %v", err)
	}
	defer recentlyAddedRows.Close()

	var recentlyAddedMovies []response.MovieResponse
	for recentlyAddedRows.Next() {
		var movie response.MovieResponse
		if err := recentlyAddedRows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Rating, &movie.CreatedAt, &movie.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan recently added movie: %v", err)
		}
		recentlyAddedMovies = append(recentlyAddedMovies, movie)
	}
	analytics["recentlyAddedMovies"] = recentlyAddedMovies

	return analytics, nil
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
