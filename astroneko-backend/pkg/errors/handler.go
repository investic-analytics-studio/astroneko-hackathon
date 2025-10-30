package errors

import (
	"strings"

	"astroneko-backend/internal/core/domain/shared"

	"github.com/gofiber/fiber/v2"
)

// Handler provides centralized error handling for HTTP handlers
type Handler struct{}

// NewHandler creates a new error handler
func NewHandler() *Handler {
	return &Handler{}
}

// HandleRequestBodyError handles request body parsing errors
func (h *Handler) HandleRequestBodyError(c *fiber.Ctx, err error) error {
	status, response := shared.NewErrorResponse("ERR_1029", "Invalid request body")
	return c.Status(status).JSON(response)
}

// HandleValidationError handles validation errors
func (h *Handler) HandleValidationError(c *fiber.Ctx, err error) error {
	status, response := shared.NewErrorResponse("ERR_1029", err.Error())
	return c.Status(status).JSON(response)
}

// HandleNotFoundError handles not found errors
func (h *Handler) HandleNotFoundError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	status, response := shared.NewErrorResponse("ERR_404", message)
	return c.Status(status).JSON(response)
}

// HandleConflictError handles conflict errors (e.g., duplicate resources)
func (h *Handler) HandleConflictError(c *fiber.Ctx, err error) error {
	status, response := shared.NewErrorResponse("ERR_409", err.Error())
	return c.Status(status).JSON(response)
}

// HandleUnauthorizedError handles unauthorized errors
func (h *Handler) HandleUnauthorizedError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	status, response := shared.NewErrorResponse("ERR_401", message)
	return c.Status(status).JSON(response)
}

// HandleInternalServerError handles internal server errors
func (h *Handler) HandleInternalServerError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	status, response := shared.NewErrorResponse("ERR_500", message)
	return c.Status(status).JSON(response)
}

// HandleServiceError handles service-level errors with appropriate HTTP status codes
func (h *Handler) HandleServiceError(c *fiber.Ctx, err error) error {
	errStr := err.Error()

	switch {
	case strings.Contains(errStr, "not found"):
		return h.HandleNotFoundError(c, errStr)
	case strings.Contains(errStr, "already exists"):
		return h.HandleConflictError(c, err)
	case strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "invalid token"):
		return h.HandleUnauthorizedError(c, errStr)
	case strings.Contains(errStr, "invalid") && strings.Contains(errStr, "request"):
		status, response := shared.NewErrorResponse("ERR_400", errStr)
		return c.Status(status).JSON(response)
	default:
		return h.HandleInternalServerError(c, errStr)
	}
}

// HandleSuccess handles successful responses
func (h *Handler) HandleSuccess(c *fiber.Ctx, data interface{}) error {
	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = data
	return c.Status(status).JSON(response)
}

// HandleSuccessWithMessage handles successful responses with custom message
func (h *Handler) HandleSuccessWithMessage(c *fiber.Ctx, message string, data interface{}) error {
	status, response := shared.NewSuccessResponse(message)
	response.Data = data
	return c.Status(status).JSON(response)
}

// HandleCreated handles created responses (201)
func (h *Handler) HandleCreated(c *fiber.Ctx, data interface{}) error {
	status, response := shared.NewSuccessResponse("SUC_201")
	response.Data = data
	return c.Status(status).JSON(response)
}

// HandleNoContent handles no content responses (204)
func (h *Handler) HandleNoContent(c *fiber.Ctx) error {
	status, response := shared.NewSuccessResponse("SUC_204")
	return c.Status(status).JSON(response)
}
