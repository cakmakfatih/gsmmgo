package controllers

import (
	"echochat/internals"
	"echochat/middlewares"
	"echochat/models"
	"echochat/services"
	"echochat/utils"
	"encoding/json"
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

	panels := c.Get("panels").(*[]models.PanelModel)

	return internals.RenderTempl(c, http.StatusOK, partials.Providers(*panels, tableData))
}

func (cr *PartialsController) createProvider(c echo.Context) error {
	user := c.Get("user").(*models.UserModel)

	panelId := c.FormValue("panel")
	providerUrl := c.FormValue("url")
	method := c.FormValue("method")
	alias := c.FormValue("alias")

	if panelId == "" || providerUrl == "" || method == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	var providerModel models.ProviderModel

	providerModel.Panel = panelId
	providerModel.URL = providerUrl
	providerModel.Method = method
	providerModel.Alias = alias

	if method != "web" {
		err := services.ApiService.CreateProvider(user, providerModel)

		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		return c.NoContent(http.StatusOK)
	}

	supportUsername := c.FormValue("supportUsername")
	supportPassword := c.FormValue("supportPassword")

	if supportUsername == "" || supportPassword == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	methodDataString, err := json.Marshal(map[string]interface{}{"support_username": supportUsername, "support_password": supportPassword})

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	providerModel.MethodData = string(methodDataString)

	err = services.ApiService.CreateProvider(user, providerModel)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}
