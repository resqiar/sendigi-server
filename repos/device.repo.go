package repos

import (
	"context"
	"log"
	"sendigi-server/configs"
	"sendigi-server/dtos"
	"sendigi-server/entities"
)

func CreateDevice(payload *dtos.DeviceInfoInput, userID string) error {
	SQL := `
        INSERT INTO device_info (
            id, device_name, device_brand, api_level, android_version,
            manufacturer, product_name, battery_level, is_charging, author_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		&payload.ID,
		&payload.DeviceName,
		&payload.DeviceBrand,
		&payload.APILevel,
		&payload.AndroidVersion,
		&payload.Manufacturer,
		&payload.ProductName,
		&payload.BatteryLevel,
		&payload.IsCharging,
		&userID,
	); err != nil {
		return err
	}

	return nil
}

func UpdateDevice(payload *dtos.DeviceInfoInput, userID string) error {
	SQL := `
		UPDATE device_info
			SET device_name = $2, device_brand = $3, api_level = $4, android_version = $5,
				manufacturer = $6, product_name = $7, battery_level = $8, is_charging = $9
		WHERE id = $1 AND author_id = $10
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		&payload.ID,
		&payload.DeviceName,
		&payload.DeviceBrand,
		&payload.APILevel,
		&payload.AndroidVersion,
		&payload.Manufacturer,
		&payload.ProductName,
		&payload.BatteryLevel,
		&payload.IsCharging,
		&userID,
	); err != nil {
		return err
	}

	return nil
}

func FindDeviceById(id string) (*entities.DeviceInfo, error) {
	var deviceInfo entities.DeviceInfo

	SQL := `
        SELECT 
            id, device_name, device_brand, api_level, android_version,
            manufacturer, product_name, battery_level, is_charging
        FROM device_info WHERE id = $1
    `
	row := configs.DB_POOL.QueryRow(context.Background(), SQL, id)
	if err := row.Scan(
		&deviceInfo.ID,
		&deviceInfo.DeviceName,
		&deviceInfo.DeviceBrand,
		&deviceInfo.APILevel,
		&deviceInfo.AndroidVersion,
		&deviceInfo.Manufacturer,
		&deviceInfo.ProductName,
		&deviceInfo.BatteryLevel,
		&deviceInfo.IsCharging,
	); err != nil {
		log.Println(err)
		return nil, err
	}

	return &deviceInfo, nil
}

func FindDevices(userID string) ([]entities.DeviceInfo, error) {
	var deviceInfos []entities.DeviceInfo

	SQL := `
        SELECT 
            id, device_name, device_brand, api_level, android_version,
            manufacturer, product_name, battery_level, is_charging
        FROM device_info WHERE author_id = $1
    `
	rows, err := configs.DB_POOL.Query(context.Background(), SQL, userID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var deviceInfo entities.DeviceInfo

		if err := rows.Scan(
			&deviceInfo.ID,
			&deviceInfo.DeviceName,
			&deviceInfo.DeviceBrand,
			&deviceInfo.APILevel,
			&deviceInfo.AndroidVersion,
			&deviceInfo.Manufacturer,
			&deviceInfo.ProductName,
			&deviceInfo.BatteryLevel,
			&deviceInfo.IsCharging,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		deviceInfos = append(deviceInfos, deviceInfo)
	}

	return deviceInfos, nil
}
