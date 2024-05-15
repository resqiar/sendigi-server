package repos

import (
	"context"
	"sendigi-server/configs"
	"sendigi-server/dtos"
	"sendigi-server/entities"

	"github.com/jackc/pgx/v5"
)

func CreateUser(input *dtos.CreateUserInput) (*entities.User, error) {
	var user entities.User

	tx, err := configs.DB_POOL.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	SQL := "INSERT INTO users(provider, username, email, picture_url) VALUES ($1, $2, $3, $4) RETURNING id, username, email"
	row := tx.QueryRow(context.Background(), SQL, input.Provider, input.Username, input.Email, input.PictureURL)
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	); err != nil {
		return nil, err
	}

	SQL = "INSERT INTO notification_config(email, user_id) VALUES ($1, $2)"
	_, err = tx.Exec(context.Background(), SQL, user.Email, user.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User

	SQL := "SELECT id, username, email, bio, picture_url FROM users WHERE email = $1"
	row := configs.DB_POOL.QueryRow(context.Background(), SQL, email)
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.PictureURL,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserByID(id string) (*entities.User, error) {
	var user entities.User

	SQL := "SELECT id, username, email, bio, picture_url FROM users WHERE id = $1"
	row := configs.DB_POOL.QueryRow(context.Background(), SQL, id)
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.PictureURL,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserNotificationConfig(userId string) (*entities.NotificationConfig, error) {
	var config entities.NotificationConfig

	SQL := "SELECT id, email, email_status, whatsapp, whatsapp_status, strategy FROM notification_config WHERE user_id = $1"
	row := configs.DB_POOL.QueryRow(context.Background(), SQL, userId)
	if err := row.Scan(
		&config.ID,
		&config.Email,
		&config.EmailStatus,
		&config.Whatsapp,
		&config.WhatsappStatus,
		&config.Strategy,
	); err != nil {
		return nil, err
	}

	return &config, nil
}

func UpdateUserNotificationConfig(payload *dtos.NotificationConfigInput, userID string) error {
	SQL := `
		UPDATE notification_config
			SET email = $1,
			email_status = $2, 
			whatsapp = $3, 
			whatsapp_status = $4,
			strategy = $5
		WHERE user_id = $6
    `
	if _, err := configs.DB_POOL.Exec(
		context.Background(),
		SQL,
		&payload.Email,
		&payload.EmailStatus,
		&payload.Whatsapp,
		&payload.WhatsappStatus,
		&payload.Strategy,
		&userID,
	); err != nil {
		return err
	}

	return nil
}
