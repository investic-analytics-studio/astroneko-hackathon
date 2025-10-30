package astro_boxing_waiting_list

import "astroneko-backend/internal/core/domain/shared"

type AstroBoxingWaitingListUser struct {
	shared.NoDeletedModel
	Email string `json:"email" gorm:"not null"`
}

func (AstroBoxingWaitingListUser) TableName() string {
	return "astro_boxing_waiting_list_users"
}
