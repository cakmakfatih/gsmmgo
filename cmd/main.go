package main

import (
	"crypto/tls"
	"echochat/controllers"
	"echochat/internals"
	"echochat/models"
	"echochat/services"
	"encoding/gob"
	"fmt"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
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
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteLaxMode,
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

	if os.Getenv("MODE") == "prod" {
		autoTLSManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("/var/www/.cache"),
		}

		s := http.Server{
			Addr:    ":443",
			Handler: e,
			TLSConfig: &tls.Config{
				GetCertificate: autoTLSManager.GetCertificate,
				NextProtos:     []string{acme.ALPNProto},
			},
		}

		if err := s.ListenAndServeTLS("./cert/cert.pem", "./cert/key.pem"); err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%v", os.Getenv("PORT"))))
	}
}
