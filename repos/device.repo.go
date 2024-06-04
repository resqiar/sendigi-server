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

func CreateDeviceActivity(payload *dtos.DeviceActivityInput, userID string) error {
	SQL := `
        INSERT INTO device_activity (
            device_id, name, description, package_name, author_id
        ) VALUES ($1, $2, $3, $4, $5)
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		&payload.DeviceID,
		&payload.Name,
		&payload.Description,
		&payload.PackageName,
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

func FindDevice(userID string, deviceID string) (*entities.DeviceInfo, error) {
	var deviceInfo entities.DeviceInfo

	SQL := `
        SELECT 
            id, device_name, device_brand, api_level, android_version,
            manufacturer, product_name, battery_level, is_charging
        FROM device_info WHERE author_id = $1 AND id = $2
    `
	row := configs.DB_POOL.QueryRow(context.Background(), SQL, userID, deviceID)
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

func FindDeviceActivities(userID string, deviceID string) ([]entities.DeviceActivity, error) {
	var deviceActivities []entities.DeviceActivity

	SQL := `
		SELECT
            da.id, da.device_id, da.name, da.package_name, da.description, da.created_at, ai.icon
        FROM device_activity AS da JOIN app_info AS ai ON da.package_name = ai.package_name WHERE da.author_id = $1 AND da.device_id = $2 ORDER BY da.created_at DESC;
    `
	rows, err := configs.DB_POOL.Query(context.Background(), SQL, userID, deviceID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var deviceActivity entities.DeviceActivity

		if err := rows.Scan(
			&deviceActivity.ID,
			&deviceActivity.DeviceID,
			&deviceActivity.Name,
			&deviceActivity.PackageName,
			&deviceActivity.Description,
			&deviceActivity.CreatedAt,
			&deviceActivity.Icon,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		deviceActivities = append(deviceActivities, deviceActivity)
	}

	return deviceActivities, nil
}
