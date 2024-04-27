package dtos

type DeviceInfoInput struct {
	ID             string `json:"androidId"`
	DeviceName     string `json:"deviceName"`
	DeviceBrand    string `json:"deviceBrand"`
	APILevel       int    `json:"apiLevel"`
	AndroidVersion string `json:"androidVersion"`
	Manufacturer   string `json:"manufacturer"`
	ProductName    string `json:"productName"`
	BatteryLevel   int    `json:"batteryLevel"`
	IsCharging     bool   `json:"isCharging"`
}
