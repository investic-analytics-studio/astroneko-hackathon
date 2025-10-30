package shared

import "time"

// BaseResponse provides common response fields
type BaseResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EmailResponse provides common email-based response
type EmailResponse struct {
	BaseResponse
	Email string `json:"email"`
}

// StatusResponse provides common status-based response
type StatusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// PaginationResponse provides pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ListResponse provides a generic list response format
type ListResponse[T any] struct {
	Data       []T                `json:"data"`
	Pagination PaginationResponse `json:"pagination,omitempty"`
}

// ActionResponse provides a generic action response
type ActionResponse struct {
	ResponseBody
	Action string `json:"action,omitempty"`
}

// NewEmailResponse creates a new email response from a model
func NewEmailResponse(id string, email string, createdAt, updatedAt time.Time) *EmailResponse {
	return &EmailResponse{
		BaseResponse: BaseResponse{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Email: email,
	}
}

// NewStatusResponse creates a new status response
func NewStatusResponse(success bool, message string) *StatusResponse {
	return &StatusResponse{
		Success: success,
		Message: message,
	}
}
