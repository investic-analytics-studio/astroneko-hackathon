package history

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a single message in a conversation history
type Message struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SessionID  uuid.UUID `gorm:"type:uuid;index:idx_message_session_id;not null" json:"session_id"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	Role       string    `gorm:"type:varchar(255);not null" json:"role"`
	UsedTokens int       `gorm:"type:int4;default:0" json:"used_tokens"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_message_created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationship
	Session Session `gorm:"foreignKey:SessionID;references:ID" json:"-"`
}

// TableName overrides the table name used by Message
func (Message) TableName() string {
	return "astroneko_message_histories"
}
