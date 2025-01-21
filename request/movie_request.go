package request

type MovieRequest struct {
	Title  string  `json:"title" validate:"required"`
	Genre  string  `json:"genre" validate:"required"`
	Year   int     `json:"year"`
	Rating float64 `json:"rating" validate:"required,gte=1,lte=5"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
}
<<<<<<< HEAD

=======
>>>>>>> 18ab6fb (useing_gorm)
type Req struct {
	PageNo   int    `json:"page" query:"pageno" validate:"min=1"`
	PageSize int    `json:"page_size" query:"page_size" validate:"min=1"`
	OrderBy  string `json:"sort_by" query:"order_by"`
	Order    string `json:"sort_order" query:"order" validate:"oneof=asc desc"`
	Filter   string `json:"filter" query:"filter"`
	Year	int     `json:"year" query:"year"`
}
