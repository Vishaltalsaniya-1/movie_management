package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
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

type AuthRequest struct {
	Username string `json:"username" orm:"unique"`
	Email    string `json:"email" orm:"unique"`
	Password string `json:"password"`
}

func init() {
	orm.RegisterModel(new(AuthRequest))
}

// func (user *AuthRequest) HashPassword(password string) error {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	if err != nil {
// 		return err
// 	}
// 	user.Password = string(bytes)
// 	return nil
// }

// func (user *AuthRequest) CheckPassword(providedPassword string) error {
// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
