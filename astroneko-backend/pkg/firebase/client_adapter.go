package firebase

import (
	"context"
	"firebase.google.com/go/v4/auth"
)

// FirebaseClientAdapter wraps the actual Firebase client to implement our interface
type FirebaseClientAdapter struct {
	client *auth.Client
}

// NewFirebaseClientAdapter creates a new adapter for the Firebase client
func NewFirebaseClientAdapter(client *auth.Client) FirebaseClientInterface {
	return &FirebaseClientAdapter{client: client}
}

// CreateUser delegates to the underlying Firebase client
func (a *FirebaseClientAdapter) CreateUser(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
	return a.client.CreateUser(ctx, user)
}

// DeleteUser delegates to the underlying Firebase client
func (a *FirebaseClientAdapter) DeleteUser(ctx context.Context, uid string) error {
	return a.client.DeleteUser(ctx, uid)
}

// GetUser delegates to the underlying Firebase client
func (a *FirebaseClientAdapter) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	return a.client.GetUser(ctx, uid)
}

// VerifyIDToken delegates to the underlying Firebase client
func (a *FirebaseClientAdapter) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return a.client.VerifyIDToken(ctx, idToken)
}