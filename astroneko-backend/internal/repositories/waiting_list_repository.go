package repositories

import (
	"context"

	"astroneko-backend/internal/core/domain/waiting_list"
	"astroneko-backend/internal/core/ports"
	waitingListPorts "astroneko-backend/internal/core/ports/waiting_list"
)

type waitingListRepository struct {
	*GenericRepository[waiting_list.WaitingListUser]
	db ports.DatabaseInterface
}

func NewWaitingListRepository(db ports.DatabaseInterface) waitingListPorts.RepositoryInterface {
	genericRepo := NewGenericRepository[waiting_list.WaitingListUser](db)
	return &waitingListRepository{
		GenericRepository: genericRepo,
		db:                db,
	}
}

func (r *waitingListRepository) Create(ctx context.Context, waitingListUser *waiting_list.WaitingListUser) (*waiting_list.WaitingListUser, error) {
	return r.GenericRepository.Create(ctx, waitingListUser)
}

func (r *waitingListRepository) GetByEmail(ctx context.Context, email string) (*waiting_list.WaitingListUser, error) {
	var waitingListUser waiting_list.WaitingListUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&waitingListUser); err != nil {
		return nil, err
	}
	return &waitingListUser, nil
}

func (r *waitingListRepository) GetByID(ctx context.Context, id string) (*waiting_list.WaitingListUser, error) {
	return r.GenericRepository.GetByID(ctx, id)
}

func (r *waitingListRepository) List(ctx context.Context, limit, offset int) ([]*waiting_list.WaitingListUser, int64, error) {
	var waitingListUsers []*waiting_list.WaitingListUser
	var count int64

	if err := r.db.WithContext(ctx).Model(&waiting_list.WaitingListUser{}).Count(&count); err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&waitingListUsers); err != nil {
		return nil, 0, err
	}

	return waitingListUsers, count, nil
}

func (r *waitingListRepository) Delete(ctx context.Context, id string) error {
	return r.GenericRepository.Delete(ctx, id)
}
