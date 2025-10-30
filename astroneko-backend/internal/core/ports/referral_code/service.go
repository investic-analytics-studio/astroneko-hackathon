package referral_code

import (
	"context"

	"astroneko-backend/internal/core/domain/referral_code"
	"github.com/google/uuid"
)

// ServiceInterface defines the contract for referral code business logic
type ServiceInterface interface {
	// General referral codes
	IsValidReferralCode(ctx context.Context, code string) (bool, error)
	CreateReferralCode(ctx context.Context, code string) (*referral_code.ReferralCode, error)
	GetReferralCodeByCode(ctx context.Context, code string) (*referral_code.ReferralCode, error)
	DeleteReferralCode(ctx context.Context, id string) error
	ListReferralCodes(ctx context.Context, limit, offset int) ([]*referral_code.ReferralCode, int64, error)

	// User referral codes
	GetOrGenerateUserReferralCodes(ctx context.Context, userID uuid.UUID) ([]*referral_code.UserReferralCode, error)
	ActivateReferralCode(ctx context.Context, userID uuid.UUID, code string) (*referral_code.ActivateReferralResponse, error)
}
