package models

type ServiceModel struct {
	Provider       string `json:"provider"`
	ServiceID      int    `json:"service_id"`
	Name           string `json:"name"`
	HasRefill      bool   `json:"has_refill"`
	RefillDuration int    `json:"refill_duration,omitempty"`
}
