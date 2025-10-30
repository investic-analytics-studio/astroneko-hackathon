package repositories

import (
	"context"

	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/internal/core/ports"
	userPorts "astroneko-backend/internal/core/ports/user"
	"github.com/google/uuid"
)

type userRepository struct {
	db ports.DatabaseInterface
}

func NewUserRepository(db ports.DatabaseInterface) userPorts.RepositoryInterface {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	user.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var user user.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*user.User, error) {
	var user user.User
	if err := r.db.WithContext(ctx).Where("firebase_uid = ?", firebaseUID).First(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *user.User) (*user.User, error) {
	if err := r.db.WithContext(ctx).Save(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Where("id = ?", userID).Delete(&user.User{})
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*user.User, int64, error) {
	var users []*user.User
	var count int64

	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&count); err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users); err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *userRepository) GetTotalUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&count); err != nil {
		return 0, err
	}
	return count, nil
}
