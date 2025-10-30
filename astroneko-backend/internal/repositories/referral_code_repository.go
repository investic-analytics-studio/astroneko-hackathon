package repositories

import (
	"context"

	"astroneko-backend/internal/core/domain/referral_code"
	"astroneko-backend/internal/core/ports"
	referralCodePorts "astroneko-backend/internal/core/ports/referral_code"

	"github.com/google/uuid"
)

const (
	WhereClauseReferralCodeEqual                = "LOWER(referral_code) = LOWER(?)"
	WhereClauseReferralCodeEqualAndNotActivated = "LOWER(referral_code) = LOWER(?) AND is_activated = false"
)

type referralCodeRepository struct {
	db ports.DatabaseInterface
}

func NewReferralCodeRepository(db ports.DatabaseInterface) referralCodePorts.RepositoryInterface {
	return &referralCodeRepository{
		db: db,
	}
}

func (r *referralCodeRepository) Create(ctx context.Context, referralCode *referral_code.ReferralCode) (*referral_code.ReferralCode, error) {
	referralCode.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(referralCode); err != nil {
		return nil, err
	}
	return referralCode, nil
}

func (r *referralCodeRepository) GetByID(ctx context.Context, id string) (*referral_code.ReferralCode, error) {
	referralCodeID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var referralCode referral_code.ReferralCode
	if err := r.db.WithContext(ctx).Where("id = ?", referralCodeID).First(&referralCode); err != nil {
		return nil, err
	}
	return &referralCode, nil
}

func (r *referralCodeRepository) GetByReferralCode(ctx context.Context, code string) (*referral_code.ReferralCode, error) {
	var referralCode referral_code.ReferralCode
	if err := r.db.WithContext(ctx).Where(WhereClauseReferralCodeEqual, code).First(&referralCode); err != nil {
		return nil, err
	}
	return &referralCode, nil
}

func (r *referralCodeRepository) IsValidReferralCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&referral_code.ReferralCode{}).Where(WhereClauseReferralCodeEqual, code).Count(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *referralCodeRepository) Update(ctx context.Context, referralCode *referral_code.ReferralCode) (*referral_code.ReferralCode, error) {
	if err := r.db.WithContext(ctx).Save(referralCode); err != nil {
		return nil, err
	}
	return referralCode, nil
}

func (r *referralCodeRepository) Delete(ctx context.Context, id string) error {
	referralCodeID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Where("id = ?", referralCodeID).Delete(&referral_code.ReferralCode{})
}

func (r *referralCodeRepository) List(ctx context.Context, limit, offset int) ([]*referral_code.ReferralCode, int64, error) {
	var referralCodes []*referral_code.ReferralCode
	var count int64

	if err := r.db.WithContext(ctx).Model(&referral_code.ReferralCode{}).Count(&count); err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&referralCodes); err != nil {
		return nil, 0, err
	}

	return referralCodes, count, nil
}

// CreateUserReferralCode creates a new user referral code
func (r *referralCodeRepository) CreateUserReferralCode(ctx context.Context, userReferralCode *referral_code.UserReferralCode) (*referral_code.UserReferralCode, error) {
	userReferralCode.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(userReferralCode); err != nil {
		return nil, err
	}
	return userReferralCode, nil
}

// GetUserReferralCodesByUserID gets all user referral codes for a specific user, sorted by created_at
func (r *referralCodeRepository) GetUserReferralCodesByUserID(ctx context.Context, userID uuid.UUID) ([]*referral_code.UserReferralCode, error) {
	var userReferralCodes []*referral_code.UserReferralCode
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at ASC").Find(&userReferralCodes); err != nil {
		return nil, err
	}
	return userReferralCodes, nil
}

// GetUserReferralCodeByCode gets a user referral code by its code
func (r *referralCodeRepository) GetUserReferralCodeByCode(ctx context.Context, code string) (*referral_code.UserReferralCode, error) {
	var userReferralCode referral_code.UserReferralCode
	if err := r.db.WithContext(ctx).Where(WhereClauseReferralCodeEqual, code).First(&userReferralCode); err != nil {
		return nil, err
	}
	return &userReferralCode, nil
}

// UpdateUserReferralCode updates a user referral code
func (r *referralCodeRepository) UpdateUserReferralCode(ctx context.Context, userReferralCode *referral_code.UserReferralCode) (*referral_code.UserReferralCode, error) {
	if err := r.db.WithContext(ctx).Save(userReferralCode); err != nil {
		return nil, err
	}
	return userReferralCode, nil
}

// IsValidUserReferralCode checks if a user referral code is valid and not activated
func (r *referralCodeRepository) IsValidUserReferralCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&referral_code.UserReferralCode{}).Where(WhereClauseReferralCodeEqualAndNotActivated, code).Count(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateReferralLog creates a new referral log entry
func (r *referralCodeRepository) CreateReferralLog(ctx context.Context, referralLog *referral_code.ReferralLog) (*referral_code.ReferralLog, error) {
	referralLog.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(referralLog); err != nil {
		return nil, err
	}
	return referralLog, nil
}

// GetReferralCodeUsageCount gets the usage count for a general referral code
func (r *referralCodeRepository) GetReferralCodeUsageCount(ctx context.Context, referralCode string) (int64, error) {
	// First get the referral code ID by the code string
	generalCode, err := r.GetByReferralCode(ctx, referralCode)
	if err != nil {
		return 0, err
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&referral_code.ReferralLog{}).
		Where("referral_code_id = ? AND code_type = ?", generalCode.ID, "general").
		Count(&count); err != nil {
		return 0, err
	}
	return count, nil
}
