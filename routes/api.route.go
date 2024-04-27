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

	user := s.Group("user")
	user.Get("/profile", middlewares.ProtectedRoute, services.GetUserProfile)

	mobile := s.Group("mobile")
	mobile.Post("/sync-user", services.MobileLoginCallbackService)
	mobile.Post("/sync-device", middlewares.ProtectedMobile, services.MobileSyncDevice)

	api := s.Group("api")
	api.Get("/devices", middlewares.ProtectedRoute, services.MobileGetDevices)
}
