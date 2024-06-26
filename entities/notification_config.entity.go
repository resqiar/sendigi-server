package entities

import "database/sql"

type NotificationConfig struct {
	ID             int
	Email          string
	EmailStatus    bool
	Whatsapp       sql.NullString
	WhatsappStatus bool
	Telegram       sql.NullString
	TelegramStatus bool
	Strategy       string
}
