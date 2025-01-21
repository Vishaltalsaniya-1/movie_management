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
<<<<<<< HEAD
=======

>>>>>>> 18ab6fb (useing_gorm)
type ListMoviesResponse struct {
	Movies      []MovieResponse `json:"movies"`
	PageNo      int             `json:"page_no"`
	PageSize    int             `json:"page_size"`
	TotalCount  int             `json:"total_count"`
	LastPages   int             `json:"last_pages"`
	CurrentPage int             `json:"current_page"`
}
<<<<<<< HEAD

// type PaginatedMoviesResponse struct {
// 	Movies      []MovieResponse `json:"movies"`
// 	PerPage     int `json:"per_page"`
// 	PageNo      int `json:"page_no"`
// 	LastPage    int `json:"last_page"`
// 	TotalPages  int `json:"total_pages"`
// }

=======
// type PaginatedMoviesResponse struct {
// 	Movies      []MovieResponse `json:"movies"`
// 	CurrentPage int             `json:"current_page"`
// 	TotalPages  int             `json:"total_pages"`
// 	TotalCount  int             `json:"total_count"`
// 	LastPage    int             `json:"last_page"`
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

>>>>>>> 18ab6fb (useing_gorm)
type AnalyticsResponse struct {
	CountByGenre       map[string]int      `json:"genreCounts,omitempty"`
	TopRatedMoviesData map[string]interface{} `json:"topRated,omitempty"` 
	RecentlyAddedCount int                 `json:"recentlyAddedCount,omitempty"` 
<<<<<<< HEAD
}


type GenreCount struct {
	Genre string `json:"genre"`
	Count int    `json:"count"`
=======
>>>>>>> 18ab6fb (useing_gorm)
}

func (MovieResponse) TableName() string {
	return "movies"
}