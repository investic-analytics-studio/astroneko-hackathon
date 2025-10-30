package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"astroneko-backend/internal/core/domain/referral_code"
	referralCodePorts "astroneko-backend/internal/core/ports/referral_code"
	userPorts "astroneko-backend/internal/core/ports/user"
	"astroneko-backend/pkg/logger"

	"github.com/google/uuid"
)

type ReferralCodeService struct {
	referralCodeRepo referralCodePorts.RepositoryInterface
	userRepo         userPorts.RepositoryInterface
	logger           logger.Logger
}

func NewReferralCodeService(referralCodeRepo referralCodePorts.RepositoryInterface, userRepo userPorts.RepositoryInterface, log logger.Logger) *ReferralCodeService {
	return &ReferralCodeService{
		referralCodeRepo: referralCodeRepo,
		userRepo:         userRepo,
		logger:           log,
	}
}

func (s *ReferralCodeService) IsValidReferralCode(ctx context.Context, code string) (bool, error) {
	return s.referralCodeRepo.IsValidReferralCode(ctx, code)
}

func (s *ReferralCodeService) CreateReferralCode(ctx context.Context, code string) (*referral_code.ReferralCode, error) {
	newReferralCode := &referral_code.ReferralCode{
		ReferralCode: code,
	}

	return s.referralCodeRepo.Create(ctx, newReferralCode)
}

func (s *ReferralCodeService) GetReferralCodeByCode(ctx context.Context, code string) (*referral_code.ReferralCode, error) {
	return s.referralCodeRepo.GetByReferralCode(ctx, code)
}

func (s *ReferralCodeService) GetReferralCodeByID(ctx context.Context, id string) (*referral_code.ReferralCode, error) {
	return s.referralCodeRepo.GetByID(ctx, id)
}

func (s *ReferralCodeService) GetReferralCodeUsageCount(ctx context.Context, referralCode string) (int64, error) {
	return s.referralCodeRepo.GetReferralCodeUsageCount(ctx, referralCode)
}

func (s *ReferralCodeService) UpdateReferralCode(ctx context.Context, referralCode *referral_code.ReferralCode) (*referral_code.ReferralCode, error) {
	return s.referralCodeRepo.Update(ctx, referralCode)
}

func (s *ReferralCodeService) DeleteReferralCode(ctx context.Context, id string) error {
	return s.referralCodeRepo.Delete(ctx, id)
}

func (s *ReferralCodeService) ListReferralCodes(ctx context.Context, limit, offset int) ([]*referral_code.ReferralCode, int64, error) {
	return s.referralCodeRepo.List(ctx, limit, offset)
}

// generateRandomCode generates a random 8-character alphanumeric code
func (s *ReferralCodeService) generateRandomCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 8

	code := make([]byte, codeLength)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}
	return string(code), nil
}

// GetOrGenerateUserReferralCodes gets existing user referral codes or generates codes to ensure user has exactly 5
func (s *ReferralCodeService) GetOrGenerateUserReferralCodes(ctx context.Context, userID uuid.UUID) ([]*referral_code.UserReferralCode, error) {
	// First check if user has activated referral feature
	user, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActivatedReferral {
		return nil, fmt.Errorf("user has not activated referral feature")
	}

	// Check if user already has referral codes
	existingCodes, err := s.referralCodeRepo.GetUserReferralCodesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user referral codes: %w", err)
	}

	// Calculate how many more codes we need to generate
	codesNeeded := 5 - len(existingCodes)

	// If user already has 5 or more codes, return existing ones
	if codesNeeded <= 0 {
		return existingCodes, nil
	}

	// Generate the remaining codes needed
	var newCodes []*referral_code.UserReferralCode
	for range codesNeeded {
		var code string
		var isUnique bool

		// Keep generating until we get a unique code
		for !isUnique {
			code, err = s.generateRandomCode()
			if err != nil {
				return nil, fmt.Errorf("failed to generate random code: %w", err)
			}

			// Only check if code exists in user referral codes (any user)
			_, err := s.referralCodeRepo.GetUserReferralCodeByCode(ctx, code)
			isUnique = err != nil // Code is unique if not found
		}

		// Create new user referral code
		newUserCode := &referral_code.UserReferralCode{
			UserID:       userID,
			ReferralCode: code,
			IsActivated:  false,
		}

		createdCode, err := s.referralCodeRepo.CreateUserReferralCode(ctx, newUserCode)
		if err != nil {
			return nil, fmt.Errorf("failed to create user referral code: %w", err)
		}

		newCodes = append(newCodes, createdCode)
	}

	// Return all codes (existing + newly created)
	allCodes := append(existingCodes, newCodes...)
	return allCodes, nil
}

// ActivateReferralCode activates a referral code and logs the usage
func (s *ReferralCodeService) ActivateReferralCode(ctx context.Context, userID uuid.UUID, code string) (*referral_code.ActivateReferralResponse, error) {
	// First check if user has already activated a referral code
	user, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.IsActivatedReferral {
		return &referral_code.ActivateReferralResponse{
			Success: false,
			Message: "User has already activated referral code",
		}, nil
	}

	// Check if it's a valid general referral code (unlimited use)
	isValidGeneral, err := s.referralCodeRepo.IsValidReferralCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to check general referral code: %w", err)
	}

	if isValidGeneral {
		// Update user's is_activated_referral to true
		user.IsActivatedReferral = true
		_, err = s.userRepo.Update(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to update user referral status: %w", err)
		}

		// Get the general referral code to obtain its ID
		generalReferralCode, err := s.referralCodeRepo.GetByReferralCode(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("failed to get general referral code: %w", err)
		}

		// Log the general referral code usage
		referralLog := &referral_code.ReferralLog{
			RedeemedByUserID: userID,
			CodeType:         "general",
			ReferralCodeID:   &generalReferralCode.ID,
		}

		_, err = s.referralCodeRepo.CreateReferralLog(ctx, referralLog)
		if err != nil {
			return nil, fmt.Errorf("failed to create referral log: %w", err)
		}

		return &referral_code.ActivateReferralResponse{
			Success: true,
			Message: "General referral code activated successfully",
		}, nil
	}

	// Check if it's a valid user referral code (single use)
	userReferralCode, err := s.referralCodeRepo.GetUserReferralCodeByCode(ctx, code)
	if err != nil {
		return &referral_code.ActivateReferralResponse{
			Success: false,
			Message: "Referral code is invalid",
		}, nil
	}

	// Check if the code is already activated
	if userReferralCode.IsActivated {
		return &referral_code.ActivateReferralResponse{
			Success: false,
			Message: "Referral code already activated",
		}, nil
	}

	// Update user's is_activated_referral to true
	user.IsActivatedReferral = true
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user referral status: %w", err)
	}

	// Activate the user referral code
	userReferralCode.IsActivated = true
	_, err = s.referralCodeRepo.UpdateUserReferralCode(ctx, userReferralCode)
	if err != nil {
		return nil, fmt.Errorf("failed to update user referral code: %w", err)
	}

	// Log the user referral code usage
	referralLog := &referral_code.ReferralLog{
		RedeemedByUserID: userID,
		CodeType:         "user",
		ReferralCodeID:   &userReferralCode.ID,
	}

	_, err = s.referralCodeRepo.CreateReferralLog(ctx, referralLog)
	if err != nil {
		return nil, fmt.Errorf("failed to create referral log: %w", err)
	}

	return &referral_code.ActivateReferralResponse{
		Success: true,
		Message: "User referral code activated successfully",
	}, nil
}
