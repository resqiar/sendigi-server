package repos

import (
	"context"
	"log"
	"sendigi-server/configs"
	"sendigi-server/dtos"
	"sendigi-server/entities"
)

func CreateRequestMessage(payload *dtos.RequestMessageInput, userID string) error {
	SQL := `
        INSERT INTO request_messages (
            device_id, message, package_name, author_id
        ) VALUES ($1, $2, $3, $4)
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		&payload.DeviceID,
		&payload.Message,
		&payload.PackageName,
		&userID,
	); err != nil {
		return err
	}

	return nil
}

func FindRequestMessages(userID string) ([]entities.RequestMessage, error) {
	var requestMessages []entities.RequestMessage

	SQL := `
		SELECT
            rm.id, rm.device_id, rm.package_name, rm.message, rm.created_at, ai.icon, ai.name, ai.lock_status
        FROM request_messages AS rm JOIN app_info AS ai ON rm.package_name = ai.package_name WHERE rm.author_id = $1 ORDER BY rm.created_at DESC;
    `
	rows, err := configs.DB_POOL.Query(context.Background(), SQL, userID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var requestMessage entities.RequestMessage

		if err := rows.Scan(
			&requestMessage.ID,
			&requestMessage.DeviceID,
			&requestMessage.PackageName,
			&requestMessage.Message,
			&requestMessage.CreatedAt,
			&requestMessage.Icon,
			&requestMessage.Name,
			&requestMessage.LockStatus,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		requestMessages = append(requestMessages, requestMessage)
	}

	return requestMessages, nil
}
