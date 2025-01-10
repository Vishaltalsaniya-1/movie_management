package request

type MovieRequest struct {
	Title  string  `json:"title" validate:"required"`
	Genre  string  `json:"genre"`
	Year   int     `json:"year" validate:"required"`
	Rating float64 `json:"rating" validate:"min=0,max=5"`
}
