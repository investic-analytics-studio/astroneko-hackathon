package handlers

import (
	"strings"

	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HistoryHTTPHandler struct {
	historyService *services.HistoryService
}

// NewHistoryHTTPHandler creates a new history HTTP handler
func NewHistoryHTTPHandler(historyService *services.HistoryService) *HistoryHTTPHandler {
	return &HistoryHTTPHandler{
		historyService: historyService,
	}
}

// GetUserSessions godoc
// @Summary Get user's conversation sessions
// @Description Retrieve all conversation sessions for the authenticated user. Supports sorting by created_at or updated_at in ascending or descending order, and searching by history_name.
// @Tags history
// @Accept json
// @Produce json
// @Param sort_by query string false "Sort field: created_at or updated_at (default: updated_at)"
// @Param sort_order query string false "Sort order: asc or desc (default: desc)"
// @Param search query string false "Search query to filter sessions by history_name (partial match)"
// @Success 200 {object} history.GetSessionsResponse
// @Failure 401 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/history/sessions [get]
func (h *HistoryHTTPHandler) GetUserSessions(c *fiber.Ctx) error {
	// Extract user from context (set by auth middleware)
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

	// Get query parameters for sorting and search
	sortBy := c.Query("sort_by", "updated_at") // default: updated_at
	sortOrder := c.Query("sort_order", "desc") // default: desc
	searchQuery := c.Query("search", "") // default: empty (no search filter)

	// Get sessions for the user
	sessionsResponse, err := h.historyService.GetUserSessions(c.Context(), userEntity.ID, sortBy, sortOrder, searchQuery)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to retrieve sessions")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = sessionsResponse
	return c.Status(status).JSON(response)
}

// GetSessionMessages godoc
// @Summary Get messages for a session
// @Description Retrieve all messages for a specific session. Validates that the session belongs to the authenticated user. Supports sorting by created_at in ascending (chronological) or descending order.
// @Tags history
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID (UUID)"
// @Param sort_order query string false "Sort order: asc or desc (default: asc for chronological)"
// @Success 200 {object} history.GetMessagesResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Failure 403 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/history/sessions/{session_id}/messages [get]
func (h *HistoryHTTPHandler) GetSessionMessages(c *fiber.Ctx) error {
	// Extract user from context (set by auth middleware)
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

	// Get session_id from URL parameter
	sessionIDParam := c.Params("session_id")
	if sessionIDParam == "" {
		status, response := shared.NewErrorResponse("ERR_400", "Session ID is required")
		return c.Status(status).JSON(response)
	}

	// Parse session_id as UUID
	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_400", "Invalid session ID format")
		return c.Status(status).JSON(response)
	}

	// Get query parameter for sorting
	sortOrder := c.Query("sort_order", "asc") // default: asc (chronological)

	// Get messages for the session (with ownership validation)
	messagesResponse, err := h.historyService.GetSessionMessages(c.Context(), userEntity.ID, sessionID, sortOrder)
	if err != nil {
		// Check for specific error types
		if strings.Contains(err.Error(), "access denied") || strings.Contains(err.Error(), "not found") {
			status, response := shared.NewErrorResponse("ERR_403", "Session not found or access denied")
			return c.Status(status).JSON(response)
		}

		status, response := shared.NewErrorResponse("ERR_500", "Failed to retrieve messages")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = messagesResponse
	return c.Status(status).JSON(response)
}

// DeleteSession godoc
// @Summary Delete a conversation session
// @Description Soft delete a conversation session by setting deleted_at timestamp. Validates that the session belongs to the authenticated user.
// @Tags history
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID (UUID)"
// @Success 200 {object} shared.ResponseBody
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Failure 403 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/history/sessions/{session_id} [delete]
func (h *HistoryHTTPHandler) DeleteSession(c *fiber.Ctx) error {
	// Extract user from context (set by auth middleware)
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

	// Get session_id from URL parameter
	sessionIDParam := c.Params("session_id")
	if sessionIDParam == "" {
		status, response := shared.NewErrorResponse("ERR_400", "Session ID is required")
		return c.Status(status).JSON(response)
	}

	// Parse session_id as UUID
	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_400", "Invalid session ID format")
		return c.Status(status).JSON(response)
	}

	// Delete the session (with ownership validation)
	err = h.historyService.DeleteSession(c.Context(), userEntity.ID, sessionID)
	if err != nil {
		// Check for specific error types
		if strings.Contains(err.Error(), "access denied") || strings.Contains(err.Error(), "not found") {
			status, response := shared.NewErrorResponse("ERR_403", "Session not found or access denied")
			return c.Status(status).JSON(response)
		}

		status, response := shared.NewErrorResponse("ERR_500", "Failed to delete session")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200", "Session deleted successfully")
	return c.Status(status).JSON(response)
}
