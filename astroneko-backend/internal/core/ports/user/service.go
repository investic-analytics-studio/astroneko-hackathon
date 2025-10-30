package user

import (
	"context"

	"astroneko-backend/internal/core/domain/user"
	"firebase.google.com/go/v4/auth"
)

// ServiceInterface defines the contract for user business logic
type ServiceInterface interface {
	CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error)
	GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*user.User, error)
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	UpdateUser(ctx context.Context, id string, req *user.UpdateUserRequest) (*user.User, error)
	DeleteUser(ctx context.Context, id string) error
	VerifyFirebaseToken(ctx context.Context, idToken string) (*auth.Token, error)
	CreateUserFromToken(ctx context.Context, token *auth.Token) (*user.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*user.AuthResponse, error)
	GoogleAuth(ctx context.Context, req *user.GoogleLoginRequest) (*user.AuthResponse, error)
	Login(ctx context.Context, req *user.LoginRequest) (*user.AuthResponse, error)
	GetTotalUsers(ctx context.Context) (int64, error)
}
