package response

import "time"

type MovieResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Genre     string    `json:"genre"`
	Year      int       `json:"year"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type ListMoviesResponse struct {
	Movies      []MovieResponse `json:"movies"`
	PageNo      int             `json:"page_no"`
	PageSize    int             `json:"per_page"`
	TotalCount  int             `json:"total_count"`
	LastPage    int             `json:"last_page"`
	CurrentPage int             `json:"current_page"`
}

type AnalyticsResponse struct {
	CountByGenre       map[string]int         `json:"genreCounts,omitempty"`
	TopRatedMoviesData map[string]interface{} `json:"topRated,omitempty"`
	RecentlyAddedCount int                    `json:"recentlyAddedCount,omitempty"`
}

type GenreCount struct {
	Genre string `json:"genre"`
	Count int    `json:"count"`
}

func (MovieResponse) TableName() string {
	return "movies"
}
