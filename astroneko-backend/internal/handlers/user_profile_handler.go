package handlers

import (
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UserProfileHandler struct {
	userService *services.UserService
	validator   validator.Validator
}

func NewUserProfileHandler(userService *services.UserService, validator validator.Validator) *UserProfileHandler {
	return &UserProfileHandler{
		userService: userService,
		validator:   validator,
	}
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags users/profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user.GetUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Router /v1/api/users/profile/{id} [get]
func (h *UserProfileHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_1029", ErrUserIDRequired)
		return c.Status(status).JSON(response)
	}

	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_404", ErrUserNotFound)
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = user.ToResponse()
	return c.Status(status).JSON(response)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users/profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body user.UpdateUserRequest true "User data"
// @Success 200 {object} user.UpdateUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Router /v1/api/users/profile/{id} [put]
func (h *UserProfileHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_1029", ErrUserIDRequired)
		return c.Status(status).JSON(response)
	}

	var req user.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", ErrInvalidRequestBody)
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	updatedUser, err := h.userService.UpdateUser(c.Context(), id, &req)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to update user")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = updatedUser.ToResponse()
	return c.Status(status).JSON(response)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users/profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user.DeleteUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Router /v1/api/users/profile/{id} [delete]
func (h *UserProfileHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_1029", ErrUserIDRequired)
		return c.Status(status).JSON(response)
	}

	err := h.userService.DeleteUser(c.Context(), id)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to delete user")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	return c.Status(status).JSON(response)
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags users/profile
// @Accept json
// @Produce json
// @Success 200 {object} user.GetUserResponse
// @Failure 401 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/users/profile/me [get]
func (h *UserProfileHandler) GetCurrentUser(c *fiber.Ctx) error {
	userData, ok := c.Locals("user").(*user.User)
	if !ok {
		status, response := shared.NewErrorResponse("ERR_401", ErrUserNotFoundInContext)
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = userData.ToResponse()
	return c.Status(status).JSON(response)
}
