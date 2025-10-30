package history

import (
	"context"

	"astroneko-backend/internal/core/domain/history"

	"github.com/google/uuid"
)

// ServiceInterface defines the contract for history business logic
type ServiceInterface interface {
	// GetUserSessions retrieves all sessions for a given user with optional sorting and search
	GetUserSessions(ctx context.Context, userID uuid.UUID, sortBy string, sortOrder string, searchQuery string) (*history.GetSessionsResponse, error)

	// GetSessionMessages retrieves all messages for a given session with optional sorting
	// Validates that the session belongs to the user
	GetSessionMessages(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, sortOrder string) (*history.GetMessagesResponse, error)

	// DeleteSession soft deletes a session (validates ownership)
	DeleteSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
}
