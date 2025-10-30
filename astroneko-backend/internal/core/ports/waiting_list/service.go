package waiting_list

import (
	"context"

	"astroneko-backend/internal/core/domain/waiting_list"
)

// ServiceInterface defines the contract for waiting list business logic
type ServiceInterface interface {
	JoinWaitingList(ctx context.Context, email string) (*waiting_list.WaitingListUser, error)
	GetWaitingListUsers(ctx context.Context, limit, offset int) ([]*waiting_list.WaitingListUser, int64, error)
	GetWaitingListUserByEmail(ctx context.Context, email string) (*waiting_list.WaitingListUser, error)
	IsInWaitingListByEmail(ctx context.Context, email string) (bool, error)
}
