package history

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a conversation session for a user
type Session struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index:idx_session_user_id;not null" json:"user_id"`
	HistoryName string         `gorm:"type:text" json:"history_name"`
	CreatedAt   time.Time      `gorm:"not null;default:now();index:idx_session_created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the table name used by Session
func (Session) TableName() string {
	return "astroneko_sessions"
}

// BeforeCreate will set a UUID rather than numeric ID.
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
