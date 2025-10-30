package services

import (
	"context"

	"astroneko-backend/internal/core/domain/user_limit"
	userPorts "astroneko-backend/internal/core/ports/user"
	userLimitPorts "astroneko-backend/internal/core/ports/user_limit"
	"astroneko-backend/pkg/logger"
)

type UserLimitService struct {
	userLimitRepo userLimitPorts.Repository
	userRepo      userPorts.RepositoryInterface
	logger        logger.Logger
}

func NewUserLimitService(userLimitRepo userLimitPorts.Repository, userRepo userPorts.RepositoryInterface, logger logger.Logger) *UserLimitService {
	return &UserLimitService{
		userLimitRepo: userLimitRepo,
		userRepo:      userRepo,
		logger:        logger,
	}
}

func (s *UserLimitService) GetUserLimit(ctx context.Context) (*user_limit.UserLimit, error) {
	userLimit, err := s.userLimitRepo.GetUserLimit(ctx)
	if err != nil {
		s.logger.Error("Failed to get user limit",
			logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	return userLimit, nil
}

func (s *UserLimitService) UpdateUserLimit(ctx context.Context, req *user_limit.UpdateUserLimitRequest) (*user_limit.UserLimit, error) {
	userLimit, err := s.userLimitRepo.UpdateUserLimit(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update user limit",
			logger.Field{Key: "error", Value: err.Error()},
			logger.Field{Key: "limit", Value: req.Limit})
		return nil, err
	}

	s.logger.Info("User limit updated successfully",
		logger.Field{Key: "limit", Value: req.Limit})

	return userLimit, nil
}

func (s *UserLimitService) IsUserOverLimitUsed(ctx context.Context) (bool, error) {
	userLimit, err := s.userLimitRepo.GetUserLimit(ctx)
	if err != nil {
		s.logger.Error("Failed to get user limit",
			logger.Field{Key: "error", Value: err.Error()})
		return false, err
	}

	totalUsers, err := s.userRepo.GetTotalUsers(ctx)
	if err != nil {
		s.logger.Error("Failed to get total users count",
			logger.Field{Key: "error", Value: err.Error()})
		return false, err
	}

	isOverLimit := totalUsers >= int64(userLimit.Limit)

	s.logger.Info("User limit check completed",
		logger.Field{Key: "total_users", Value: totalUsers},
		logger.Field{Key: "user_limit", Value: userLimit.Limit},
		logger.Field{Key: "is_over_limit", Value: isOverLimit})

	return isOverLimit, nil
}
