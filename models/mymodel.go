package models

import "time"

type Movie struct {
    ID        uint      `gorm:"primaryKey"`
    Title     string    `gorm:"type:varchar(255);not null"`
    Genre     string    `gorm:"type:varchar(100);not null"`
    Year      int       `gorm:"not null"`
    Rating    float64   `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}