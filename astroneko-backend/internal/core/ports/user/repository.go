package user

import (
	"context"

	"astroneko-backend/internal/core/domain/user"
)

// RepositoryInterface defines the contract for user data operations
type RepositoryInterface interface {
	Create(ctx context.Context, user *user.User) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, user *user.User) (*user.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*user.User, int64, error)
	GetTotalUsers(ctx context.Context) (int64, error)
}
