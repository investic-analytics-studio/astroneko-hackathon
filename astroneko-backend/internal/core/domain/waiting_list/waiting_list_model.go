package waiting_list

import (
	"astroneko-backend/internal/core/domain/shared"
	"github.com/google/uuid"
)

type WaitingListUser struct {
	shared.NoDeletedModel
	Email string `json:"email" gorm:"not null"`
}

func (w *WaitingListUser) SetID(id uuid.UUID) {
	w.ID = id
}

func (WaitingListUser) TableName() string {
	return "astroneko_waiting_list_users"
}
