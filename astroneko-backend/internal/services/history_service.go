package services

import (
	"context"
	"fmt"

	"astroneko-backend/internal/core/domain/history"
	historyPorts "astroneko-backend/internal/core/ports/history"
	"astroneko-backend/pkg/logger"

	"github.com/google/uuid"
)

// HistoryService provides business logic for conversation history
type HistoryService struct {
	historyRepo historyPorts.RepositoryInterface
	logger      logger.Logger
}

// NewHistoryService creates a new history service instance
func NewHistoryService(historyRepo historyPorts.RepositoryInterface, log logger.Logger) *HistoryService {
	return &HistoryService{
		historyRepo: historyRepo,
		logger:      log,
	}
}

// GetUserSessions retrieves all sessions for a given user with optional sorting and search
// sortBy: "created_at" or "updated_at" (default: "updated_at")
// sortOrder: "asc" or "desc" (default: "desc")
// searchQuery: optional text to search in history_name (partial match)
func (s *HistoryService) GetUserSessions(ctx context.Context, userID uuid.UUID, sortBy string, sortOrder string, searchQuery string) (*history.GetSessionsResponse, error) {
	sessions, err := s.historyRepo.GetSessionsByUserID(ctx, userID, sortBy, sortOrder, searchQuery)
	if err != nil {
		s.logger.Error("Failed to retrieve user sessions",
			logger.Field{Key: "module", Value: "history_service"},
			logger.Field{Key: "user_id", Value: userID.String()},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to retrieve sessions: %w", err)
	}

	// Transform to response format
	sessionSummaries := make([]history.SessionSummary, 0, len(sessions))
	for _, session := range sessions {
		sessionSummaries = append(sessionSummaries, history.SessionSummary{
			SessionID:   session.ID,
			HistoryName: session.HistoryName,
			CreatedAt:   session.CreatedAt,
			UpdatedAt:   session.UpdatedAt,
		})
	}

	return &history.GetSessionsResponse{
		Sessions: sessionSummaries,
		Total:    len(sessionSummaries),
	}, nil
}

// GetSessionMessages retrieves all messages for a given session with optional sorting
// Validates that the session belongs to the user before returning messages
// sortOrder: "asc" or "desc" (default: "asc" for chronological order)
func (s *HistoryService) GetSessionMessages(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, sortOrder string) (*history.GetMessagesResponse, error) {
	// First, validate session ownership
	isOwner, err := s.historyRepo.ValidateSessionOwnership(ctx, sessionID, userID)
	if err != nil {
		s.logger.Error("Failed to validate session ownership",
			logger.Field{Key: "module", Value: "history_service"},
			logger.Field{Key: "user_id", Value: userID.String()},
			logger.Field{Key: "session_id", Value: sessionID.String()},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to validate session ownership: %w", err)
	}

	if !isOwner {
		s.logger.Warn("Unauthorized session access attempt",
			logger.Field{Key: "module", Value: "history_service"},
			logger.Field{Key: "user_id", Value: userID.String()},
			logger.Field{Key: "session_id", Value: sessionID.String()})
		return nil, fmt.Errorf("session not found or access denied")
	}

	// Get session details
	session, err := s.historyRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		s.logger.Error("Failed to retrieve session",
			logger.Field{Key: "module", Value: "history_service"},
			logger.Field{Key: "session_id", Value: sessionID.String()},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	// Get messages for the session
	messages, err := s.historyRepo.GetMessagesBySessionID(ctx, sessionID, sortOrder)
	if err != nil {
		s.logger.Error("Failed to retrieve session messages",
			logger.Field{Key: "module", Value: "history_service"},
			logger.Field{Key: "session_id", Value: sessionID.String()},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	// Transform to response format with JSON extraction
	messageDetails := make([]history.MessageDetail, 0, len(messages))
	for _, msg := range messages {
		// Extract JSON from message if present
		cleanedMessage, card, meaning := history.ExtractJSONFromMessage(msg.Message)

		messageDetails = append(messageDetails, history.MessageDetail{
			ID:         msg.ID,
			Message:    cleanedMessage,
			Role:       msg.Role,
			UsedTokens: msg.UsedTokens,
			CreatedAt:  msg.CreatedAt,
			Card:       card,
			Meaning:    meaning,
		})
	}

	return &history.GetMessagesResponse{
		SessionID:   session.ID,
		HistoryName: session.HistoryName,
		Messages:    messageDetails,
		Total:       len(messageDetails),
	}, nil
}

// DeleteSession soft deletes a session by setting the deleted_at timestamp
// Validates that the session belongs to the user before deletion
func (s *HistoryService) DeleteSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	// First, validate session ownership
	isOwner, err := s.historyRepo.ValidateSessionOwnership(ctx, sessionID, userID)
	if err != nil {
		s.logger.Error("Failed to validate session ownership for deletion",
			logger.Field{Key: "module", Value: "history_service"},
			logger.Field{Key: "user_id", Value: userID.String()},
			logger.Field{Key: "session_id", Value: sessionID.String()},
			logger.Field{Key: "error", Value: err.Error()})
		return fmt.Errorf("failed to validate session ownership: %w", err)
	}

	if !isOwner {
		s.logger.Warn("Unauthorized session deletion attempt",
			logger.Field{Key: "module", Value: "history_service"},
			logger.Field{Key: "user_id", Value: userID.String()},
			logger.Field{Key: "session_id", Value: sessionID.String()})
		return fmt.Errorf("session not found or access denied")
	}

	// Delete the session
	err = s.historyRepo.DeleteSession(ctx, sessionID)
	if err != nil {
		s.logger.Error("Failed to delete session",
			logger.Field{Key: "module", Value: "history_service"},
			logger.Field{Key: "user_id", Value: userID.String()},
			logger.Field{Key: "session_id", Value: sessionID.String()},
			logger.Field{Key: "error", Value: err.Error()})
		return fmt.Errorf("failed to delete session: %w", err)
	}

	s.logger.Info("Session deleted successfully",
		logger.Field{Key: "module", Value: "history_service"},
		logger.Field{Key: "user_id", Value: userID.String()},
		logger.Field{Key: "session_id", Value: sessionID.String()})

	return nil
}
