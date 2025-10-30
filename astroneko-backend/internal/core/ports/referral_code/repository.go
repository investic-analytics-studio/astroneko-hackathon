package referral_code

import (
	"context"

	"astroneko-backend/internal/core/domain/referral_code"
	"github.com/google/uuid"
)

// RepositoryInterface defines the contract for referral code data operations
type RepositoryInterface interface {
	// General referral codes
	Create(ctx context.Context, referralCode *referral_code.ReferralCode) (*referral_code.ReferralCode, error)
	GetByID(ctx context.Context, id string) (*referral_code.ReferralCode, error)
	GetByReferralCode(ctx context.Context, code string) (*referral_code.ReferralCode, error)
	IsValidReferralCode(ctx context.Context, code string) (bool, error)
	Update(ctx context.Context, referralCode *referral_code.ReferralCode) (*referral_code.ReferralCode, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*referral_code.ReferralCode, int64, error)

	// User referral codes
	CreateUserReferralCode(ctx context.Context, userReferralCode *referral_code.UserReferralCode) (*referral_code.UserReferralCode, error)
	GetUserReferralCodesByUserID(ctx context.Context, userID uuid.UUID) ([]*referral_code.UserReferralCode, error)
	GetUserReferralCodeByCode(ctx context.Context, code string) (*referral_code.UserReferralCode, error)
	UpdateUserReferralCode(ctx context.Context, userReferralCode *referral_code.UserReferralCode) (*referral_code.UserReferralCode, error)
	IsValidUserReferralCode(ctx context.Context, code string) (bool, error)

	// Referral logs
	CreateReferralLog(ctx context.Context, referralLog *referral_code.ReferralLog) (*referral_code.ReferralLog, error)
	GetReferralCodeUsageCount(ctx context.Context, referralCode string) (int64, error)
}
