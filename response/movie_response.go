package response

type MovieResponse struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Genre     string  `json:"genre"`
	Year      int     `json:"year"`
	Rating    float64 `json:"rating"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"update_at"`
}

// type PaginatedMoviesResponse struct {
// 	Movies      []MovieResponse `json:"movies"`
// 	PerPage     int `json:"per_page"`     
// 	PageNo      int `json:"page_no"`      
// 	LastPage    int `json:"last_page"`   
// 	TotalPages  int `json:"total_pages"`  
// }

// type AnalyticsResponse struct {
// 	CountByGenre   map[string]int  `json:"count_by_genre,omitempty"`
// 	TopRatedMovies []MovieResponse `json:"top_rated_movies,omitempty"`
// 	RecentlyAdded  []MovieResponse `json:"recently_added,omitempty"`
// }

type GenreCount struct {
    Genre string `json:"genre"`
    Count int    `json:"count"`
}