package handlers

import (
	"strings"

	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/waiting_list"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type WaitingListHTTPHandler struct {
	waitingListService *services.WaitingListService
	validator          validator.Validator
}

func NewWaitingListHTTPHandler(waitingListService *services.WaitingListService, validator validator.Validator) *WaitingListHTTPHandler {
	return &WaitingListHTTPHandler{
		waitingListService: waitingListService,
		validator:          validator,
	}
}

// JoinWaitingList godoc
// @Summary Join waiting list
// @Description Add user to the waiting list by email
// @Tags waiting_list
// @Accept json
// @Produce json
// @Param request body waiting_list.JoinWaitingListRequest true "Join waiting list request"
// @Success 201 {object} waiting_list.JoinWaitingListResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 409 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/waiting-list/join [post]
func (h *WaitingListHTTPHandler) JoinWaitingList(c *fiber.Ctx) error {
	var req waiting_list.JoinWaitingListRequest

	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", err.Error())
		return c.Status(status).JSON(response)
	}

	waitingListUser, err := h.waitingListService.JoinWaitingList(c.Context(), req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "already exists in waiting list") {
			status, response := shared.NewErrorResponse("ERR_1034")
			return c.Status(status).JSON(response)
		}
		if strings.Contains(err.Error(), "failed to add user to waiting list") {
			status, response := shared.NewErrorResponse("ERR_1035")
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Failed to join waiting list")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_201")
	response.Data = waitingListUser.ToResponse()
	return c.Status(status).JSON(response)
}

// IsInWaitingListByEmail godoc
// @Summary Check if email is in waiting list
// @Description Check if an email exists in the waiting list
// @Tags waiting_list
// @Accept json
// @Produce json
// @Param request body waiting_list.CheckWaitingListRequest true "Check waiting list request"
// @Success 200 {object} waiting_list.IsInWaitingListResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/waiting-list/check [post]
func (h *WaitingListHTTPHandler) IsInWaitingListByEmail(c *fiber.Ctx) error {
	var req waiting_list.CheckWaitingListRequest

	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_400", err.Error())
		return c.Status(status).JSON(response)
	}

	isInWaitingList, err := h.waitingListService.IsInWaitingListByEmail(c.Context(), req.Email)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_500", "Failed to check waiting list status")
		return c.Status(status).JSON(response)
	}

	response := waiting_list.NewIsInWaitingListResponse(isInWaitingList)
	status, resp := shared.NewSuccessResponse("SUC_200")
	resp.Data = response
	return c.Status(status).JSON(resp)
}
