package crm_user

import (
	"context"

	"astroneko-backend/internal/core/domain/crm_user"
)

type ServiceInterface interface {
	CreateUser(ctx context.Context, req *crm_user.CreateCRMUserRequest) (*crm_user.CRMUser, error)
	Login(ctx context.Context, req *crm_user.CRMLoginRequest) (*crm_user.CRMLoginResponse, error)
	GetUserByID(ctx context.Context, id string) (*crm_user.CRMUser, error)
	ValidateToken(token string) (*crm_user.JWTClaims, error)
}
