package handlers

import (
	"strings"

	"astroneko-backend/internal/core/domain/referral_code"
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHTTPHandler struct {
	userService         *services.UserService
	referralCodeService *services.ReferralCodeService
	validator           validator.Validator
}

func NewUserHTTPHandler(userService *services.UserService, referralCodeService *services.ReferralCodeService, validator validator.Validator) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService:         userService,
		referralCodeService: referralCodeService,
		validator:           validator,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with Firebase authentication
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "User data"
// @Success 201 {object} user.CreateUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/users [post]
func (h *UserHTTPHandler) CreateUser(c *fiber.Ctx) error {
	var req user.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", ErrInvalidRequestBody)
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	newUser, err := h.userService.CreateUser(c.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			status, response := shared.NewErrorResponse("ERR_409", err.Error())
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Failed to create user")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_201")
	response.Data = newUser.ToResponse()
	return c.Status(status).JSON(response)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user.GetUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Router /v1/api/users/{id} [get]
func (h *UserHTTPHandler) GetUserByID(c *fiber.Ctx) error {
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
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body user.UpdateUserRequest true "Updated user data"
// @Success 200 {object} user.UpdateUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Router /v1/api/users/{id} [put]
func (h *UserHTTPHandler) UpdateUser(c *fiber.Ctx) error {
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

	updatedUser, err := h.userService.UpdateUser(c.Context(), id, &req)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_404", ErrUserNotFound)
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = updatedUser.ToResponse()
	return c.Status(status).JSON(response)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 {object} user.DeleteUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Router /v1/api/users/{id} [delete]
func (h *UserHTTPHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_1029", ErrUserIDRequired)
		return c.Status(status).JSON(response)
	}

	if err := h.userService.DeleteUser(c.Context(), id); err != nil {
		status, response := shared.NewErrorResponse("ERR_404", ErrUserNotFound)
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_204")
	return c.Status(status).JSON(response)
}

// RefreshToken godoc
// @Summary Refresh Firebase token
// @Description Refresh an expired Firebase token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body user.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} user.RefreshTokenResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/auth/refresh [post]
func (h *UserHTTPHandler) RefreshToken(c *fiber.Ctx) error {
	var req user.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", ErrInvalidRequestBody)
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	authResp, err := h.userService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "failed to refresh token") {
			status, response := shared.NewErrorResponse("ERR_401", "Invalid or expired refresh token")
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Failed to refresh token")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = authResp
	return c.Status(status).JSON(response)
}

// Login godoc
// @Summary Login with email and password
// @Description Authenticate user with email and password, check both database and Firebase
// @Tags auth
// @Accept json
// @Produce json
// @Param login body user.LoginRequest true "Login credentials"
// @Success 200 {object} user.RefreshTokenResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/auth/login [post]
func (h *UserHTTPHandler) Login(c *fiber.Ctx) error {
	var req user.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", ErrInvalidRequestBody)
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	authResp, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			status, response := shared.NewErrorResponse("ERR_401", "Invalid email or password")
			return c.Status(status).JSON(response)
		}
		if strings.Contains(err.Error(), "authentication failed") {
			status, response := shared.NewErrorResponse("ERR_401", "Invalid email or password")
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Login failed")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = authResp
	return c.Status(status).JSON(response)
}

// AuthenticateWithFirebase godoc
// @Summary Authenticate with Firebase token
// @Description Verify Firebase ID token and return user info (simplified endpoint)
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {firebase_id_token}"
// @Success 200 {object} user.GetUserResponse
// @Failure 401 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Router /v1/api/auth/firebase [post]
func (h *UserHTTPHandler) AuthenticateWithFirebase(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		status, response := shared.NewErrorResponse("ERR_401", "Authorization header required")
		return c.Status(status).JSON(response)
	}

	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := h.userService.VerifyFirebaseToken(c.Context(), idToken)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_401", "Invalid Firebase token")
		return c.Status(status).JSON(response)
	}

	// Get existing user (don't auto-create here, use dedicated endpoints for that)
	existingUser, err := h.userService.GetUserByFirebaseUID(c.Context(), token.UID)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_404", "User not found in database")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = existingUser.ToResponse()
	return c.Status(status).JSON(response)
}

// GoogleAuth godoc
// @Summary Authenticate with Google OAuth
// @Description Register new user or login existing user with Google ID token (Firebase user creation handled in frontend)
// @Tags auth
// @Accept json
// @Produce json
// @Param google_auth body user.GoogleLoginRequest true "Google authentication tokens"
// @Success 200 {object} user.RefreshTokenResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/auth/google [post]
func (h *UserHTTPHandler) GoogleAuth(c *fiber.Ctx) error {
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
		status, response := shared.NewErrorResponse("ERR_401", "Google authentication failed")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = authResp
	return c.Status(status).JSON(response)
}

// GetMe godoc
// @Summary Get current user information
// @Description Get current authenticated user information from JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} user.GetUserResponse
// @Failure 401 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/auth/me [get]
func (h *UserHTTPHandler) GetMe(c *fiber.Ctx) error {
	userFromContext := c.Locals("user")
	if userFromContext == nil {
		status, response := shared.NewErrorResponse("ERR_401", ErrUserNotFoundInContext)
		return c.Status(status).JSON(response)
	}

	userEntity, ok := userFromContext.(*user.User)
	if !ok {
		status, response := shared.NewErrorResponse("ERR_401", ErrInvalidUserDataInContext)
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = userEntity.ToResponse()
	return c.Status(status).JSON(response)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} shared.ResponseBody
// @Router /v1/api/auth/logout [post]
func (h *UserHTTPHandler) Logout(c *fiber.Ctx) error {
	status, response := shared.NewSuccessResponse("SUC_200")
	return c.Status(status).JSON(response)
}

// GetUserReferralCodes godoc
// @Summary Get user referral codes
// @Description Get or generate user referral codes (5 codes with 8 characters each). Requires user to have is_activated_referral = true
// @Tags referral
// @Accept json
// @Produce json
// @Success 200 {object} referral_code.GetUserReferralCodesResponse
// @Failure 401 {object} shared.ResponseBody
// @Failure 403 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/auth/referral/codes [get]
func (h *UserHTTPHandler) GetUserReferralCodes(c *fiber.Ctx) error {
	userFromContext := c.Locals("user")
	if userFromContext == nil {
		status, response := shared.NewErrorResponse("ERR_401", ErrUserNotFoundInContext)
		return c.Status(status).JSON(response)
	}

	userEntity, ok := userFromContext.(*user.User)
	if !ok {
		status, response := shared.NewErrorResponse("ERR_401", ErrInvalidUserDataInContext)
		return c.Status(status).JSON(response)
	}

	userID, err := uuid.Parse(userEntity.ID.String())
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", "Invalid user ID")
		return c.Status(status).JSON(response)
	}

	userReferralCodes, err := h.referralCodeService.GetOrGenerateUserReferralCodes(c.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "user has not activated referral feature") {
			status, response := shared.NewErrorResponse("ERR_403", "Referral feature not activated")
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Failed to get user referral codes")
		return c.Status(status).JSON(response)
	}

	// Convert to response format
	var codeResponses []referral_code.UserReferralCodeResponse
	for _, code := range userReferralCodes {
		codeResponses = append(codeResponses, *code.ToResponse())
	}

	responseData := referral_code.GetUserReferralCodesResponse{
		Codes: codeResponses,
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = responseData
	return c.Status(status).JSON(response)
}

// ActivateReferral godoc
// @Summary Activate referral code for authenticated user
// @Description Activate referral code with improved validation and logging. Sets is_activated_referral = true in user table. Returns error if user has already activated a referral code.
// @Tags referral
// @Accept json
// @Produce json
// @Param referral body referral_code.ActivateReferralRequest true "Referral code"
// @Success 200 {object} referral_code.ActivateReferralResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/auth/referral/activate [post]
func (h *UserHTTPHandler) ActivateReferral(c *fiber.Ctx) error {
	userFromContext := c.Locals("user")
	if userFromContext == nil {
		status, response := shared.NewErrorResponse("ERR_401", ErrUserNotFoundInContext)
		return c.Status(status).JSON(response)
	}

	userEntity, ok := userFromContext.(*user.User)
	if !ok {
		status, response := shared.NewErrorResponse("ERR_401", ErrInvalidUserDataInContext)
		return c.Status(status).JSON(response)
	}

	var req referral_code.ActivateReferralRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", ErrInvalidRequestBody)
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	userID, err := uuid.Parse(userEntity.ID.String())
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", "Invalid user ID")
		return c.Status(status).JSON(response)
	}

	activationResponse, err := h.referralCodeService.ActivateReferralCode(c.Context(), userID, req.ReferralCode)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to activate referral code")
		return c.Status(status).JSON(response)
	}

	if !activationResponse.Success {
		// Use specific error codes for referral activation failures
		if strings.Contains(activationResponse.Message, "already activated") {
			status, response := shared.NewErrorResponse("ERR_1033")
			return c.Status(status).JSON(response)
		}
		// Default to invalid referral code for other failures
		status, response := shared.NewErrorResponse("ERR_1032")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = activationResponse
	return c.Status(status).JSON(response)
}

// GetTotalUsers godoc
// @Summary Get total number of users
// @Description Get the total count of users in the system (CRM access required)
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} user.GetTotalUsersResponse
// @Failure 401 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/crm/users/total [get]
func (h *UserHTTPHandler) GetTotalUsers(c *fiber.Ctx) error {
	totalUsers, err := h.userService.GetTotalUsers(c.Context())
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to get total users")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = &user.GetTotalUsersResponse{
		TotalUsers: totalUsers,
	}
	return c.Status(status).JSON(response)
}
