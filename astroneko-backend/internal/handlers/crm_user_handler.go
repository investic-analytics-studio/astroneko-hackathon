package handlers

import (
	"strings"

	"astroneko-backend/internal/core/domain/crm_user"
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type CRMUserHTTPHandler struct {
	crmUserService *services.CRMUserService
	validator      validator.Validator
}

func NewCRMUserHTTPHandler(crmUserService *services.CRMUserService, validator validator.Validator) *CRMUserHTTPHandler {
	return &CRMUserHTTPHandler{
		crmUserService: crmUserService,
		validator:      validator,
	}
}

// CreateCRMUser godoc
// @Summary Create a new CRM user
// @Description Create a new CRM user with username and password
// @Tags crm-auth
// @Accept json
// @Produce json
// @Param user body crm_user.CreateCRMUserRequest true "CRM user data"
// @Success 201 {object} crm_user.CRMUserResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 409 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/crm/users [post]
func (h *CRMUserHTTPHandler) CreateCRMUser(c *fiber.Ctx) error {
	var req crm_user.CreateCRMUserRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	newUser, err := h.crmUserService.CreateUser(c.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "username already exists") {
			status, response := shared.NewErrorResponse("ERR_409", "Username already exists")
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Failed to create CRM user")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_201")
	response.Data = newUser.ToResponse()
	return c.Status(status).JSON(response)
}

// CRMLogin godoc
// @Summary Login with CRM credentials
// @Description Authenticate CRM user with username and password, returns JWT token
// @Tags crm-auth
// @Accept json
// @Produce json
// @Param login body crm_user.CRMLoginRequest true "CRM login credentials"
// @Success 200 {object} crm_user.CRMLoginResponse
// @Failure 400 {object} shared.ResponseBody
// @Failure 401 {object} shared.ResponseBody
// @Failure 500 {object} shared.ResponseBody
// @Router /v1/api/crm/auth/login [post]
func (h *CRMUserHTTPHandler) CRMLogin(c *fiber.Ctx) error {
	var req crm_user.CRMLoginRequest
	if err := c.BodyParser(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", "Invalid request body")
		return c.Status(status).JSON(response)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		status, response := shared.NewErrorResponse("ERR_1029", err.Error())
		return c.Status(status).JSON(response)
	}

	loginResp, err := h.crmUserService.Login(c.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			status, response := shared.NewErrorResponse("ERR_401", "Invalid username or password")
			return c.Status(status).JSON(response)
		}
		status, response := shared.NewErrorResponse("ERR_500", "Login failed")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = loginResp
	return c.Status(status).JSON(response)
}

// GetCRMMe godoc
// @Summary Get current CRM user information
// @Description Get current authenticated CRM user information from JWT token
// @Tags crm-auth
// @Accept json
// @Produce json
// @Success 200 {object} crm_user.CRMUserResponse
// @Failure 401 {object} shared.ResponseBody
// @Failure 404 {object} shared.ResponseBody
// @Security BearerAuth
// @Router /v1/api/crm/auth/me [get]
func (h *CRMUserHTTPHandler) GetCRMMe(c *fiber.Ctx) error {
	userFromContext := c.Locals("crm_user")
	if userFromContext == nil {
		status, response := shared.NewErrorResponse("ERR_401", "CRM user not found in context")
		return c.Status(status).JSON(response)
	}

	userEntity, ok := userFromContext.(*crm_user.CRMUser)
	if !ok {
		status, response := shared.NewErrorResponse("ERR_401", "Invalid CRM user data in context")
		return c.Status(status).JSON(response)
	}

	status, response := shared.NewSuccessResponse("SUC_200")
	response.Data = userEntity.ToResponse()
	return c.Status(status).JSON(response)
}
