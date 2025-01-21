package models

import (
	"time"
)

type Movie struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	Genre     string    `json:"genre" gorm:"not null"`
	Year      int       `json:"year" gorm:"not null"`
	Rating    float64   `json:"rating" gorm:"type:decimal(3,1);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
