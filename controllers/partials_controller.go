package controllers

import (
	"echochat/middlewares"
	"echochat/models"
	"echochat/services"
	"log"

	"github.com/labstack/echo/v4"
)

type PartialsController struct {
	rg *echo.Group
}

func NewPartialsController(rg *echo.Group) *PartialsController {
	return &PartialsController{
		rg: rg,
	}
}

func (c *PartialsController) register() {
	c.rg.GET("providers", providersHandler, middlewares.AuthGuardMiddleware)
}

func providersHandler(c echo.Context) error {
	user := c.Get("user").(*models.UserModel)
	providers, err := services.ApiService.GetProviders(user)

	if err != nil {
		log.Println(err)

		return nil
	}

	log.Println(providers)

	return nil
}
