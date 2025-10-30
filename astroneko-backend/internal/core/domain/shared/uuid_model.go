package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type NoDeletedModel struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid();omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoDeletedUintModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
