package main

import (
	"echochat/controllers"
	"echochat/internals"
	"echochat/models"
	"echochat/services"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func staticCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/assets") {
			c.Response().Header().Set("Cache-Control", "public, max-age=86399")
		}

		return next(c)
	}
}

func main() {
	internals.InitConfig()
	services.InitAPI()

	e := echo.New()

	e.Use(staticCache)

	gob.Register(&models.UserSession{})
	gob.Register(&[]models.PanelModel{})

	e.Use(middleware.Gzip())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookieDomain:   os.Getenv("COOKIE_DOMAIN"),
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET_TOKEN")))))
	e.Static("assets", "./assets")

	indexRG := e.Group("/")
	authRG := e.Group("/auth/")
	partialsRG := e.Group("/partials/")

	indexController := controllers.NewIndexController(indexRG)
	authController := controllers.NewAuthController(authRG)
	partialsController := controllers.NewPartialsController(partialsRG)

	controllers.RegisterControllers([]controllers.Controller{indexController, authController, partialsController})

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%v", os.Getenv("PORT"))))
}
