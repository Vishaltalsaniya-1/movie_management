package db

import (
	"fmt"
	"log"
	"movie_management/config"
	"movie_management/models"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Mysql) (*gorm.DB, error) {
	mysqlDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME,
	)
	db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to MySql: %v", err)
		return nil, err
	}
	if err := db.AutoMigrate(&models.Movie{}); err != nil {
		log.Fatalf("Error during auto-migration: %v", err)
		return nil, err
	}
	log.Println("Database migration completed!")

	DB = db
	return DB, nil
}
