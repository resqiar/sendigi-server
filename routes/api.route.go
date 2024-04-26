package routes

import (
	"sendigi-server/middlewares"
	"sendigi-server/services"

	"github.com/gofiber/fiber/v2"
)

func InitAPIRoutes(s *fiber.App) {
	auth := s.Group("auth")

	auth.Get("/google", services.GoogleLoginService)
	auth.Get("/google/callback", services.GoogleLoginCallbackService)

	auth.Get("/check", middlewares.ProtectedRoute, func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	auth.Get("/logout", services.SendLogout)
}
