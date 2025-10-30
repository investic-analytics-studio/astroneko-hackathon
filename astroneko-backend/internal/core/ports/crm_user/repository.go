package crm_user

import (
	"context"

	"astroneko-backend/internal/core/domain/crm_user"
)

type RepositoryInterface interface {
	Create(ctx context.Context, user *crm_user.CRMUser) (*crm_user.CRMUser, error)
	GetByID(ctx context.Context, id string) (*crm_user.CRMUser, error)
	GetByUsername(ctx context.Context, username string) (*crm_user.CRMUser, error)
	Update(ctx context.Context, user *crm_user.CRMUser) (*crm_user.CRMUser, error)
	Delete(ctx context.Context, id string) error
}
