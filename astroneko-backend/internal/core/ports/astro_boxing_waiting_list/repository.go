package astro_boxing_waiting_list

import (
	"context"

	"astroneko-backend/internal/core/domain/astro_boxing_waiting_list"
)

type RepositoryInterface interface {
	Create(ctx context.Context, user *astro_boxing_waiting_list.AstroBoxingWaitingListUser) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error)
	GetByEmail(ctx context.Context, email string) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error)
	GetByID(ctx context.Context, id string) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error)
	List(ctx context.Context, limit, offset int) ([]*astro_boxing_waiting_list.AstroBoxingWaitingListUser, int64, error)
	Delete(ctx context.Context, id string) error
}
