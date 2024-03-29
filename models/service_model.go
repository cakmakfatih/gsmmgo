package models

type ServiceModel struct {
	Provider  string `json:"provider"`
	ServiceID int    `json:"service_id"`
	Name      string `json:"name"`
}
