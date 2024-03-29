package controllers

import (
	"echochat/internals"
	"echochat/middlewares"
	"echochat/models"
	"echochat/services"
	"echochat/utils"
	"log"
	"net/http"

	partials "echochat/templates/partials"

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

func (cr *PartialsController) register() {
	providersRG := cr.rg.Group("providers")

	providersRG.GET("", cr.getProviders, middlewares.AuthGuardMiddleware)
	providersRG.POST("", cr.createProvider, middlewares.AuthGuardMiddleware)
}

func (cr *PartialsController) getProviders(c echo.Context) error {
	user := c.Get("user").(*models.UserModel)
	providers, err := services.ApiService.GetProviders(user)

	if err != nil {
		log.Println(err)

		return nil
	}

	tableData := utils.ProvidersToTableData(providers)

	return internals.RenderTempl(c, http.StatusOK, partials.Providers(tableData))
}

func (cr *PartialsController) createProvider(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
