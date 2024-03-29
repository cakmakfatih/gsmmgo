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
	"strconv"

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

	providersRG.GET("", cr.getProviders, middlewares.AuthGuardMiddleware, middlewares.RedirectIfNotAuthenticated)
	providersRG.POST("", cr.createProvider, middlewares.AuthGuardMiddleware, middlewares.RedirectIfNotAuthenticated)
	providersRG.POST("/services", cr.createServiceFromProvider, middlewares.AuthGuardMiddleware, middlewares.RedirectIfNotAuthenticated)
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
		addedProvider, err := services.ApiService.CreateProvider(user, providerModel)

		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		return internals.RenderTempl(c, http.StatusCreated, partials.AddedProviderResponse(addedProvider))
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

	addedProvider, err := services.ApiService.CreateProvider(user, providerModel)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return internals.RenderTempl(c, http.StatusCreated, partials.AddedProviderResponse(addedProvider))
}

func (cr *PartialsController) createServiceFromProvider(c echo.Context) error {
	user := c.Get("user").(*models.UserModel)

	providerId := c.FormValue("provider")
	serviceId := c.FormValue("serviceId")
	serviceName := c.FormValue("serviceName")

	if providerId == "" || serviceId == "" || serviceName == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	id, err := strconv.Atoi(serviceId)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	serviceModel := models.ServiceModel{ServiceID: id, Name: serviceName, Provider: providerId}

	addedService, err := services.ApiService.CreateService(user, serviceModel)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return internals.RenderTempl(c, http.StatusCreated, partials.AddedServiceFromProviderResponse(addedService))
}
