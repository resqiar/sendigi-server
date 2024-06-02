package dtos

import "time"

type RequestMessageInput struct {
	Message     string `json:"message" validate:"required"`
	PackageName string `json:"packageName" validate:"required"`
	DeviceID    string `json:"deviceId" validate:"required"`
	UserID      string
	CreatedAt   time.Time
}

type MessageInput struct {
	Message     string `json:"message"`
	PackageName string `json:"packageName" validate:"required"`
	LockStatus  bool   `json:"lockStatus"`
	UserID      string
}
