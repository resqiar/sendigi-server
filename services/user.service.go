package services

import (
	"sendigi-server/dtos"
	"sendigi-server/entities"
	"sendigi-server/repos"
	"sendigi-server/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateUser(payload *dtos.GooglePayload) (*entities.User, error) {
	// format the given name from the provider
	formattedName := utils.FormatUsername(payload.GivenName)

	newUser := &dtos.CreateUserInput{
		Username:   formattedName,
		Email:      payload.Email,
		Provider:   "google",
		PictureURL: payload.Picture,
	}

	result, err := repos.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := repos.FindUserByID(userID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(user)
}
