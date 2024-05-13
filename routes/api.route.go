package routes

import (
	"sendigi-server/middlewares"
	"sendigi-server/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func InitAPIRoutes(s *fiber.App) {
	// health monitor
	s.Get("/health", monitor.New(monitor.Config{
		Title: "SenDigi Monitoring System",
	}))

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
	mobile.Post("/sync-device-activity", middlewares.ProtectedMobile, services.MobileSyncDeviceActivity)
	mobile.Post("/sync-app", middlewares.ProtectedMobile, services.MobileSyncApp)

	api := s.Group("api")
	api.Get("/devices", middlewares.ProtectedRoute, services.MobileGetDevices)
	api.Get("/activities", middlewares.ProtectedRoute, services.MobileGetActivities)
	api.Get("/apps", middlewares.ProtectedRoute, services.MobileGetApps)
	api.Post("/apps/update", middlewares.ProtectedRoute, services.WebUpdateApps)
}
