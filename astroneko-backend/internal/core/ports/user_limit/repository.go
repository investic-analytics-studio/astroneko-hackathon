package user_limit

import (
	"context"

	"astroneko-backend/internal/core/domain/user_limit"
)

type Repository interface {
	GetUserLimit(ctx context.Context) (*user_limit.UserLimit, error)
	UpdateUserLimit(ctx context.Context, req *user_limit.UpdateUserLimitRequest) (*user_limit.UserLimit, error)
}
