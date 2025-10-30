package routes

import (
	"astroneko-backend/internal/handlers"
	"astroneko-backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupReferralCodeRoutes sets up all referral code management routes
func SetupReferralCodeRoutes(api fiber.Router, handler *handlers.ReferralCodeHTTPHandler, crmAuthMiddleware *middleware.CRMAuthMiddleware) {
	// General referral codes group (admin only)
	referralCodesGroup := api.Group("/referral-codes")

	// Apply CRM authentication middleware to all routes
	referralCodesGroup.Use(crmAuthMiddleware.RequireAuth)

	// CRUD operations for general referral codes
	referralCodesGroup.Post("/", handler.CreateReferralCode)
	referralCodesGroup.Get("/", handler.ListReferralCodes)
	referralCodesGroup.Get("/:id", handler.GetReferralCodeByID)
	referralCodesGroup.Put("/:id", handler.UpdateReferralCode)
	referralCodesGroup.Delete("/:id", handler.DeleteReferralCode)

	// Additional utility endpoints
	referralCodesGroup.Get("/code/:code", handler.GetReferralCodeByCode)
	referralCodesGroup.Get("/validate/:code", handler.ValidateReferralCode)
}
