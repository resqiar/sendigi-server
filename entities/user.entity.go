package entities

import (
	"database/sql"
	"time"
)

type User struct {
	ID         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Provider   string
	Username   string
	Email      string
	Bio        sql.NullString
	PictureURL string
}
