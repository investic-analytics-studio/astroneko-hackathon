package firebase

import (
	"context"
	"firebase.google.com/go/v4/auth"
)

// FirebaseClientInterface defines the contract for Firebase Auth operations
type FirebaseClientInterface interface {
	CreateUser(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error)
	DeleteUser(ctx context.Context, uid string) error
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}