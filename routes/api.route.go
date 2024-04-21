package routes

import (
	"sendigi-server/services"

	"github.com/gofiber/fiber/v2"
)

func InitAPIRoutes(s *fiber.App) {
	auth := s.Group("auth")

	auth.Get("/google", services.GoogleLoginService)
	auth.Get("/google/callback", services.GoogleLoginCallbackService)
}
