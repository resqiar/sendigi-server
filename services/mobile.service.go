package services

import (
	"encoding/json"
	"log"
	"sendigi-server/dtos"
	"sendigi-server/repos"

	"github.com/gofiber/fiber/v2"
)

func MobileSyncDevice(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	rawPayload := c.Body()

	var payload dtos.DeviceInfoInput

	err := json.Unmarshal([]byte(rawPayload), &payload)
	if err != nil {
		log.Println("Unmarshal Error:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// check if the device already exist
	existingDevice, _ := repos.FindDeviceById(payload.ID)
	if existingDevice == nil { // device info is new
		err := repos.CreateDevice(&payload, userID)
		if err != nil {
			log.Printf("Failed to register device: %v", err)
			return c.JSON(fiber.Map{
				"status": fiber.StatusInternalServerError,
			})
		}
	} else { // device info is present and need to be updated
		err := repos.UpdateDevice(&payload, userID)
		if err != nil {
			log.Printf("Failed to update device: %v", err)
			return c.JSON(fiber.Map{
				"status": fiber.StatusInternalServerError,
			})
		}
	}

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
	})
}

func MobileGetDevices(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// check if the device already exist
	devices, err := repos.FindDevices(userID)
	if err != nil {
		log.Printf("Failed to get devices: %v", err)
		return c.JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{"data": devices})
}

func MobileSyncApp(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	rawPayload := c.Body()

	var payload dtos.AppInfoInput

	err := json.Unmarshal([]byte(rawPayload), &payload)
	if err != nil {
		log.Println("Unmarshal Error:", err)
		return c.JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
		})
	}

	// check if the device already exist
	existingApp, _ := repos.FindAppByPackageName(payload.PackageName)
	if existingApp == nil { // app info is new
		err := repos.CreateAppInfo(&payload, userID)
		if err != nil {
			log.Printf("Failed to register device: %v", err)
			return c.JSON(fiber.Map{
				"status": fiber.StatusInternalServerError,
			})
		}
	} else { // device info is present and need to be updated
		err := repos.UpdateAppInfo(&payload, userID)
		if err != nil {
			log.Printf("Failed to update device: %v", err)
			return c.JSON(fiber.Map{
				"status": fiber.StatusInternalServerError,
			})
		}
	}

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
	})
}

func MobileGetApps(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// check if the device already exist
	devices, err := repos.FindApps(userID)
	if err != nil {
		log.Printf("Failed to get apps: %v", err)
		return c.JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{"data": devices})
}
