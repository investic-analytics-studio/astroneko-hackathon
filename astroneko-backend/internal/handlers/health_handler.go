package handlers

import (
	"time"

	"astroneko-backend/internal/core/domain/user"

	"github.com/gofiber/fiber/v2"
)

type HealthHTTPHandler struct{}

func NewHealthHTTPHandler() *HealthHTTPHandler {
	return &HealthHTTPHandler{}
}

type HealthResponse struct {
	Status    string             `json:"status"`
	Message   string             `json:"message"`
	Service   string             `json:"service"`
	Timestamp string             `json:"timestamp"`
	User      *user.UserResponse `json:"user,omitempty"`
}

type SimpleHealthResponse struct {
	Status string `json:"status"`
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the API is healthy and running, shows authenticated user if token provided
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Security BearerAuth
// @Router /health-check [get]
func (h *HealthHTTPHandler) HealthCheck(c *fiber.Ctx) error {
	response := HealthResponse{
		Status:    "OK",
		Message:   "API is healthy",
		Service:   "astroneko-backend",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Check if user is authenticated (added by middleware)
	if userData := c.Locals("user"); userData != nil {
		if authUser, ok := userData.(*user.User); ok {
			response.User = authUser.ToResponse()
			response.Message = "API is healthy - Authenticated user found"
		}
	}

	return c.Status(200).JSON(response)
}

// SimpleHealthCheck godoc
// @Summary Simple health check endpoint
// @Description Simple health check that returns basic status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} SimpleHealthResponse
// @Router /health [get]
func (h *HealthHTTPHandler) SimpleHealthCheck(c *fiber.Ctx) error {
	response := SimpleHealthResponse{
		Status: "ok",
	}
	return c.Status(200).JSON(response)
}
