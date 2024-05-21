package dtos

type NotificationConfigInput struct {
	Email          string `json:"email" validate:"required"`
	EmailStatus    bool   `json:"emailStatus"`
	Whatsapp       string `json:"whatsapp"`
	WhatsappStatus bool   `json:"whatsappStatus"`
	Telegram       string `json:"telegram"`
	TelegramStatus bool   `json:"telegramStatus"`
	Strategy       string `json:"strategy" validate:"required,oneof=OFF LOCKED ALL"`
}
