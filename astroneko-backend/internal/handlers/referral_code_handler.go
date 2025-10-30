package handlers

import (
	"strconv"

	"astroneko-backend/internal/core/domain/referral_code"
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

const (
	ErrFailedToValidateReferralCode = "Failed to validate referral code"
	ErrReferralCodeIDRequired       = "Referral code ID is required"
	ErrReferralCodeNotFound         = "Referral code not found"
)

type ReferralCodeHTTPHandler struct {
	referralCodeService *services.ReferralCodeService
	validator           validator.Validator
}

func NewReferralCodeHTTPHandler(referralCodeService *services.ReferralCodeService, validator validator.Validator) *ReferralCodeHTTPHandler {
	return &ReferralCodeHTTPHandler{
		referralCodeService: referralCodeService,
		validator:           validator,
	}
}

// CreateReferralCode godoc
// @Summary Create a new general referral code
// @Description Create a new general referral code (admin only)
// @Tags referral-codes
// @Accept json
// @Produce json
// @Param referral_code body referral_code.CreateReferralCodeRequest true "Referral code data"
// @Success 201 {object} referral_code.ReferralCodeResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 409 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/referral-codes [post]
func (h *ReferralCodeHTTPHandler) CreateReferralCode(c *fiber.Ctx) error {
	var req referral_code.CreateReferralCodeRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	isValid, err := h.referralCodeService.IsValidReferralCode(c.Context(), req.ReferralCode)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", ErrFailedToValidateReferralCode)
		return c.Status(status).JSON(response)
	}

	if isValid {
		status, response := shared.NewErrorResponse("ERR_409", "Referral code already exists")
		return c.Status(status).JSON(response)
	}

	newReferralCode, err := h.referralCodeService.CreateReferralCode(c.Context(), req.ReferralCode)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to create referral code")
		return c.Status(status).JSON(response)
	}

	// Get usage count for the newly created referral code (should be 0)
	usageCount, err := h.referralCodeService.GetReferralCodeUsageCount(c.Context(), newReferralCode.ReferralCode)
	if err != nil {
		usageCount = 0 // Default to 0 if we can't get the count
	}

	status, response := shared.NewSuccessResponse("SUC_201")
	response.Data = newReferralCode.ToResponse(usageCount)
	return c.Status(status).JSON(response)
}

// GetReferralCodeByID godoc
// @Summary Get referral code by ID
// @Description Get a general referral code by its ID
// @Tags referral-codes
// @Accept json
// @Produce json
// @Param id path string true "Referral Code ID"
// @Success 200 {object} referral_code.ReferralCodeResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/referral-codes/{id} [get]
func (h *ReferralCodeHTTPHandler) GetReferralCodeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_1029", ErrReferralCodeIDRequired)
		return c.Status(status).JSON(response)
	}

	referralCode, err := h.referralCodeService.GetReferralCodeByID(c.Context(), id)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_404", ErrReferralCodeNotFound)
		return c.Status(status).JSON(response)
	}

	// Get usage count for this referral code
	usageCount, err := h.referralCodeService.GetReferralCodeUsageCount(c.Context(), referralCode.ReferralCode)
	if err != nil {
		usageCount = 0 // Default to 0 if we can't get the count
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = referralCode.ToResponse(usageCount)
	return c.Status(status).JSON(response)
}

// GetReferralCodeByCode godoc
// @Summary Get referral code by code
// @Description Get a general referral code by its code value
// @Tags referral-codes
// @Accept json
// @Produce json
// @Param code path string true "Referral Code"
// @Success 200 {object} referral_code.ReferralCodeResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/referral-codes/code/{code} [get]
func (h *ReferralCodeHTTPHandler) GetReferralCodeByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		status, response := shared.NewErrorResponse("ERR_1029", "Referral code is required")
		return c.Status(status).JSON(response)
	}

	referralCode, err := h.referralCodeService.GetReferralCodeByCode(c.Context(), code)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_404", ErrReferralCodeNotFound)
		return c.Status(status).JSON(response)
	}

	// Get usage count for this referral code
	usageCount, err := h.referralCodeService.GetReferralCodeUsageCount(c.Context(), referralCode.ReferralCode)
	if err != nil {
		usageCount = 0 // Default to 0 if we can't get the count
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = referralCode.ToResponse(usageCount)
	return c.Status(status).JSON(response)
}

// UpdateReferralCode godoc
// @Summary Update referral code
// @Description Update a general referral code
// @Tags referral-codes
// @Accept json
// @Produce json
// @Param id path string true "Referral Code ID"
// @Param referral_code body referral_code.UpdateReferralCodeRequest true "Updated referral code data"
// @Success 200 {object} referral_code.ReferralCodeResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Failure 409 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/referral-codes/{id} [put]
func (h *ReferralCodeHTTPHandler) UpdateReferralCode(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_1029", ErrReferralCodeIDRequired)
		return c.Status(status).JSON(response)
	}

	var req referral_code.UpdateReferralCodeRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	existingReferralCode, err := h.referralCodeService.GetReferralCodeByID(c.Context(), id)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_404", ErrReferralCodeNotFound)
		return c.Status(status).JSON(response)
	}

	if existingReferralCode.ReferralCode != req.ReferralCode {
		isValid, err := h.referralCodeService.IsValidReferralCode(c.Context(), req.ReferralCode)
		if err != nil {
			status, response := shared.NewErrorResponse("ERR_500", ErrFailedToValidateReferralCode)
			return c.Status(status).JSON(response)
		}

		if isValid {
			status, response := shared.NewErrorResponse("ERR_409", "Referral code already exists")
			return c.Status(status).JSON(response)
		}
	}

	existingReferralCode.ReferralCode = req.ReferralCode
	updatedReferralCode, err := h.referralCodeService.UpdateReferralCode(c.Context(), existingReferralCode)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to update referral code")
		return c.Status(status).JSON(response)
	}

	// Get usage count for the updated referral code
	usageCount, err := h.referralCodeService.GetReferralCodeUsageCount(c.Context(), updatedReferralCode.ReferralCode)
	if err != nil {
		usageCount = 0 // Default to 0 if we can't get the count
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = updatedReferralCode.ToResponse(usageCount)
	return c.Status(status).JSON(response)
}

// DeleteReferralCode godoc
// @Summary Delete referral code
// @Description Delete a general referral code
// @Tags referral-codes
// @Accept json
// @Produce json
// @Param id path string true "Referral Code ID"
// @Success 204 {object} shared.ResponseBody
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/referral-codes/{id} [delete]
func (h *ReferralCodeHTTPHandler) DeleteReferralCode(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_1029", ErrReferralCodeIDRequired)
		return c.Status(status).JSON(response)
	}

	if err := h.referralCodeService.DeleteReferralCode(c.Context(), id); err != nil {
		status, response := shared.NewErrorResponse("ERR_404", ErrReferralCodeNotFound)
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_204")
	return c.Status(status).JSON(response)
}

// ListReferralCodes godoc
// @Summary List all referral codes
// @Description Get a paginated list of all general referral codes
// @Tags referral-codes
// @Accept json
// @Produce json
// @Param limit query int false "Number of items to return (default: 10, max: 100)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Success 200 {object} referral_code.ListReferralCodesResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/referral-codes [get]
func (h *ReferralCodeHTTPHandler) ListReferralCodes(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	referralCodes, total, err := h.referralCodeService.ListReferralCodes(c.Context(), limit, offset)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to list referral codes")
		return c.Status(status).JSON(response)
	}

	var codeResponses []referral_code.ReferralCodeResponse
	for _, code := range referralCodes {
		// Get usage count for each referral code
		usageCount, err := h.referralCodeService.GetReferralCodeUsageCount(c.Context(), code.ReferralCode)
		if err != nil {
			usageCount = 0 // Default to 0 if we can't get the count
		}
		codeResponses = append(codeResponses, *code.ToResponse(usageCount))
	}

	responseData := referral_code.ListReferralCodesResponse{
		Codes:  codeResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = responseData
	return c.Status(status).JSON(response)
}

// ValidateReferralCode godoc
// @Summary Validate referral code
// @Description Check if a referral code is valid
// @Tags referral-codes
// @Accept json
// @Produce json
// @Param code path string true "Referral Code"
// @Success 200 {object} shared.ResponseBody
// @Failure 400 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/referral-codes/validate/{code} [get]
func (h *ReferralCodeHTTPHandler) ValidateReferralCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		status, response := shared.NewErrorResponse("ERR_1029", "Referral code is required")
		return c.Status(status).JSON(response)
	}

	isValid, err := h.referralCodeService.IsValidReferralCode(c.Context(), code)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", ErrFailedToValidateReferralCode)
		return c.Status(status).JSON(response)
	}

	if !isValid {
		status, response := shared.NewErrorResponse("ERR_404", "Referral code is invalid")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = fiber.Map{"message": "Referral code is valid", "valid": true}
	return c.Status(status).JSON(response)
}
