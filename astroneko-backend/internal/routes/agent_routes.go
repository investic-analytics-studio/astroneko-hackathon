package routes

import (
	"astroneko-backend/internal/handlers"
	"astroneko-backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAgentRoutes(api fiber.Router, agentHandler *handlers.AgentHTTPHandler, authMiddleware *middleware.AuthMiddleware, guestRateLimit *middleware.GuestRateLimitMiddleware) {
	agent := api.Group("/agent")

	// Clear state requires authentication
	agent.Post("/clear-state", authMiddleware.RequireAuth, agentHandler.ClearState)

	agent.Post("/reply",
		authMiddleware.OptionalAuthWithReferralCheck,
		guestRateLimit.GuestOrAuthRateLimit("/api/v1/agent/reply", 3),
		middleware.SetupAgentReplyRateLimitMiddleware(),
		agentHandler.Reply,
	)
}
