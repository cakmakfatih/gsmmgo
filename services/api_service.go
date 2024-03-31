package services

import (
	"bytes"
	"echochat/entities"
	"echochat/internals"
	"echochat/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type apiService struct {
	h *http.Client
}

var ApiService apiService

func InitAPI() {
	ApiService = apiService{
		h: &http.Client{},
	}
}

func (s *apiService) GetProviders(u *models.UserModel) ([]models.ProviderModel, error) {
	var providers []models.ProviderModel
	resp, err := internals.GetWithToken("/api/collections/providers/records?sort=-created&max=2000&perPage=2000", u.Token)

	if err != nil {
		return providers, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return providers, errors.New("failed response status")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return providers, err
	}

	var providersResponse entities.PbListResponseEntity

	if err := json.Unmarshal(body, &providersResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return providers, err
	}

	for _, item := range providersResponse.Items {
		provider := models.ProviderModel{
			ID:     item["id"].(string),
			Panel:  item["panel"].(string),
			URL:    item["url"].(string),
			Method: item["method"].(string),
		}

		if val, ok := item["method_data"].(map[string]interface{}); ok {
			provider.SetMethodDataReadableFromJsonString(val)
		}

		if val, ok := item["alias"].(string); ok {
			provider.Alias = val
		}

		providers = append(providers, provider)
	}

	return providers, nil
}

func (s *apiService) GetServicesFromProviderID(u *models.UserModel, providerId string) ([]models.ServiceModel, error) {
	var services []models.ServiceModel
	resp, err := internals.GetWithToken(fmt.Sprintf("/api/collections/services/records?sort=-created&perPage=2000&filter=(provider='%v')", providerId), u.Token)

	if err != nil {
		return services, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return services, errors.New("failed response status")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return services, err
	}

	var servicesResponse entities.PbListResponseEntity

	if err := json.Unmarshal(body, &servicesResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return services, err
	}

	for _, item := range servicesResponse.Items {
		service := models.ServiceModel{
			ID:             item["id"].(string),
			Name:           item["name"].(string),
			ServiceID:      int(item["service_id"].(float64)),
			HasRefill:      item["has_refill"].(bool),
			RefillDuration: int(item["refill_duration"].(float64)),
		}

		services = append(services, service)
	}

	return services, nil
}

func (s *apiService) UpdateService(u *models.UserModel, service models.ServiceModel) (models.ServiceModel, error) {
	var result models.ServiceModel
	serviceJSON, err := json.Marshal(service)

	if err != nil {
		return result, err
	}

	resp, err := internals.NewRequestWithToken(http.MethodPatch, fmt.Sprintf("/api/collections/services/records/%v", service.ID), bytes.NewBuffer(serviceJSON), u.Token)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return result, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing JSON:", err)

		return result, err
	}

	return result, nil
}

func (s *apiService) GetProvider(u *models.UserModel, providerId string) (models.ProviderModel, error) {
	var result models.ProviderModel

	resp, err := internals.GetWithToken(fmt.Sprintf("/api/collections/providers/records/%v", providerId), u.Token)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)

		return result, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing JSON:", err)

		return result, err
	}

	return result, nil
}

func (s *apiService) CreateService(u *models.UserModel, service models.ServiceModel) (models.ServiceModel, error) {
	var result models.ServiceModel
	serviceJSON, err := json.Marshal(service)

	if err != nil {
		return result, err
	}

	resp, err := internals.PostWithToken("/api/collections/services/records", bytes.NewBuffer(serviceJSON), u.Token)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return result, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing JSON:", err)

		return result, err
	}

	return result, nil
}

func (s *apiService) CreateProvider(u *models.UserModel, provider models.ProviderModel) (models.ProviderModel, error) {
	var result models.ProviderModel
	providerJSON, err := json.Marshal(provider)

	if err != nil {
		return result, err
	}

	resp, err := internals.PostWithToken("/api/collections/providers/records", bytes.NewBuffer(providerJSON), u.Token)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return result, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing JSON:", err.Error())

		return result, err
	}

	result.SetMethodDataReadableFromSelf()

	return result, nil
}

func (s *apiService) UpdateProvider(u *models.UserModel, provider models.ProviderModel) error {
	return nil
}

func (s *apiService) DeleteProvider(u *models.UserModel, id string) error {
	return nil
}

func (s *apiService) GetPanels(u *models.UserModel) ([]models.PanelModel, error) {
	var panels []models.PanelModel

	resp, err := internals.GetWithToken("/api/collections/panels/records", u.Token)
	if err != nil {
		return panels, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return panels, errors.New("failed response status")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return panels, err
	}

	var panelsResponse entities.PbListResponseEntity
	if err := json.Unmarshal(body, &panelsResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return panels, err
	}

	for _, item := range panelsResponse.Items {
		panel := models.PanelModel{
			ID:              item["id"].(string),
			LoginURL:        item["login_url"].(string),
			SupportUsername: item["support_username"].(string),
			SupportPassword: item["support_password"].(string),
		}

		if val, ok := item["telegram_token"].(string); ok {
			panel.TelegramToken = val
		}

		if val, ok := item["whatsapp_token"].(string); ok {
			panel.WhatsAppToken = val
		}

		panels = append(panels, panel)
	}

	return panels, nil
}
