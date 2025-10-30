package referral_code

import (
	"astroneko-backend/internal/core/domain/shared"

	"github.com/google/uuid"
)

type UserReferralCode struct {
	shared.NoDeletedModel
	UserID       uuid.UUID `json:"user_id" gorm:"not null;type:uuid"`
	ReferralCode string    `json:"referral_code" gorm:"not null"`
	IsActivated  bool      `json:"is_activated" gorm:"default:false;not null"`
}

func (UserReferralCode) TableName() string {
	return "astroneko_user_referral_codes"
}

type UserReferralCodeResponse struct {
	ID           uuid.UUID `json:"id"`
	ReferralCode string    `json:"referral_code"`
	IsActivated  bool      `json:"is_activated"`
}

func (u *UserReferralCode) ToResponse() *UserReferralCodeResponse {
	return &UserReferralCodeResponse{
		ID:           u.ID,
		ReferralCode: u.ReferralCode,
		IsActivated:  u.IsActivated,
	}
}

type GetUserReferralCodesResponse struct {
	Codes []UserReferralCodeResponse `json:"codes"`
}
