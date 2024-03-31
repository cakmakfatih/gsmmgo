package controllers

import (
	"echochat/internals"
	"echochat/middlewares"
	"echochat/models"
	"echochat/services"
	"echochat/utils"
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
	providersRG.DELETE("", cr.deleteProviders, middlewares.AuthGuardMiddleware, middlewares.RedirectIfNotAuthenticated)
	providersRG.POST("/services", cr.createServiceFromProvider, middlewares.AuthGuardMiddleware, middlewares.RedirectIfNotAuthenticated)
	providersRG.GET("/:id", cr.getProvider, middlewares.AuthGuardMiddleware, middlewares.RedirectIfNotAuthenticated)
	providersRG.PATCH("/services/:id", cr.updateServiceFromProvider, middlewares.AuthGuardMiddleware, middlewares.RedirectIfNotAuthenticated)
}

func (cr *PartialsController) deleteProviders(c echo.Context) error {
	return c.NoContent(http.StatusBadRequest)
}

func (cr *PartialsController) updateServiceFromProvider(c echo.Context) error {
	user := c.Get("user").(*models.UserModel)

	id := c.Param("id")
	serviceID := c.FormValue("editServiceID")
	serviceName := c.FormValue("editServiceName")
	refillDuration := c.FormValue("editRefillDuration")
	providerID := c.FormValue("editProviderID")

	if id == "" || serviceID == "" || serviceName == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	serviceIDInt, err := strconv.Atoi(serviceID)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var serviceModel models.ServiceModel

	serviceModel.ID = id
	serviceModel.ServiceID = serviceIDInt
	serviceModel.Name = serviceName
	serviceModel.Provider = providerID

	if refillDuration == "NoRefill" || refillDuration == "0" || refillDuration == "" {
		serviceModel.RefillDuration = 0
		serviceModel.HasRefill = false
	} else {
		refillDurationInt, err := strconv.Atoi(refillDuration)

		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		serviceModel.HasRefill = true
		serviceModel.RefillDuration = refillDurationInt
	}

	_, err = services.ApiService.UpdateService(user, serviceModel)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return internals.RenderTempl(c, http.StatusOK, partials.EditedServiceResponse())
}

func (cr *PartialsController) getProvider(c echo.Context) error {
	user := c.Get("user").(*models.UserModel)

	provider, err := services.ApiService.GetProvider(user, c.Param("id"))

	if err != nil {
		log.Println(err)

		return nil
	}

	services, err := services.ApiService.GetServicesFromProviderID(user, provider.ID)

	if err != nil {
		log.Println(err)

		return nil
	}

	return internals.RenderTempl(c, http.StatusOK, partials.EditProviderResponse(provider, services))
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

	if method == "whatsapp" || method == "manual" {
		addedProvider, err := services.ApiService.CreateProvider(user, providerModel)

		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		return internals.RenderTempl(c, http.StatusCreated, partials.AddedProviderResponse(addedProvider))
	}

	if method == "telegram" {
		telegramChatId := c.FormValue("telegramChatId")

		if telegramChatId == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		providerModel.MethodData = map[string]string{"telegram_chat_id": telegramChatId}
	} else if method == "web" {
		supportUsername := c.FormValue("supportUsername")
		supportPassword := c.FormValue("supportPassword")

		if supportUsername == "" || supportPassword == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		providerModel.MethodData = map[string]string{"support_username": supportUsername, "support_password": supportPassword}
	}

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
	refillDuration := c.FormValue("refillDuration")

	if providerId == "" || serviceId == "" || serviceName == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	id, err := strconv.Atoi(serviceId)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	serviceModel := models.ServiceModel{ServiceID: id, Name: serviceName, Provider: providerId}

	if refillDuration == "" {
		serviceModel.HasRefill = false
	} else {
		refillDurationInt, err := strconv.Atoi(refillDuration)

		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		serviceModel.HasRefill = true
		serviceModel.RefillDuration = refillDurationInt
	}

	addedService, err := services.ApiService.CreateService(user, serviceModel)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if c.FormValue("from") == "edit" {
		return internals.RenderTempl(c, http.StatusCreated, partials.AddServiceFromEditProvider(addedService))
	}

	return internals.RenderTempl(c, http.StatusCreated, partials.AddedServiceFromProviderResponse(addedService))
}
