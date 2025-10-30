package middleware

import (
	"strings"

	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type CRMAuthMiddleware struct {
	crmUserService *services.CRMUserService
	logger         logger.Logger
}

func NewCRMAuthMiddleware(crmUserService *services.CRMUserService, logger logger.Logger) *CRMAuthMiddleware {
	return &CRMAuthMiddleware{
		crmUserService: crmUserService,
		logger:         logger,
	}
}

func (m *CRMAuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		status, response := shared.NewErrorResponse("ERR_401", "Authorization header required")
		return c.Status(status).JSON(response)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		status, response := shared.NewErrorResponse("ERR_401", "Invalid authorization header format")
		return c.Status(status).JSON(response)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		status, response := shared.NewErrorResponse("ERR_401", "Token required")
		return c.Status(status).JSON(response)
	}

	claims, err := m.crmUserService.ValidateToken(token)
	if err != nil {
		m.logger.Error("Failed to validate CRM token", logger.Field{Key: "error", Value: err.Error()})
		status, response := shared.NewErrorResponse("ERR_401", "Invalid or expired token")
		return c.Status(status).JSON(response)
	}

	// Get full user details
	user, err := m.crmUserService.GetUserByID(c.Context(), claims.UserID.String())
	if err != nil {
		m.logger.Error("Failed to get CRM user",
			logger.Field{Key: "error", Value: err.Error()},
			logger.Field{Key: "user_id", Value: claims.UserID.String()})
		status, response := shared.NewErrorResponse("ERR_401", "User not found")
		return c.Status(status).JSON(response)
	}

	// Store user in context for handlers to use
	c.Locals("crm_user", user)
	c.Locals("crm_user_id", claims.UserID.String())

	return c.Next()
}
