package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type Movie struct {
	Id        int       `orm:"auto"`
	Title     string    `orm:"size(255)"`
	Genre     string    `orm:"size(100)"`
	Year      int       `orm:"size(4)"`
	Rating    float64   `orm:"digits(3);decimals(1)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Movie))
}
