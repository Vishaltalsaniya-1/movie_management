package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func Connect() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}
	fmt.Println("Connected to database!")

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS movies (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL UNIQUE,
		genre VARCHAR(50) NOT NULL,
		year INT NOT NULL,
		rating FLOAT CHECK (rating BETWEEN 0 AND 5),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}
	fmt.Println("Movies table created (if not already exists)!")
	return nil
}

func GetDB() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	return db, nil
}
