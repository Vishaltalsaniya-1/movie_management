package managers

import (
	"movie_management/models"
	"movie_management/utils"

	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, email, password string) error {
	o := orm.NewOrm()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := models.AuthRequest{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	_, err = o.Insert(&user)
	return err
}

func AuthenticateUser(email, password string) (string, error) {
	o := orm.NewOrm()
	user := models.AuthRequest{}

	err := o.QueryTable("AuthRequest").Filter("Email", email).One(&user)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	return utils.GenerateJWT(user.Email, user.Username)
}
