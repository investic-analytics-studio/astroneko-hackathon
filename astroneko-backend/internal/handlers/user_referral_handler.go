package handlers

import (
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserReferralHandler struct {
	userService         *services.UserService
	referralCodeService *services.ReferralCodeService
	validator           validator.Validator
}

func NewUserReferralHandler(userService *services.UserService, referralCodeService *services.ReferralCodeService, validator validator.Validator) *UserReferralHandler {
	return &UserReferralHandler{
		userService:         userService,
		referralCodeService: referralCodeService,
		validator:           validator,
	}
}

// ActivateReferral godoc
// @Summary Activate referral code
// @Description Activate a referral code for a user
// @Tags users/referral
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param referral body user.ActivateReferralRequest true "Referral data"
// @Success 200 {object} user.ActivateReferralResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Router /v1/api/users/{id}/referral/activate [post]
func (h *UserReferralHandler) ActivateReferral(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_1029", ErrUserIDRequired)
		return c.Status(status).JSON(response)
	}

	var req user.ActivateReferralRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", ErrInvalidRequestBody)
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	// Use ReferralCodeService to activate the referral code
	referralResp, err := h.referralCodeService.ActivateReferralCode(c.Context(), uuid.MustParse(id), req.ReferralCode)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to activate referral")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = referralResp
	return c.Status(status).JSON(response)
}

// GetUserReferralCodes godoc
// @Summary Get user referral codes
// @Description Get all referral codes for a user
// @Tags users/referral
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} referral_code.GetUserReferralCodesResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Router /v1/api/users/{id}/referral/codes [get]
func (h *UserReferralHandler) GetUserReferralCodes(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_1029", ErrUserIDRequired)
		return c.Status(status).JSON(response)
	}

	referralCodes, err := h.referralCodeService.GetOrGenerateUserReferralCodes(c.Context(), uuid.MustParse(id))
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to get referral codes")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = referralCodes
	return c.Status(status).JSON(response)
}
