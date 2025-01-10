package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"movie_management/models"
	"movie_management/request"
	"time"
)

func CreateMovie(db *sql.DB, movie *request.MovieRequest) (*models.Movie, error) {
	currentYear := time.Now().Year()
	if movie.Year < 1900 || movie.Year > currentYear {
		return nil, errors.New("year must be between 1900 and the current year")
	}
	if movie.Rating < 0 || movie.Rating > 5 {
		return nil, errors.New("rating must be between 0 and 5")
	}

	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM movies WHERE title = ?", movie.Title).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check movie existence: %v", err)
	}
	if exists > 0 {
		return nil, errors.New("movie with this title already exists")
	}
	log.Println("service------------>")
	stmt, err := db.Prepare("INSERT INTO movies (title, genre, year, rating) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(movie.Title, movie.Genre, movie.Year, movie.Rating)
	if err != nil {
		return nil, fmt.Errorf("failed to execute insert statement: %v", err)
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve last insert ID: %v", err)
	}

	return &models.Movie{
		ID:     int(lastInsertID),
		Title:  movie.Title,
		Genre:  movie.Genre,
		Year:   movie.Year,
		Rating: movie.Rating,
	}, nil
}

func UpdateMovie(db *sql.DB, id int, movie *request.MovieRequest) (*models.Movie, error) {
	stmt, err := db.Prepare("UPDATE movies SET title = ?, genre = ?, year = ?, rating = ? WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(movie.Title, movie.Genre, movie.Year, movie.Rating, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update statement: %v", err)
	}

	return &models.Movie{
		ID:     id,
		Title:  movie.Title,
		Genre:  movie.Genre,
		Year:   movie.Year,
		Rating: movie.Rating,
	}, nil
}

func DeleteMovie(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %v", err)
	}
	return nil
}

func ListMovies(db *sql.DB, genre string, year int, title string, limit int, offset int, sort string) ([]models.Movie, error) {
    query := "SELECT id, title, genre, year, rating FROM movies WHERE 1=1"
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
    log.Println("list_service---------->")

    switch sort {
    case "rating_asc":
        query += " ORDER BY rating ASC"
    case "rating_desc":
        query += " ORDER BY rating DESC"
    case "year_asc":
        query += " ORDER BY year ASC"
    case "year_desc":
        query += " ORDER BY year DESC"
    case "title_asc":
        query += " ORDER BY title ASC"
    case "title_desc":
        query += " ORDER BY title DESC"
    default:
        query += " ORDER BY id DESC"
    }

    
    query += " LIMIT ? OFFSET ?"
    args = append(args, limit, offset)

    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve movies: %v", err)
    }
    defer rows.Close()

    var movies []models.Movie
    for rows.Next() {
        var movie models.Movie
        if err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Year, &movie.Rating); err != nil {
            return nil, fmt.Errorf("failed to scan movie: %v", err)
        }
        movies = append(movies, movie)
    }

    return movies, nil
}


func GetAnalytics(db *sql.DB) (*models.AnalyticsResponse, error) {
	var totalMovies int
	var averageRating float32

	err := db.QueryRow("SELECT COUNT(*), AVG(rating) FROM movies").Scan(&totalMovies, &averageRating)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate analytics: %v", err)
	}

	return &models.AnalyticsResponse{
		TotalMovies:   totalMovies,
		AverageRating: averageRating,
	}, nil
}
