package main

import (
	"echochat/controllers"
	"echochat/internals"
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	internals.InitConfig()

	e := echo.New()

	e.Use(middleware.Gzip())
	e.Use(middleware.CSRF())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET_TOKEN")))))
	e.Static("assets", "./assets")

	indexRG := e.Group("/")
	partialsRG := e.Group("/partials/")

	indexController := controllers.NewIndexController(indexRG)
	partialsController := controllers.NewPartialsController(partialsRG)

	controllers.RegisterControllers([]controllers.Controller{indexController, partialsController})

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%v", os.Getenv("PORT"))))
}
