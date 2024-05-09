package dtos

import "database/sql"

type AppInfo struct {
	ID              string
	Name            string
	PackageName     string
	LockStatus      bool
	Icon            string
	TimeUsage       int64
	DateLocked      sql.NullString
	TimeStartLocked sql.NullString
	TimeEndLocked   sql.NullString
	AuthorID        string
}

type AppInfoInput struct {
	Name            string `json:"name"`
	PackageName     string `json:"packageName"`
	LockStatus      bool   `json:"lockStatus"`
	Icon            string `json:"icon"`
	TimeUsage       int64  `json:"timeUsage"`
	DateLocked      string `json:"dateLocked"`
	TimeStartLocked string `json:"timeStartLocked"`
	TimeEndLocked   string `json:"timeEndLocked"`
}
