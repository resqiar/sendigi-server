package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sendigi-server/configs"
	"sendigi-server/constants"
	"sendigi-server/dtos"
	"sendigi-server/repos"
	"sendigi-server/utils"
	"time"

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

func MobileSyncDeviceActivity(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	rawPayload := c.Body()

	var payload dtos.DeviceActivityInput

	err := json.Unmarshal([]byte(rawPayload), &payload)
	if err != nil {
		log.Println("Unmarshal Error:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	timeDiff, err := configs.RedisStore.Get(fmt.Sprintf("%s_last_inserted_time", userID))
	if err != nil {
		log.Println("Failed to retrieve redis value:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// if there is a last insertion before, if not, just skip the whole check
	if string(timeDiff) != "" {
		parsedDate, err := time.Parse(time.RFC3339, string(timeDiff))
		if err != nil {
			log.Println("Failed to parse redis value:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// check the last insertion time saved in redis.
		// if time from the last caller is less than 5 seconds,
		// return as no content and just skip the next process entirely.
		if time.Since(parsedDate).Seconds() < 5 {
			configs.RedisStore.Set(fmt.Sprintf("%s_last_inserted_time", userID), []byte(time.Now().Format(time.RFC3339)), 1*time.Minute)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	if err := repos.CreateDeviceActivity(&payload, userID); err != nil {
		log.Printf("Failed to insert device activity: %v", err)
		return c.JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
		})
	}

	payloadChan := make(chan dtos.DeviceActivityInput)

	go func() {
		// get the payload from the channel,
		// if the channel is not having any value yet, it will block the goroutine
		// until the payload is present and vice versa.
		payloadCopy := <-payloadChan

		payloadCopy.CreatedAt = time.Now()

		config, err := repos.FindUserNotificationConfig(userID)
		if err != nil {
			log.Printf("Failed to get notification config: %v", err)
		}

		// if strategy for current user is off, just skip sending to mq
		if config.Strategy == constants.NOTIF_STRATEGY_OFF {
			return
		}

		appInfo, err := repos.FindAppByPackageName(payloadCopy.PackageName, userID)
		if err != nil {
			log.Printf("Failed to get app by package name: %v", err)
			return
		}

		switch {
		case config.Strategy == constants.NOTIF_STRATEGY_LOCKED && appInfo.LockStatus:
			sendToNotifMQ(userID, &payload)
		case config.Strategy == constants.NOTIF_STRATEGY_ALL:
			sendToNotifMQ(userID, &payload)
		}
	}()

	// Send the original payload to the channel
	payloadChan <- payload

	// set the key as inserted inside redis for max 10s
	// why save it longer since the data itself only needed for 3s
	configs.RedisStore.Set(fmt.Sprintf("%s_last_inserted_time", userID), []byte(time.Now().Format(time.RFC3339)), 10*time.Second)

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
	})
}

func sendToNotifMQ(userID string, payload *dtos.DeviceActivityInput) error {
	// force addition of user id
	payload.UserID = userID

	ch, err := configs.InitChannel()
	if err != nil {
		log.Printf("[Notif Queue] Failed to init channel: %v", err)
		return err
	}

	q, ctx, cancel, err := configs.InitNotifQueue(ch, userID)
	if err != nil {
		log.Printf("[Notif Queue] Failed to init queue: %v", err)
		return err
	}
	defer cancel()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Notif Queue] Failed to marshall payload: %v", err)
		return err
	}

	if err := configs.SendToNotifMQ(ch, q, ctx, jsonPayload); err != nil {
		log.Printf("[Notif Queue] Failed to send to queue: %v", err)
		return err
	}

	return nil
}

func MobileGetDevices(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	devices, err := repos.FindDevices(userID)
	if err != nil {
		log.Printf("Failed to get devices: %v", err)
		return c.JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{"data": devices})
}

func MobileGetActivities(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	activities, err := repos.FindDeviceActivities(userID)
	if err != nil {
		log.Printf("Failed to get device activities: %v", err)
		return c.JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{"data": activities})
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

	// check if the app already exist
	existingApp, err := repos.FindAppByPackageName(payload.PackageName, userID)
	if existingApp == nil || err != nil { // app info is new
		url, err := utils.UploadImages(payload.Icon, payload.PackageName)
		if err != nil {
			log.Println("Error Uploading Icon:", err)
			return c.JSON(fiber.Map{
				"status": fiber.StatusInternalServerError,
			})
		}

		// force icon to use remotely stored icon
		payload.Icon = url

		err = repos.CreateAppInfo(&payload, userID)
		if err != nil {
			log.Printf("Failed to register app info: %v", err)
			return c.JSON(fiber.Map{
				"status": fiber.StatusInternalServerError,
			})
		}
	} else { // app info is present and need to be updated
		err := repos.UpdateAppInfo(&payload, userID)
		if err != nil {
			log.Printf("Failed to update app info: %v", err)
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
