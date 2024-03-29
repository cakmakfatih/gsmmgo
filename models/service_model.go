package models

type ServiceModel struct {
	Provider  string `json:"provider"`
	ServiceID string `json:"service_id"`
	Name      string `json:"name"`
}
