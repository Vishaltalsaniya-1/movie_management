package response

import "time"

type MovieResponse struct {
	ID        int       `json:"id" orm:"column(id)"`
	Title     string    `json:"title" orm:"column(title)"`
	Genre     string    `json:"genre" orm:"column(genre)"`
	Year      int       `json:"year" orm:"column(year)"`
	Rating    float64   `json:"rating" orm:"column(rating)"`
	CreatedAt time.Time `json:"created_at" orm:"column(created_at)"`
	UpdatedAt time.Time `json:"updated_at" orm:"column(updated_at)"`
}


type ErrorResponse struct {
    Message string `json:"message"`
}

type ListMoviesResponse struct {
	Movies      []MovieResponse `json:"movies"`
	PageNo      int             `json:"page_no"`
	PageSize    int             `json:"per_page"`
	TotalCount  int             `json:"total_count"`
	LastPage    int             `json:"last_page"`
	CurrentPage int             `json:"current_page"`
}


type GenreCount struct {
	Genre string `json:"genre"`
	Count int    `json:"count"`
}

func (MovieResponse) TableName() string {
	return "movies"
}
