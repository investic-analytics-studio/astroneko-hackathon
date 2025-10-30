package handlers

import (
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user_limit"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UserLimitHTTPHandler struct {
	userLimitService *services.UserLimitService
	validator        validator.Validator
}

func NewUserLimitHTTPHandler(userLimitService *services.UserLimitService, validator validator.Validator) *UserLimitHTTPHandler {
	return &UserLimitHTTPHandler{
		userLimitService: userLimitService,
		validator:        validator,
	}
}

// GetUserLimit godoc
// @Summary Get user limit
// @Description Get the current user limit configuration
// @Tags user-limit
// @Accept json
// @Produce json
// @Success 200 {object} user_limit.GetUserLimitResponse
// @Failure 404 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/user-limit [get]
func (h *UserLimitHTTPHandler) GetUserLimit(c *fiber.Ctx) error {
	userLimit, err := h.userLimitService.GetUserLimit(c.Context())
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to get user limit")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = userLimit.ToResponse()
	return c.Status(status).JSON(response)
}

// UpdateUserLimit godoc
// @Summary Update user limit
// @Description Update the user limit configuration (CRM access required)
// @Tags user-limit
// @Accept json
// @Produce json
// @Param user_limit body user_limit.UpdateUserLimitRequest true "Updated user limit data"
// @Success 200 {object} user_limit.UpdateUserLimitResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/crm/user-limit [put]
func (h *UserLimitHTTPHandler) UpdateUserLimit(c *fiber.Ctx) error {
	var req user_limit.UpdateUserLimitRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	updatedUserLimit, err := h.userLimitService.UpdateUserLimit(c.Context(), &req)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to update user limit")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = updatedUserLimit.ToResponse()
	return c.Status(status).JSON(response)
}

// IsUserOverLimitUsed godoc
// @Summary Check if user limit is exceeded
// @Description Check if the total number of users exceeds the configured limit
// @Tags user-limit
// @Accept json
// @Produce json
// @Success 200 {object} user_limit.IsUserOverLimitUsedResponse
// @Failure 401 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/user-limit/check [get]
func (h *UserLimitHTTPHandler) IsUserOverLimitUsed(c *fiber.Ctx) error {
	isOverLimit, err := h.userLimitService.IsUserOverLimitUsed(c.Context())
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to check user limit")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = &user_limit.IsUserOverLimitUsedResponse{
		IsOverLimit: isOverLimit,
	}
	return c.Status(status).JSON(response)
}
