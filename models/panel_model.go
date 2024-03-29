package models

type PanelModel struct {
	ID              string `json:"id"`
	LoginURL        string `json:"login_url"`
	SupportUsername string `json:"support_username"`
	SupportPassword string `json:"support_password"`
	TelegramToken   string `json:"telegram_token"`
	WhatsAppToken   string `json:"whatsapp_token"`
}
