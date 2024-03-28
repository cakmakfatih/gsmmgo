package controllers

import (
	"echochat/internals"
	pages "echochat/templates/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IndexController struct {
	rg *echo.Group
}

func (c *IndexController) register() {
	c.rg.GET("", indexHandler)
	c.rg.GET("panel", panelHandler)
}

func NewIndexController(rg *echo.Group) *IndexController {
	return &IndexController{
		rg: rg,
	}
}

func indexHandler(c echo.Context) error {
	return internals.RenderTempl(c, http.StatusOK, pages.Home())
}

func panelHandler(c echo.Context) error {
	return internals.RenderTempl(c, http.StatusOK, pages.Panel())
}
