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

type PaginatedMoviesResponse struct {
	Movies      []MovieResponse `json:"movies"`
	CurrentPage int             `json:"current_page"`
	TotalPages  int             `json:"total_pages"`
	TotalCount  int             `json:"total_count"`
	LastPage    int             `json:"last_page"`
}
type AnalyticsResponse struct {
	CountByGenre   map[string]int  `json:"count_by_genre,omitempty"`
	TopRatedMovies []MovieResponse `json:"top_rated_movies,omitempty"`
	RecentlyAdded  []MovieResponse `json:"recently_added,omitempty"`
}
