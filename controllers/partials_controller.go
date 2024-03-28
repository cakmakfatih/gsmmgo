package controllers

import "github.com/labstack/echo/v4"

type PartialsController struct {
	rg *echo.Group
}

func NewPartialsController(rg *echo.Group) *PartialsController {
	return &PartialsController{
		rg: rg,
	}
}

func (c *PartialsController) register() {
	c.rg.GET("providers", providersHandler)
}

func providersHandler(c echo.Context) error {
	return nil
}
