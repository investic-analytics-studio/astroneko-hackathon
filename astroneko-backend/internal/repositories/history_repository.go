package repositories

import (
	"context"
	"fmt"

	"astroneko-backend/internal/core/domain/history"
	"astroneko-backend/internal/core/ports"
	historyPorts "astroneko-backend/internal/core/ports/history"

	"github.com/google/uuid"
)

type historyRepository struct {
	db ports.DatabaseInterface
}

// NewHistoryRepository creates a new history repository instance
func NewHistoryRepository(db ports.DatabaseInterface) historyPorts.RepositoryInterface {
	return &historyRepository{
		db: db,
	}
}

// GetSessionsByUserID retrieves all sessions for a specific user
// sortBy: "created_at" or "updated_at" (default: "updated_at")
// sortOrder: "asc" or "desc" (default: "desc")
// searchQuery: optional text to search in history_name (partial match)
func (r *historyRepository) GetSessionsByUserID(ctx context.Context, userID uuid.UUID, sortBy string, sortOrder string, searchQuery string) ([]history.Session, error) {
	var sessions []history.Session

	// Validate and set defaults
	if sortBy == "" {
		sortBy = "updated_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Whitelist sortBy to prevent SQL injection
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "updated_at"
	}

	// Whitelist sortOrder to prevent SQL injection
	validSortOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}
	if !validSortOrders[sortOrder] {
		sortOrder = "desc"
	}

	// Convert to uppercase for SQL
	if sortOrder == "asc" {
		sortOrder = "ASC"
	} else {
		sortOrder = "DESC"
	}

	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// Add search filter if searchQuery is provided
	if searchQuery != "" {
		query = query.Where("history_name ILIKE ?", "%"+searchQuery+"%")
	}

	err := query.
		Order(orderClause).
		Find(&sessions)

	if err != nil {
		return nil, fmt.Errorf("failed to get sessions for user %s: %w", userID, err)
	}

	return sessions, nil
}

// GetSessionByID retrieves a session by its ID
func (r *historyRepository) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*history.Session, error) {
	var session history.Session

	err := r.db.WithContext(ctx).
		Where("id = ?", sessionID).
		First(&session)

	if err != nil {
		return nil, fmt.Errorf("failed to get session %s: %w", sessionID, err)
	}

	return &session, nil
}

// ValidateSessionOwnership checks if a session belongs to a specific user
// Returns true if the session belongs to the user, false otherwise
func (r *historyRepository) ValidateSessionOwnership(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&history.Session{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Count(&count)

	if err != nil {
		return false, fmt.Errorf("failed to validate session ownership: %w", err)
	}

	return count > 0, nil
}

// GetMessagesBySessionID retrieves all messages for a specific session
// sortOrder: "asc" or "desc" (default: "asc" for chronological order)
func (r *historyRepository) GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID, sortOrder string) ([]history.Message, error) {
	var messages []history.Message

	// Validate and set default
	if sortOrder == "" {
		sortOrder = "asc"
	}

	// Whitelist sortOrder to prevent SQL injection
	validSortOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}
	if !validSortOrders[sortOrder] {
		sortOrder = "asc"
	}

	// Convert to uppercase for SQL
	if sortOrder == "asc" {
		sortOrder = "ASC"
	} else {
		sortOrder = "DESC"
	}

	orderClause := fmt.Sprintf("created_at %s", sortOrder)

	err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order(orderClause).
		Find(&messages)

	if err != nil {
		return nil, fmt.Errorf("failed to get messages for session %s: %w", sessionID, err)
	}

	return messages, nil
}

// DeleteSession soft deletes a session by setting the deleted_at timestamp
func (r *historyRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&history.Session{}, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session %s: %w", sessionID, err)
	}

	return nil
}
