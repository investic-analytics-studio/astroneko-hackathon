package history

import (
	"time"

	"github.com/google/uuid"
)

// SortOrder represents the sorting order
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// SortField represents the field to sort by
type SortField string

const (
	SortByCreatedAt SortField = "created_at"
	SortByUpdatedAt SortField = "updated_at"
)

// SessionSummary represents a summary of a session for listing
type SessionSummary struct {
	SessionID   uuid.UUID `json:"session_id"`
	HistoryName string    `json:"history_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetSessionsResponse represents the response for listing user sessions
type GetSessionsResponse struct {
	Sessions []SessionSummary `json:"sessions"`
	Total    int              `json:"total"`
}

// MessageDetail represents a message in the history
type MessageDetail struct {
	ID         uuid.UUID `json:"id"`
	Message    string    `json:"message"`
	Role       string    `json:"role"`
	UsedTokens int       `json:"used_tokens"`
	CreatedAt  time.Time `json:"created_at"`
	Card       string    `json:"card,omitempty"`
	Meaning    string    `json:"meaning,omitempty"`
}

// GetMessagesResponse represents the response for getting messages in a session
type GetMessagesResponse struct {
	SessionID   uuid.UUID       `json:"session_id"`
	HistoryName string          `json:"history_name"`
	Messages    []MessageDetail `json:"messages"`
	Total       int             `json:"total"`
}
