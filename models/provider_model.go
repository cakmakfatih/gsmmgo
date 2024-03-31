package models

type ProviderModel struct {
	ID                 string            `json:"id"`
	Panel              string            `json:"panel"`
	URL                string            `json:"url"`
	Alias              string            `json:"alias"`
	Method             string            `json:"method"`
	MethodData         map[string]string `json:"method_data,omitempty"`
	MethodDataReadable string            `json:"-"`
}

func (p *ProviderModel) SetMethodDataReadableFromJsonString(methodData map[string]interface{}) {
	if p.Method == "telegram" {
		p.MethodDataReadable = methodData["telegram_chat_id"].(string)
	} else if p.Method == "web" {
		p.MethodDataReadable += methodData["support_username"].(string) + ":" + methodData["support_password"].(string)
	}
}

func (p *ProviderModel) SetMethodDataReadableFromSelf() {
	if p.Method == "telegram" {
		p.MethodDataReadable = p.MethodData["telegram_chat_id"]
	} else if p.Method == "web" {
		p.MethodDataReadable += p.MethodData["support_username"] + ":" + p.MethodData["support_password"]
	}
}
