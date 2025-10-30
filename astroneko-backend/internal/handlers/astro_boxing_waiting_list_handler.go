package handlers

import (
	"strings"

	"astroneko-backend/internal/core/domain/astro_boxing_waiting_list"
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AstroBoxingWaitingListHTTPHandler struct {
	astroBoxingWaitingListService *services.AstroBoxingWaitingListService
	validator                     validator.Validator
}

func NewAstroBoxingWaitingListHTTPHandler(astroBoxingWaitingListService *services.AstroBoxingWaitingListService, validator validator.Validator) *AstroBoxingWaitingListHTTPHandler {
	return &AstroBoxingWaitingListHTTPHandler{
		astroBoxingWaitingListService: astroBoxingWaitingListService,
		validator:                     validator,
	}
}

// JoinAstroBoxingWaitingList godoc
// @Summary Join astro boxing waiting list
// @Description Add user to the astro boxing waiting list
// @Tags astro_boxing_waiting_list
// @Accept json
// @Produce json
// @Param request body astro_boxing_waiting_list.JoinAstroBoxingWaitingListRequest true "Join request"
// @Success 201 {object} astro_boxing_waiting_list.JoinAstroBoxingWaitingListResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 409 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/astro-boxing-waiting-list/join [post]
func (h *AstroBoxingWaitingListHTTPHandler) JoinAstroBoxingWaitingList(c *fiber.Ctx) error {
	var req astro_boxing_waiting_list.JoinAstroBoxingWaitingListRequest

	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", err.Error())
		return c.Status(status).JSON(response)
	}

	user, err := h.astroBoxingWaitingListService.JoinWaitingList(c.Context(), req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "already exists in waiting list") {
			status, response := shared.NewErrorResponse("ERR_1034")
			return c.Status(status).JSON(response)
		}
		if strings.Contains(err.Error(), "failed to add user to waiting list") {
			status, response := shared.NewErrorResponse("ERR_1035")
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Failed to join astro boxing waiting list")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_201")
	response.Data = user.ToResponse()
	return c.Status(status).JSON(response)
}

// IsInAstroBoxingWaitingListByEmail godoc
// @Summary Check if email is in astro boxing waiting list
// @Description Check if an email exists in the astro boxing waiting list
// @Tags astro_boxing_waiting_list
// @Accept json
// @Produce json
// @Param email query string true "Email to check"
// @Success 200 {object} astro_boxing_waiting_list.IsInAstroBoxingWaitingListResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/astro-boxing-waiting-list/check [get]
func (h *AstroBoxingWaitingListHTTPHandler) IsInAstroBoxingWaitingListByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		status, response := shared.NewErrorResponse("ERR_400", "Email query parameter is required")
		return c.Status(status).JSON(response)
	}

	type emailValidation struct {
		Email string `validate:"required,email"`
	}
	emailStruct := emailValidation{Email: email}
	if err := h.validator.ValidateStruct(&emailStruct); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", "Invalid email format")
		return c.Status(status).JSON(response)
	}

	isInWaitingList, err := h.astroBoxingWaitingListService.IsInWaitingListByEmail(c.Context(), email)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to check astro boxing waiting list status")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = astro_boxing_waiting_list.IsInAstroBoxingWaitingListResponse{
		IsInWaitingList: isInWaitingList,
	}
	return c.Status(status).JSON(response)
}

// GetAstroBoxingWaitingListUsers godoc
// @Summary Get astro boxing waiting list users
// @Description Get paginated list of astro boxing waiting list users
// @Tags astro_boxing_waiting_list
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} astro_boxing_waiting_list.AstroBoxingWaitingListUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/astro-boxing-waiting-list/users [get]
func (h *AstroBoxingWaitingListHTTPHandler) GetAstroBoxingWaitingListUsers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	users, total, err := h.astroBoxingWaitingListService.GetWaitingListUsers(c.Context(), limit, offset)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to get astro boxing waiting list users")
		return c.Status(status).JSON(response)
	}

	var responses []*astro_boxing_waiting_list.AstroBoxingWaitingListUserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = map[string]interface{}{
		"users":  responses,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
	return c.Status(status).JSON(response)
}

// DeleteAstroBoxingWaitingListUser godoc
// @Summary Delete astro boxing waiting list user
// @Description Delete a user from astro boxing waiting list by ID
// @Tags astro_boxing_waiting_list
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} shared.ResponseBody
// @Failure 400 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/astro-boxing-waiting-list/users/{id} [delete]
func (h *AstroBoxingWaitingListHTTPHandler) DeleteAstroBoxingWaitingListUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		status, response := shared.NewErrorResponse("ERR_400", "User ID is required")
		return c.Status(status).JSON(response)
	}

	err := h.astroBoxingWaitingListService.DeleteUser(c.Context(), id)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to delete user from astro boxing waiting list")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	return c.Status(status).JSON(response)
}
