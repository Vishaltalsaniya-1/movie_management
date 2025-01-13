package models

import "time"

type Movie struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Genre     string    `json:"genre"`
	Year      int       `json:"year"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
