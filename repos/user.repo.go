package repos

import (
	"context"
	"sendigi-server/configs"
	"sendigi-server/dtos"
	"sendigi-server/entities"
)

func CreateUser(input *dtos.CreateUserInput) (*entities.User, error) {
	var user entities.User

	SQL := "INSERT INTO users(provider, username, email, picture_url) VALUES ($1, $2, $3, $4) RETURNING id, username, email"
	row := configs.DB_POOL.QueryRow(context.Background(), SQL, input.Provider, input.Username, input.Email, input.PictureURL)
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	); err != nil {
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
