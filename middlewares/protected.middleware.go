package middlewares

import (
	"log"
	"sendigi-server/configs"
	"sendigi-server/repos"
	"sendigi-server/utils"

	"github.com/gofiber/fiber/v2"
)

func ProtectedRoute(c *fiber.Ctx) error {
	sess, err := configs.SessionStore.Get(c)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userID := sess.Get("ID")
	if userID == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	_, err = repos.FindUserByID(userID.(string))
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// save user id from session into local key value
	c.Locals("userID", userID)

	return c.Next()
}

func ProtectedMobile(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
		})
	}

	userID := utils.ParseJWT(token)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
		})
	}

	_, err := repos.FindUserByID(userID)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// save to locals
	c.Locals("userID", userID)

	return c.Next()
}
