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
			time_end_locked,
			recurring
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		&payload.Recurring,
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
			time_usage = $4,
			date_locked = $6, 
			time_start_locked = $7, 
			time_end_locked = $8,
			recurring = $9
		WHERE package_name = $2 AND author_id = $5
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		&payload.Name,
		&payload.PackageName,
		&payload.LockStatus,
		&payload.TimeUsage,
		&userID,
		&payload.DateLocked,
		&payload.TimeStartLocked,
		&payload.TimeEndLocked,
		&payload.Recurring,
	); err != nil {
		return err
	}

	return nil
}

func UpdateAppInfoWeb(payload *dtos.AppInfoInput, userID string) error {
	SQL := `
		UPDATE app_info
			SET name = $1,
			lock_status = $3, 
			date_locked = $5, 
			time_start_locked = $6, 
			time_end_locked = $7,
			recurring = $8
		WHERE package_name = $2 AND author_id = $4
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		&payload.Name,
		&payload.PackageName,
		&payload.LockStatus,
		&userID,
		&payload.DateLocked,
		&payload.TimeStartLocked,
		&payload.TimeEndLocked,
		&payload.Recurring,
	); err != nil {
		return err
	}

	return nil
}

func FindAppByPackageName(packageName string, userId string) (*dtos.AppInfo, error) {
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
			time_end_locked,
			recurring
        FROM app_info WHERE package_name = $1 AND author_id = $2
    `
	row := configs.DB_POOL.QueryRow(context.Background(), SQL, packageName, userId)
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
		&appInfo.Recurring,
	); err != nil {
		log.Println(err)
		return nil, err
	}

	return &appInfo, nil
}

func ResetAppLockByPackageName(packageName string, userID string) error {
	SQL := `
		UPDATE app_info
			SET lock_status = $1, 
			date_locked = $2, 
			time_start_locked = $3, 
			time_end_locked = $4
		WHERE package_name = $5 AND author_id = $6
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		false,
		"",
		"",
		"",
		packageName,
		userID,
	); err != nil {
		return err
	}

	return nil
}

func SetAppLockByPackageName(packageName string, userID string) error {
	SQL := `
		UPDATE app_info
			SET lock_status = $1
		WHERE package_name = $2 AND author_id = $3
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		true,
		packageName,
		userID,
	); err != nil {
		return err
	}

	return nil
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
			time_end_locked,
			recurring
        FROM app_info WHERE author_id = $1 ORDER BY time_usage DESC
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
			&appInfo.Recurring,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		appInfos = append(appInfos, appInfo)
	}

	return appInfos, nil
}
