package models

type Movie struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Genre  string  `json:"genre"`
	Year   int     `json:"year"`
	Rating float64 `json:"rating"`
}

type AnalyticsResponse struct {
	TotalMovies   int     `json:"total_movies"`
	AverageRating float32 `json:"average_rating"`
}
