package routes

import (
	"astroneko-backend/internal/handlers"
	"astroneko-backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, userHandler *handlers.UserHTTPHandler, authMiddleware *middleware.AuthMiddleware) {
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", userHandler.Login)
	auth.Post("/google", userHandler.GoogleAuth)
	auth.Post("/firebase", userHandler.AuthenticateWithFirebase)
	auth.Post("/refresh", userHandler.RefreshToken)
	auth.Post("/logout", userHandler.Logout)

	// Protected routes
	auth.Get("/me", authMiddleware.RequireAuth, userHandler.GetMe)

	// Referral code endpoints
	auth.Get("/referral/codes", authMiddleware.RequireAuth, userHandler.GetUserReferralCodes)
	auth.Post("/referral/activate", authMiddleware.RequireAuth, userHandler.ActivateReferral)
}
