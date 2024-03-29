package middlewares

import (
	"echochat/models"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AuthGuardMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)

		if err != nil {
			c.Set("is_authenticated", false)

			return next(c)
		}

		user, exists := sess.Values["user"].(*models.UserSession)

		if !exists {
			c.Set("is_authenticated", false)

			return next(c)
		}

		csrf, err := c.Cookie("_csrf")

		if err != nil {
			c.Set("is_authenticated", false)

			return next(c)
		}

		if user.CSRF != csrf.Value {
			c.Set("is_authenticated", false)

			return next(c)
		}

		c.Set("is_authenticated", true)
		c.Set("user", &user.User)

		return next(c)
	}
}
