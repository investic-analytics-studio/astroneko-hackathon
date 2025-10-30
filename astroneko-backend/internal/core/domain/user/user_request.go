package user

import "time"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	IsActivatedReferral *bool      `json:"is_activated_referral"`
	LatestLoginAt       *time.Time `json:"latest_login_at"`
	ProfileImageURL     *string    `json:"profile_image_url"`
	DisplayName         *string    `json:"display_name"`
}

type GetUserByFirebaseUIDRequest struct {
	FirebaseUID string `json:"firebase_uid" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type GoogleLoginRequest struct {
	IDToken      string `json:"id_token" validate:"required"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ActivateReferralRequest struct {
	ReferralCode string `json:"referral_code" validate:"required"`
}
