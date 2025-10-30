package services

import (
	"context"

	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/waiting_list"
	waitingListPorts "astroneko-backend/internal/core/ports/waiting_list"
	"astroneko-backend/pkg/logger"
)

type WaitingListService struct {
	waitingListRepo waitingListPorts.RepositoryInterface
	logger          logger.Logger
}

func NewWaitingListService(waitingListRepo waitingListPorts.RepositoryInterface, log logger.Logger) *WaitingListService {
	return &WaitingListService{
		waitingListRepo: waitingListRepo,
		logger:          log,
	}
}

func (s *WaitingListService) JoinWaitingList(ctx context.Context, email string) (*waiting_list.WaitingListUser, error) {
	// Check if user already exists in waiting list
	existingUser, err := s.waitingListRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		s.logger.Warn("User already exists in waiting list",
			logger.Field{Key: "module", Value: "waiting_list_service"},
			logger.Field{Key: "email", Value: email})
		return nil, shared.ErrWaitingListUserAlreadyExists
	}

	// Create new waiting list user
	newWaitingListUser := &waiting_list.WaitingListUser{
		Email: email,
	}

	waitingListUser, err := s.waitingListRepo.Create(ctx, newWaitingListUser)
	if err != nil {
		s.logger.Error("Failed to add user to waiting list",
			logger.Field{Key: "module", Value: "waiting_list_service"},
			logger.Field{Key: "email", Value: email},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrWaitingListUserCreationFailed
	}

	s.logger.Info("User successfully added to waiting list",
		logger.Field{Key: "module", Value: "waiting_list_service"},
		logger.Field{Key: "user_id", Value: waitingListUser.ID.String()},
		logger.Field{Key: "email", Value: email})

	return waitingListUser, nil
}

func (s *WaitingListService) GetWaitingListUsers(ctx context.Context, limit, offset int) ([]*waiting_list.WaitingListUser, int64, error) {
	return s.waitingListRepo.List(ctx, limit, offset)
}

func (s *WaitingListService) GetWaitingListUserByEmail(ctx context.Context, email string) (*waiting_list.WaitingListUser, error) {
	return s.waitingListRepo.GetByEmail(ctx, email)
}

func (s *WaitingListService) IsInWaitingListByEmail(ctx context.Context, email string) (bool, error) {
	_, err := s.waitingListRepo.GetByEmail(ctx, email)
	if err != nil {
		// If user not found, they're not in the waiting list
		return false, nil
	}
	// If no error, user exists in waiting list
	return true, nil
}
