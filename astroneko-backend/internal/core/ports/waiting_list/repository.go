package waiting_list

import (
	"context"

	"astroneko-backend/internal/core/domain/waiting_list"
)

// RepositoryInterface defines the contract for waiting list data operations
type RepositoryInterface interface {
	Create(ctx context.Context, waitingListUser *waiting_list.WaitingListUser) (*waiting_list.WaitingListUser, error)
	GetByEmail(ctx context.Context, email string) (*waiting_list.WaitingListUser, error)
	GetByID(ctx context.Context, id string) (*waiting_list.WaitingListUser, error)
	List(ctx context.Context, limit, offset int) ([]*waiting_list.WaitingListUser, int64, error)
	Delete(ctx context.Context, id string) error
}
