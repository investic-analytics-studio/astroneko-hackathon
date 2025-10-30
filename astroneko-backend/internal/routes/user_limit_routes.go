package routes

import (
	"astroneko-backend/internal/handlers"
	"astroneko-backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserLimitRoutes(api fiber.Router, userLimitHandler *handlers.UserLimitHTTPHandler, crmAuthMiddleware *middleware.CRMAuthMiddleware, authMiddleware *middleware.AuthMiddleware) {
	// Public route for getting user limit
	api.Get("/user-limit", userLimitHandler.GetUserLimit)

	// User authenticated route for checking user limit
	api.Get("/user-limit/check", authMiddleware.RequireAuth, userLimitHandler.IsUserOverLimitUsed)

	// CRM protected route for updating user limit
	crm := api.Group("/crm")
	crm.Put("/user-limit", crmAuthMiddleware.RequireAuth, userLimitHandler.UpdateUserLimit)
}
