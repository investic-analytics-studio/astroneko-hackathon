package astro_boxing_waiting_list

import (
	"context"

	"astroneko-backend/internal/core/domain/astro_boxing_waiting_list"
)

type ServiceInterface interface {
	JoinWaitingList(ctx context.Context, email string) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error)
	GetWaitingListUsers(ctx context.Context, limit, offset int) ([]*astro_boxing_waiting_list.AstroBoxingWaitingListUser, int64, error)
	GetWaitingListUserByEmail(ctx context.Context, email string) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error)
	IsInWaitingListByEmail(ctx context.Context, email string) (bool, error)
	DeleteUser(ctx context.Context, id string) error
}
