package entities

import "time"

type RequestMessage struct {
	ID          string
	Icon        string
	Message     string
	Ack         bool
	Name        string
	PackageName string
	LockStatus  bool
	DeviceID    string
	CreatedAt   time.Time
}
