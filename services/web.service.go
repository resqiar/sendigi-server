package services

import (
	"encoding/json"
	"log"
	"sendigi-server/configs"
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

	// update message queue async-ly to trigger android device
	go func() {
		err := SendStateToMobile(userID, payload)
		if err != nil {
			log.Println(err)
		}
	}()

	return c.SendStatus(fiber.StatusOK)
}

func SendStateToMobile(userID string, payload dtos.AppInfoInput) error {
	// get device information
	devices, err := repos.FindDevices(userID)
	if err != nil {
		log.Printf("[Queue] Failed to get devices: %v", err)
		return err
	}

	ch, err := configs.InitChannel()
	if err != nil {
		log.Printf("[Queue] Failed to init channel: %v", err)
		return err
	}

	q, ctx, cancel, err := configs.InitMobileQueue(ch, userID, devices[0].ID)
	if err != nil {
		log.Printf("[Queue] Failed to init queue: %v", err)
		return err
	}
	defer cancel()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Queue] Failed to marshall payload: %v", err)
		return err
	}

	// send to mobile queue
	if err := configs.SendToMQ(ch, q, ctx, jsonPayload); err != nil {
		log.Printf("[Queue] Failed to send to queue: %v", err)
		return err
	}

	return nil
}

func GetNotificationConfig(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	config, err := repos.FindUserNotificationConfig(userID)
	if err != nil {
		log.Println("Error getting notification config:", err)
		// send raw error (unprocessable entity)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"data": config,
	})
}

func UpdateNotificationConfig(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var payload dtos.NotificationConfigInput

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

	err := repos.UpdateUserNotificationConfig(&payload, userID)
	if err != nil {
		log.Println("Error updating notification config:", err)
		// send raw error (unprocessable entity)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
