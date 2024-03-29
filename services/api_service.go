package services

import (
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
	resp, err := internals.GetWithToken("/api/collections/providers/records", u.Token)

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

		if val, ok := item["alias"].(string); ok {
			provider.Alias = val
		}

		if val, ok := item["method_data"].(string); ok {
			provider.MethodData = val
		}

		providers = append(providers, provider)
	}

	return providers, nil
}

func (s *apiService) GetProvider(u *models.UserModel, id string) error {
	return nil
}

func (s *apiService) CreateProvider(u *models.UserModel, provider models.ProviderModel) error {
	return nil
}

func (s *apiService) UpdateProvider(u *models.UserModel, provider models.ProviderModel) error {
	return nil
}

func (s *apiService) DeleteProvider(u *models.UserModel, id string) error {
	return nil
}
