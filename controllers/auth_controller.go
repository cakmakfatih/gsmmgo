package controllers

import (
	"echochat/models"
	"echochat/services"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	rg *echo.Group
}

func (c *AuthController) register() {
	c.rg.POST("sign-in", signInHandler)
	c.rg.GET("sign-out", signOutHandler)
}

func NewAuthController(rg *echo.Group) *AuthController {
	return &AuthController{
		rg: rg,
	}
}

func signOutHandler(c echo.Context) error {
	sess, err := session.Get("session", c)

	if err != nil {
		log.Println(err)

		return c.JSON(http.StatusInternalServerError, err)
	}

	sess.Options.MaxAge = -1

	if err := sess.Save(c.Request(), c.Response().Writer); err != nil {
		log.Println(err)

		return c.JSON(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("HX-Refresh", "true")

	return nil
}

func signInHandler(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := models.LoginUser(&models.UserLoginRequest{Identity: email, Password: password})

	if err != nil {
		log.Println(err)

		return c.NoContent(http.StatusUnauthorized)
	}

	sess, err := session.Get("session", c)

	if err != nil {
		log.Println(err)

		return c.NoContent(http.StatusUnauthorized)
	}

	panels, err := services.ApiService.GetPanels(&user)

	if err != nil {
		log.Println(err)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["panels"] = &panels
	sess.Values["user"] = &models.UserSession{
		CSRF: c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
		User: user,
	}

	err = sess.Save(c.Request(), c.Response())

	if err != nil {
		log.Println(err)

		return c.JSON(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("HX-Refresh", "true")

	return nil
}
