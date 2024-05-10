package services

import (
	"log"
	"sendigi-server/dtos"
	"sendigi-server/repos"
	"sendigi-server/utils"

	"github.com/gofiber/fiber/v2"
)

func WebUpdateApps(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var payload dtos.AppInfoInput

	// bind the body parser into payload
	if err := c.BodyParser(&payload); err != nil {
		log.Println("Parsing Error:", err)
		// send raw error (unprocessable entity)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// validate the payload using class-validator
	if err := utils.ValidateInput(payload); err != "" {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err,
		})
	}

	// TODO: Implement better message later on
	if payload.PackageName == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := repos.UpdateAppInfoWeb(&payload, userID); err != nil {
		log.Printf("Failed to update apps: %v", err)
		return c.JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
