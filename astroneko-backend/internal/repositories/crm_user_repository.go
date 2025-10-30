package repositories

import (
	"context"

	"astroneko-backend/internal/core/domain/crm_user"
	"astroneko-backend/internal/core/ports"
	crmUserPorts "astroneko-backend/internal/core/ports/crm_user"

	"github.com/google/uuid"
)

type crmUserRepository struct {
	db ports.DatabaseInterface
}

func NewCRMUserRepository(db ports.DatabaseInterface) crmUserPorts.RepositoryInterface {
	return &crmUserRepository{
		db: db,
	}
}

func (r *crmUserRepository) Create(ctx context.Context, user *crm_user.CRMUser) (*crm_user.CRMUser, error) {
	user.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *crmUserRepository) GetByID(ctx context.Context, id string) (*crm_user.CRMUser, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var user crm_user.CRMUser
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *crmUserRepository) GetByUsername(ctx context.Context, username string) (*crm_user.CRMUser, error) {
	var user crm_user.CRMUser
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *crmUserRepository) Update(ctx context.Context, user *crm_user.CRMUser) (*crm_user.CRMUser, error) {
	if err := r.db.WithContext(ctx).Save(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *crmUserRepository) Delete(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Where("id = ?", userID).Delete(&crm_user.CRMUser{})
}
