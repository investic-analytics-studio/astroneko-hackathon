package user_limit

import (
	"astroneko-backend/internal/core/domain/shared"
)

type UserLimit struct {
	shared.NoDeletedModel
	Limit int `json:"limit" gorm:"default:10000;not null"`
}

func (UserLimit) TableName() string {
	return "astroneko_user_limit"
}

func (ul *UserLimit) ToResponse() *UserLimitResponse {
	return &UserLimitResponse{
		ID:        ul.ID.String(),
		Limit:     ul.Limit,
		CreatedAt: ul.CreatedAt,
		UpdatedAt: ul.UpdatedAt,
	}
}
