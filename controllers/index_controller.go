package controllers

import (
	"echochat/internals"
	components "echochat/templates/components"
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
	sizeClasses := []string{
		"w-36", "w-52", "w-28", "",
	}
	headers := []string{
		"ID", "URL", "Method", "Method Data",
	}
	values := [][]string{}

	for i := 0; i < 5; i++ {
		values = append(values, []string{"ubqeu79jpsxa48j", "https://1kview.com", "Telegram", "Empty"})
	}

	return internals.RenderTempl(c, http.StatusOK, pages.Panel(components.Table(sizeClasses, headers, values)))
}
