package repositories

import (
	"context"

	"astroneko-backend/internal/core/domain/user_limit"
	"astroneko-backend/internal/core/ports"
	userLimitPorts "astroneko-backend/internal/core/ports/user_limit"
)

type userLimitRepository struct {
	db ports.DatabaseInterface
}

func NewUserLimitRepository(db ports.DatabaseInterface) userLimitPorts.Repository {
	return &userLimitRepository{
		db: db,
	}
}

func (r *userLimitRepository) GetUserLimit(ctx context.Context) (*user_limit.UserLimit, error) {
	var userLimit user_limit.UserLimit

	// Get the first (and only) row from the table
	if err := r.db.WithContext(ctx).First(&userLimit); err != nil {
		return nil, err
	}

	return &userLimit, nil
}

func (r *userLimitRepository) UpdateUserLimit(ctx context.Context, req *user_limit.UpdateUserLimitRequest) (*user_limit.UserLimit, error) {
	var userLimit user_limit.UserLimit

	// Get the first row
	if err := r.db.WithContext(ctx).First(&userLimit); err != nil {
		return nil, err
	}

	// Update the limit
	userLimit.Limit = req.Limit

	// Save the changes
	if err := r.db.WithContext(ctx).Save(&userLimit); err != nil {
		return nil, err
	}

	return &userLimit, nil
}
