package controllers

import (
	"echochat/internals"
	"echochat/middlewares"
	pages "echochat/templates/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IndexController struct {
	rg *echo.Group
}

func (c *IndexController) register() {
	c.rg.GET("", indexHandler, middlewares.AuthGuardMiddleware)
	c.rg.GET("panel", panelHandler, middlewares.AuthGuardMiddleware)
}

func NewIndexController(rg *echo.Group) *IndexController {
	return &IndexController{
		rg: rg,
	}
}

func indexHandler(c echo.Context) error {
	if c.Get("is_authenticated") == true {
		return c.Redirect(http.StatusPermanentRedirect, "/panel")
	}

	return internals.RenderTempl(c, http.StatusOK, pages.Home())
}

func panelHandler(c echo.Context) error {
	if c.Get("is_authenticated") == false {
		return c.Redirect(http.StatusPermanentRedirect, "/")
	}

	return internals.RenderTempl(c, http.StatusOK, pages.Panel())
}
