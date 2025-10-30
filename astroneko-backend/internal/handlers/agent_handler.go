package handlers

import (
	"strings"

	"astroneko-backend/internal/core/domain/agent"
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AgentHTTPHandler struct {
	agentService *services.AgentService
	validator    validator.Validator
}

func NewAgentHTTPHandler(agentService *services.AgentService, validator validator.Validator) *AgentHTTPHandler {
	return &AgentHTTPHandler{
		agentService: agentService,
		validator:    validator,
	}
}

// ClearState godoc
// @Summary Clear agent state for authenticated user
// @Description Clear the conversation state for the cat fortune agent
// @Tags agent
// @Accept json
// @Produce json
// @Param clear_state body agent.ClearStateRequest false "Clear state request (user_id will be set from auth)"
// @Success 200 {object} agent.ClearStateResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/agent/clear-state [post]
func (h *AgentHTTPHandler) ClearState(c *fiber.Ctx) error {
	userFromContext := c.Locals("user")
	if userFromContext == nil {
		status, response := shared.NewErrorResponse("ERR_401", "User not found in context")
		return c.Status(status).JSON(response)
	}

	userEntity, ok := userFromContext.(*user.User)
	if !ok {
		status, response := shared.NewErrorResponse("ERR_401", "Invalid user data in context")
		return c.Status(status).JSON(response)
	}

	var req agent.ClearStateRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	agentResponse, err := h.agentService.ClearState(c.Context(), userEntity.ID.String(), req)
	if err != nil {
		if strings.Contains(err.Error(), "failed to make request") {
			status, response := shared.NewErrorResponse("ERR_502", "External service unavailable")
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Failed to clear agent state")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = agentResponse
	return c.Status(status).JSON(response)
}

// Reply godoc
// @Summary Send message to agent and get reply
// @Description Send a message to the cat fortune agent and receive a response. Works for both authenticated users (unlimited) and guests (3 requests/day). For authenticated users, user_id is automatically extracted from auth token. For guests, session fingerprint is used. session_id is optional.
// @Tags agent
// @Accept json
// @Produce json
// @Param reply body agent.ReplyRequest true "Message to send to agent (session_id is optional)"
// @Success 200 {object} agent.ReplyResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 429 {object} shared.ResponseBody "Rate limit exceeded for guest users"
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/agent/reply [post]
func (h *AgentHTTPHandler) Reply(c *fiber.Ctx) error {
	// Get user from context (optional - may be nil for guests)
	userFromContext := c.Locals("user")

	// Determine user ID (use authenticated user ID or generate guest session ID)
	var userID string
	var isGuest bool

	if userFromContext == nil {
		// Guest user OR logged-in user without activated referral (both treated as guest with 3/day limit)
		isGuest = true
		// Use IP + user agent hash as guest session ID (from fingerprint)
		fingerprint := c.Locals("guest_fingerprint")
		if fingerprint != nil {
			userID = fingerprint.(string)
		} else {
			// Fallback: use IP as session ID
			userID = "guest_" + c.IP()
		}
	} else {
		// Authenticated user with activated referral (unlimited access)
		userEntity, ok := userFromContext.(*user.User)
		if !ok {
			status, response := shared.NewErrorResponse("ERR_401", "Invalid user data in context")
			return c.Status(status).JSON(response)
		}
		userID = userEntity.ID.String()
		isGuest = false
	}

	var req agent.ReplyRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", err.Error())
		return c.Status(status).JSON(response)
	}

	if strings.TrimSpace(req.Text) == "" {
		status, response := shared.NewErrorResponse("ERR_400", "Text cannot be empty")
		return c.Status(status).JSON(response)
	}

	// Set UserID in request body to send to LLM API
	// For logged-in users: use user_id from auth context
	// For guest users: use guest fingerprint
	req.UserID = userID

	agentResponse, err := h.agentService.Reply(c.Context(), userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "failed to make request") {
			status, response := shared.NewErrorResponse("ERR_502", "External service unavailable")
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Failed to get agent reply")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = agentResponse

	// Add guest indicator to response for transparency
	if isGuest {
		response.Meta = map[string]interface{}{
			"is_guest": true,
			"message":  "You are using guest mode. Sign in for unlimited requests.",
		}
	}

	return c.Status(status).JSON(response)
}
