package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func Connect() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "movie_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	log.Println("Connected to database!")

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS movies (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL UNIQUE,
		genre VARCHAR(50) NOT NULL,
		year YEAR NOT NULL,
		rating DECIMAL(3,1) CHECK (rating BETWEEN 0 AND 5),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	if _, err = db.Exec(createTableQuery); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	log.Println("Movies table is ready!")
	return nil
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func GetDB() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	return db, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
