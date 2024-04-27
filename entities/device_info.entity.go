package entities

type DeviceInfo struct {
	ID             string
	DeviceName     string
	DeviceBrand    string
	APILevel       int
	AndroidVersion string
	Manufacturer   string
	ProductName    string
	BatteryLevel   int
	IsCharging     bool
}
