package handlers

import (
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/errors"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UserAuthHandler struct {
	userService  *services.UserService
	validator    validator.Validator
	errorHandler *errors.Handler
}

func NewUserAuthHandler(userService *services.UserService, validator validator.Validator) *UserAuthHandler {
	return &UserAuthHandler{
		userService:  userService,
		validator:    validator,
		errorHandler: errors.NewHandler(),
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with Firebase authentication
// @Tags users/auth
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "User data"
// @Success 201 {object} user.CreateUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/users/auth/register [post]
func (h *UserAuthHandler) CreateUser(c *fiber.Ctx) error {
	var req user.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.HandleRequestBodyError(c, err)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		return h.errorHandler.HandleValidationError(c, err)
	}

	newUser, err := h.userService.CreateUser(c.Context(), &req)
	if err != nil {
		return h.errorHandler.HandleServiceError(c, err)
	}

	return h.errorHandler.HandleCreated(c, newUser.ToResponse())
}

// GoogleAuth godoc
// @Summary Authenticate with Google
// @Description Authenticate user with Google OAuth
// @Tags users/auth
// @Accept json
// @Produce json
// @Param auth body user.GoogleLoginRequest true "Google auth data"
// @Success 200 {object} user.AuthResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Router /v1/api/users/auth/google [post]
func (h *UserAuthHandler) GoogleAuth(c *fiber.Ctx) error {
	var req user.GoogleLoginRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", ErrInvalidRequestBody)
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	authResp, err := h.userService.GoogleAuth(c.Context(), &req)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_401", "Authentication failed")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = authResp
	return c.Status(status).JSON(response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags users/auth
// @Accept json
// @Produce json
// @Param refresh body user.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} user.RefreshTokenResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Router /v1/api/users/auth/refresh [post]
func (h *UserAuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req user.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", ErrInvalidRequestBody)
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	tokenResp, err := h.userService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_401", "Invalid refresh token")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = tokenResp
	return c.Status(status).JSON(response)
}
