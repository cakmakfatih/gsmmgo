package controllers

import (
	"echochat/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	rg *echo.Group
}

func (c *AuthController) register() {
	c.rg.POST("sign-in", signInHandler)
}

func NewAuthController(rg *echo.Group) *AuthController {
	return &AuthController{
		rg: rg,
	}
}

func signInHandler(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := models.LoginUser(&models.UserLoginRequest{Identity: email, Password: password})

	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, user)
}
