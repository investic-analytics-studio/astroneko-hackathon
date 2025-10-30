package services

import (
	"context"

	"astroneko-backend/internal/core/domain/astro_boxing_waiting_list"
	"astroneko-backend/internal/core/domain/shared"
	astroBoxingWaitingListPorts "astroneko-backend/internal/core/ports/astro_boxing_waiting_list"
	"astroneko-backend/pkg/logger"
)

type AstroBoxingWaitingListService struct {
	astroBoxingWaitingListRepo astroBoxingWaitingListPorts.RepositoryInterface
	logger                     logger.Logger
}

func NewAstroBoxingWaitingListService(astroBoxingWaitingListRepo astroBoxingWaitingListPorts.RepositoryInterface, log logger.Logger) *AstroBoxingWaitingListService {
	return &AstroBoxingWaitingListService{
		astroBoxingWaitingListRepo: astroBoxingWaitingListRepo,
		logger:                     log,
	}
}

func (s *AstroBoxingWaitingListService) JoinWaitingList(ctx context.Context, email string) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error) {
	existingUser, err := s.astroBoxingWaitingListRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		s.logger.Warn("User already exists in astro boxing waiting list",
			logger.Field{Key: "module", Value: "astro_boxing_waiting_list_service"},
			logger.Field{Key: "email", Value: email})
		return nil, shared.ErrWaitingListUserAlreadyExists
	}

	newUser := &astro_boxing_waiting_list.AstroBoxingWaitingListUser{
		Email: email,
	}

	user, err := s.astroBoxingWaitingListRepo.Create(ctx, newUser)
	if err != nil {
		s.logger.Error("Failed to add user to astro boxing waiting list",
			logger.Field{Key: "module", Value: "astro_boxing_waiting_list_service"},
			logger.Field{Key: "email", Value: email},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrWaitingListUserCreationFailed
	}

	s.logger.Info("User successfully added to astro boxing waiting list",
		logger.Field{Key: "module", Value: "astro_boxing_waiting_list_service"},
		logger.Field{Key: "user_id", Value: user.ID.String()},
		logger.Field{Key: "email", Value: email})

	return user, nil
}

func (s *AstroBoxingWaitingListService) GetWaitingListUsers(ctx context.Context, limit, offset int) ([]*astro_boxing_waiting_list.AstroBoxingWaitingListUser, int64, error) {
	return s.astroBoxingWaitingListRepo.List(ctx, limit, offset)
}

func (s *AstroBoxingWaitingListService) GetWaitingListUserByEmail(ctx context.Context, email string) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error) {
	return s.astroBoxingWaitingListRepo.GetByEmail(ctx, email)
}

func (s *AstroBoxingWaitingListService) IsInWaitingListByEmail(ctx context.Context, email string) (bool, error) {
	_, err := s.astroBoxingWaitingListRepo.GetByEmail(ctx, email)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *AstroBoxingWaitingListService) DeleteUser(ctx context.Context, id string) error {
	return s.astroBoxingWaitingListRepo.Delete(ctx, id)
}
