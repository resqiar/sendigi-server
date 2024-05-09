package repos

import (
	"context"
	"log"
	"sendigi-server/configs"
	"sendigi-server/dtos"
)

func CreateAppInfo(payload *dtos.AppInfoInput, userID string) error {
	SQL := `
		INSERT INTO app_info (
			name,
			package_name,
			lock_status,
			icon, 
			time_usage, 
			author_id, 
			date_locked, 
			time_start_locked, 
			time_end_locked
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		&payload.Name,
		&payload.PackageName,
		&payload.LockStatus,
		&payload.Icon,
		&payload.TimeUsage,
		&userID,
		&payload.DateLocked,
		&payload.TimeStartLocked,
		&payload.TimeEndLocked,
	); err != nil {
		return err
	}

	return nil
}

func UpdateAppInfo(payload *dtos.AppInfoInput, userID string) error {
	SQL := `
		UPDATE app_info
			SET name = $1,
			lock_status = $3, 
			icon = $4, 
			time_usage = $5,
			date_locked = $7, 
			time_start_locked = $8, 
			time_end_locked = $9
		WHERE package_name = $2 AND author_id = $6
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		&payload.Name,
		&payload.PackageName,
		&payload.LockStatus,
		&payload.Icon,
		&payload.TimeUsage,
		&userID,
		&payload.DateLocked,
		&payload.TimeStartLocked,
		&payload.TimeEndLocked,
	); err != nil {
		return err
	}

	return nil
}

func FindAppByPackageName(packageName string) (*dtos.AppInfo, error) {
	var appInfo dtos.AppInfo

	SQL := `
        SELECT 
			id, 
			name, 
			package_name, 
			lock_status, 
			icon, 
			time_usage,
			date_locked, 
			time_start_locked, 
			time_end_locked
        FROM app_info WHERE package_name = $1
    `
	row := configs.DB_POOL.QueryRow(context.Background(), SQL, packageName)
	if err := row.Scan(
		&appInfo.ID,
		&appInfo.Name,
		&appInfo.PackageName,
		&appInfo.LockStatus,
		&appInfo.Icon,
		&appInfo.TimeUsage,
		&appInfo.DateLocked,
		&appInfo.TimeStartLocked,
		&appInfo.TimeEndLocked,
	); err != nil {
		log.Println(err)
		return nil, err
	}

	return &appInfo, nil
}

func FindApps(userID string) ([]dtos.AppInfo, error) {
	var appInfos []dtos.AppInfo

	SQL := `
        SELECT 
			id, 
			name, 
			package_name, 
			lock_status, 
			icon, 
			time_usage,
			date_locked, 
			time_start_locked, 
			time_end_locked
        FROM app_info WHERE author_id = $1
    `
	rows, err := configs.DB_POOL.Query(context.Background(), SQL, userID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var appInfo dtos.AppInfo

		if err := rows.Scan(
			&appInfo.ID,
			&appInfo.Name,
			&appInfo.PackageName,
			&appInfo.LockStatus,
			&appInfo.Icon,
			&appInfo.TimeUsage,
			&appInfo.DateLocked,
			&appInfo.TimeStartLocked,
			&appInfo.TimeEndLocked,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		appInfos = append(appInfos, appInfo)
	}

	return appInfos, nil
}
