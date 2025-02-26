package controller

import (
	"movie_management/managers"
	"movie_management/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// var validater = validator.New()

func Register(c echo.Context) error {
	// var req request.Auth
	req := new(models.AuthRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// if err := validater.Struct(req); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	// }
	err := managers.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User registration failed"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User registered successfully"})
}

func Login(c echo.Context) error {
	req := new(models.AuthRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	token, err := managers.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
