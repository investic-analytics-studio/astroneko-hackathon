package history

import (
	"context"

	"astroneko-backend/internal/core/domain/history"

	"github.com/google/uuid"
)

// RepositoryInterface defines the contract for history data operations
type RepositoryInterface interface {
	// Session operations
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID, sortBy string, sortOrder string, searchQuery string) ([]history.Session, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*history.Session, error)
	ValidateSessionOwnership(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (bool, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error

	// Message operations
	GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID, sortOrder string) ([]history.Message, error)
}
