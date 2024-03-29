package models

type ProviderModel struct {
	ID         string `json:"id"`
	Panel      string `json:"panel"`
	URL        string `json:"url"`
	Alias      string `json:"alias"`
	Method     string `json:"method"`
	MethodData string `json:"method_data"`
}
