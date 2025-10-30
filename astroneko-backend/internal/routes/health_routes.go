package routes

import (
	"astroneko-backend/internal/handlers"
	"astroneko-backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupHealthRoutes(app *fiber.App, api fiber.Router, healthHandler *handlers.HealthHTTPHandler, authMiddleware *middleware.AuthMiddleware) {
	// Simple health check (no auth)
	app.Get("/health", healthHandler.SimpleHealthCheck)

	// Health check with optional auth (shows user data if authenticated)
	app.Get("/health-check", authMiddleware.OptionalAuth, healthHandler.HealthCheck)
}
