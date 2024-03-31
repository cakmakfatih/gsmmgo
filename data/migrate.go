package data

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
	"os"
)

type PerfectPanelServiceResponse struct {
	Categories []PerfectPanelServiceCategory `json:"data"`
}

type PerfectPanelServiceCategory struct {
	ID         int                   `json:"id"`
	Title      string                `json:"title"`
	Visibility bool                  `json:"visibility"`
	Services   []PerfectPanelService `json:"services"`
	Icon       struct {
		IconType string `json:"icon_type"`
		Icon     string `json:"icon"`
	} `json:"icon"`
}

type PerfectPanelProviderResponse struct {
	Providers []PerfectPanelProvider `json:"data"`
}

type PerfectPanelProvider struct {
	ID                    int         `json:"id"`
	Name                  string      `json:"name"`
	Type                  int         `json:"type"`
	Refill                bool        `json:"refill"`
	Cancel                bool        `json:"cancel"`
	Currency              string      `json:"currency"`
	ProviderRate          bool        `json:"provider_rate"`
	ServiceAutoMin        bool        `json:"-"`
	ServiceAutoMax        bool        `json:"-"`
	ServiceAutoIncrement  bool        `json:"-"`
	ServiceAutoRate       bool        `json:"-"`
	ServiceAutoUpdateRate bool        `json:"-"`
	ServiceAutoStatus     bool        `json:"-"`
	Import                bool        `json:"import"`
	ServiceDescription    bool        `json:"service_description"`
	Form                  bool        `json:"form"`
	DefaultServiceType    interface{} `json:"default_service_type"`
	Subscriptions         bool        `json:"subscriptions"`
	ShowProviderServiceID bool        `json:"show_provider_service_id"`
	Position              int         `json:"position"`
}

type PerfectPanelService struct {
	ID               int           `json:"id"`
	Name             string        `json:"name"`
	Mode             int           `json:"mode"`
	Type             int           `json:"type"`
	FakeSubscription bool          `json:"fake_subscription"`
	ProviderID       int           `json:"provider_id"`
	ProviderError    bool          `json:"provider_error"`
	ProviderType     int           `json:"provider_type"`
	Options          []interface{} `json:"options"`
	Rate             struct {
		Custom      string      `json:"custom"`
		CustomFloat float64     `json:"custom_float"`
		Provider    string      `json:"provider"`
		Converted   interface{} `json:"converted"`
	} `json:"rate"`
	Min struct {
		Custom   int `json:"custom"`
		Provider int `json:"provider"`
	} `json:"min"`
	Max struct {
		Custom   int `json:"custom"`
		Provider int `json:"provider"`
	} `json:"max"`
	Status                        int         `json:"status"`
	AutoRate                      bool        `json:"auto_rate"`
	AutoStatus                    bool        `json:"auto_status"`
	AutoIncrement                 bool        `json:"auto_increment"`
	AutoRateFixed                 string      `json:"auto_rate_fixed"`
	AutoRatePercent               string      `json:"auto_rate_percent"`
	AutoMin                       bool        `json:"auto_min"`
	AutoMax                       bool        `json:"auto_max"`
	ProviderServiceID             string      `json:"provider_service_id"`
	CanEditMassRates              bool        `json:"can_edit_mass_rates"`
	ProviderCurrency              string      `json:"provider_currency"`
	ProviderShowProviderServiceID bool        `json:"provider_show_provider_service_id"`
	SyncError                     bool        `json:"sync_error"`
	LastSuccessSync               interface{} `json:"last_success_sync"`
}

func getServices() []PerfectPanelService {
	var result []PerfectPanelService
	jsonFile, err := os.Open("./data/services.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var perfectPanelServiceCategoryResponse PerfectPanelServiceResponse
	err = json.Unmarshal([]byte(byteValue), &perfectPanelServiceCategoryResponse)

	if err != nil {
		fmt.Println(err)
	}

	for _, v := range perfectPanelServiceCategoryResponse.Categories {
		result = append(result, v.Services...)
	}

	return result
}

func getProviders() []PerfectPanelProvider {
	jsonFile, err := os.Open("./data/providers.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var perfectPanelProviderResponse PerfectPanelProviderResponse
	err = json.Unmarshal([]byte(byteValue), &perfectPanelProviderResponse)

	if err != nil {
		fmt.Println(err)
	}

	return perfectPanelProviderResponse.Providers
}

func getProviderFromService(providerID int, providers *[]PerfectPanelProvider) *PerfectPanelProvider {
	var provider PerfectPanelProvider

	for _, v := range *providers {
		if v.ID == providerID {
			provider = v

			return &provider
		}
	}

	return &provider
}

func addServiceToDb(serviceModel *models.ServiceModel) (models.ServiceModel, error) {
	var result models.ServiceModel
	serviceJSON, err := json.Marshal(*serviceModel)

	if err != nil {
		return result, err
	}

	resp, err := internals.C.PbAdmin.Post("/api/collections/services/records", bytes.NewBuffer(serviceJSON))
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

func pushPerfectPanelServiceToDatabase(service *PerfectPanelService, provider *PerfectPanelProvider) error {
	resp, err := internals.C.PbAdmin.Get(fmt.Sprintf("/api/collections/providers/records?filter=(url='%v')", provider.Name))

	if err != nil {
		fmt.Println(err)

		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed response status")
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return errors.New("couldn't read resp body")
	}

	var providersResponse entities.PbListResponseEntity

	if err := json.Unmarshal(body, &providersResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return err
	}

	if len(providersResponse.Items) == 0 {
		return fmt.Errorf("no provider with URL: %v", provider.Name)
	}

	providerModel := models.ProviderModel{
		ID:     providersResponse.Items[0]["id"].(string),
		Panel:  providersResponse.Items[0]["panel"].(string),
		URL:    providersResponse.Items[0]["url"].(string),
		Method: providersResponse.Items[0]["method"].(string),
	}

	serviceModel := &models.ServiceModel{
		ServiceID:      service.ID,
		Name:           service.Name,
		Provider:       providerModel.ID,
		HasRefill:      false,
		RefillDuration: 0,
	}

	addServiceToDb(serviceModel)

	return nil
}

func MigratePerfectPanelServices() {
	services := getServices()
	providers := getProviders()

	var failedServices []PerfectPanelService

	totalAddedServices := 0

	for _, service := range services {
		if service.Status == 1 {
			provider := getProviderFromService(service.ProviderID, &providers)

			if provider.Name != "" {
				err := pushPerfectPanelServiceToDatabase(&service, provider)

				if err != nil {
					failedServices = append(failedServices, service)
				} else {
					totalAddedServices += 1
				}
			}
		}
	}

	fmt.Printf("Added %v services.", totalAddedServices)

	var failedTxt string

	for _, service := range failedServices {
		failedTxt += fmt.Sprint(service.ID) + " " + service.Name + " " + fmt.Sprint(service.ProviderID) + "\n"
	}

	os.WriteFile("./data/failed_to_add.txt", []byte(failedTxt), 0644)
}
