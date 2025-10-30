package routes

import (
	"astroneko-backend/internal/handlers"
	"astroneko-backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupHistoryRoutes configures all history-related routes
func SetupHistoryRoutes(api fiber.Router, historyHandler *handlers.HistoryHTTPHandler, authMiddleware *middleware.AuthMiddleware) {
	history := api.Group("/history")

	// Protected routes - require authentication
	history.Get("/sessions", authMiddleware.RequireAuth, historyHandler.GetUserSessions)
	history.Get("/sessions/:session_id/messages", authMiddleware.RequireAuth, historyHandler.GetSessionMessages)
	history.Delete("/sessions/:session_id", authMiddleware.RequireAuth, historyHandler.DeleteSession)
}
