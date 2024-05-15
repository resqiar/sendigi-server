package dtos

import "database/sql"

type NotificationConfigInput struct {
	Email          string         `json:"email" validate:"required"`
	EmailStatus    bool           `json:"emailStatus"`
	Whatsapp       sql.NullString `json:"whatsapp"`
	WhatsappStatus bool           `json:"whatsappStatus"`
	Strategy       string         `json:"strategy" validate:"required,oneof=OFF LOCKED ALL"`
}
