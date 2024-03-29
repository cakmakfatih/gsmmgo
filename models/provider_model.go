package models

type ProviderModel struct {
	ID         string `json:"id"`
	Panel      string `json:"panel"`
	URL        string `json:"url"`
	Alias      string `json:"alias"`
	Method     string `json:"method"`
	MethodData string `json:"method_data"`
}

func (p *ProviderModel) SetMethodDataFromJsonString(methodData map[string]interface{}) {
	if val, ok := methodData["support_username"]; ok {
		p.MethodData += val.(string)
	}

	if val, ok := methodData["support_password"]; ok {
		p.MethodData = p.MethodData + ":" + val.(string)
	}
}
