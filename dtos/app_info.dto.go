package dtos

type AppInfo struct {
	ID          string
	Name        string
	PackageName string
	LockStatus  bool
	Icon        string
	TimeUsage   int64
	AuthorID    string
}

type AppInfoInput struct {
	Name        string `json:"name"`
	PackageName string `json:"packageName"`
	LockStatus  bool   `json:"lockStatus"`
	Icon        string `json:"icon"`
	TimeUsage   int64  `json:"timeUsage"`
}
