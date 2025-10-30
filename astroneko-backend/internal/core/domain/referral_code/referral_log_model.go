package referral_code

import (
	"astroneko-backend/internal/core/domain/shared"

	"github.com/google/uuid"
)

type ReferralLog struct {
	shared.NoDeletedModel
	RedeemedByUserID uuid.UUID  `json:"redeemed_by_user_id" gorm:"not null;type:uuid"`
	CodeType         string     `json:"code_type" gorm:"not null"`
	ReferralCodeID   *uuid.UUID `json:"referral_code_id" gorm:"type:uuid"`
}

func (ReferralLog) TableName() string {
	return "astroneko_referral_logs"
}

type ActivateReferralRequest struct {
	ReferralCode string `json:"referral_code" validate:"required"`
}

type ActivateReferralResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
