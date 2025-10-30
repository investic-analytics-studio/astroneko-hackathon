package routes

import (
	"astroneko-backend/internal/handlers"
	"astroneko-backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupCRMRoutes sets up all CRM-related routes
func SetupCRMRoutes(api fiber.Router, crmUserHandler *handlers.CRMUserHTTPHandler, userHandler *handlers.UserHTTPHandler, crmAuthMiddleware *middleware.CRMAuthMiddleware) {
	// CRM group
	crm := api.Group("/crm")

	// User management
	crm.Post("/users", crmUserHandler.CreateCRMUser)
	crm.Get("/users/total", crmAuthMiddleware.RequireAuth, userHandler.GetTotalUsers)

	// Authentication
	crmAuth := crm.Group("/auth")
	crmAuth.Post("/login", crmUserHandler.CRMLogin)

	// Protected routes
	crmAuth.Get("/me", crmAuthMiddleware.RequireAuth, crmUserHandler.GetCRMMe)
}
