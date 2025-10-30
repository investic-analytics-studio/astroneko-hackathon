package referral_code

import (
	"astroneko-backend/internal/core/domain/shared"

	"github.com/google/uuid"
)

type ReferralCode struct {
	shared.NoDeletedModel
	ReferralCode string `json:"referral_code" gorm:"not null"`
}

func (ReferralCode) TableName() string {
	return "astroneko_general_referral_codes"
}

type CreateReferralCodeRequest struct {
	ReferralCode string `json:"referral_code" validate:"required"`
}

type UpdateReferralCodeRequest struct {
	ReferralCode string `json:"referral_code" validate:"required"`
}

type ReferralCodeResponse struct {
	ID           uuid.UUID `json:"id"`
	ReferralCode string    `json:"referral_code"`
	UsedCount    int64     `json:"used_count"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func (r *ReferralCode) ToResponse(usedCount int64) *ReferralCodeResponse {
	return &ReferralCodeResponse{
		ID:           r.ID,
		ReferralCode: r.ReferralCode,
		UsedCount:    usedCount,
		CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type ListReferralCodesResponse struct {
	Codes  []ReferralCodeResponse `json:"codes"`
	Total  int64                  `json:"total"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
}
