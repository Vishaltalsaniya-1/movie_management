package response


type MovieResponse struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Genre  string  `json:"genre"`
	Year   int     `json:"year"`
	Rating float64 `json:"rating"`
}

type MovieListResponse struct {
	Movies []MovieResponse `json:"movies"`
	Total  int             `json:"total"`  
}

type AnalyticsResponse struct {
    TotalMovies int     `json:"total_movies"`
    AverageRating float64 `json:"average_rating"`
}
