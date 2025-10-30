package repositories

import (
	"context"

	"astroneko-backend/internal/core/domain/astro_boxing_waiting_list"
	"astroneko-backend/internal/core/ports"
	astroBoxingWaitingListPorts "astroneko-backend/internal/core/ports/astro_boxing_waiting_list"
	"github.com/google/uuid"
)

type astroBoxingWaitingListRepository struct {
	db ports.DatabaseInterface
}

func NewAstroBoxingWaitingListRepository(db ports.DatabaseInterface) astroBoxingWaitingListPorts.RepositoryInterface {
	return &astroBoxingWaitingListRepository{
		db: db,
	}
}

func (r *astroBoxingWaitingListRepository) Create(ctx context.Context, user *astro_boxing_waiting_list.AstroBoxingWaitingListUser) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error) {
	user.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *astroBoxingWaitingListRepository) GetByEmail(ctx context.Context, email string) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error) {
	var user astro_boxing_waiting_list.AstroBoxingWaitingListUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *astroBoxingWaitingListRepository) GetByID(ctx context.Context, id string) (*astro_boxing_waiting_list.AstroBoxingWaitingListUser, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var user astro_boxing_waiting_list.AstroBoxingWaitingListUser
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *astroBoxingWaitingListRepository) List(ctx context.Context, limit, offset int) ([]*astro_boxing_waiting_list.AstroBoxingWaitingListUser, int64, error) {
	var users []*astro_boxing_waiting_list.AstroBoxingWaitingListUser
	var count int64

	if err := r.db.WithContext(ctx).Model(&astro_boxing_waiting_list.AstroBoxingWaitingListUser{}).Count(&count); err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&users); err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *astroBoxingWaitingListRepository) Delete(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Where("id = ?", userID).Delete(&astro_boxing_waiting_list.AstroBoxingWaitingListUser{})
}
