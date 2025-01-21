package db

import (
	"fmt"
	"log"
	"movie_management/models"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}
	if err := DB.AutoMigrate(&models.Movie{}); err != nil {
		return fmt.Errorf("error migrating database: %v", err)
	}

	log.Println("Connected to database!")
	return nil
}

func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func GetDB() (*gorm.DB, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	return DB, nil
}

// func getEnv(key, fallback string) string {
// 	value := os.Getenv(key)
// 	if value == "" {
// 		return fallback
// 	}
// 	return value
// }
