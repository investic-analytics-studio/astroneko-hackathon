package auth

import "github.com/google/uuid"

type TokenPayloadRequest struct {
	UserID *uuid.UUID `json:"user_id,omitempty"`
	Email  *string    `json:"email,omitempty"`
	RoleID *int64     `json:"role_id,omitempty"`
	Role   *string    `json:"role,omitempty"`
}
