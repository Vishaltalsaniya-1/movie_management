package request

type MovieRequest struct {
	Title  string  `json:"title" validate:"required" orm:"column(id)"`
	Genre  string  `json:"genre" validate:"required"`
	Year   int     `json:"year"`
	Rating float64 `json:"rating" validate:"required,gte=1,lte=5"`
}

type Req struct {
	PageNo   int    `json:"page" query:"pageno" validate:"min=1"`
	PageSize int    `json:"page_size" query:"page_size" validate:"min=1"`
	OrderBy  string `json:"sort_by" query:"order_by"`
	Order    string `json:"sort_order" query:"order" validate:"oneof=asc desc"`
	Filter   string `json:"filter" query:"filter"`
	Year     int    `json:"year" query:"year"`
}

// type Auth struct {
// 	Username string `json:"username" validate:"required"`
// 	Email    string `json:"email" validate:"required"`
// 	Password string `json:"password"`
// }

// func init() {
// 	orm.RegisterModel(new(AuthRequest))
// }
