package request

type MovieRequest struct {
	Title  string  `json:"title" validate:"required"`
	Genre  string  `json:"genre" validate:"required"`
	Year   int     `json:"year"`
	Rating float64 `json:"rating" validate:"required,gte=1,lte=5"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
}
