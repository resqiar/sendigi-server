package entities

import (
	"database/sql"
	"time"
)

type DeviceActivity struct {
	ID          string
	Name        string
	Icon        string
	Description sql.NullString
	PackageName string
	DeviceID    string
	CreatedAt   time.Time
}
