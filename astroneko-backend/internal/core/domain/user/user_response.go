package user

import (
	"time"

	"astroneko-backend/internal/core/domain/shared"
)

type UserResponse struct {
	ID                  string     `json:"id"`
	Email               string     `json:"email"`
	IsActivatedReferral bool       `json:"is_activated_referral"`
	LatestLoginAt       *time.Time `json:"latest_login_at"`
	FirebaseUID         string     `json:"firebase_uid"`
	ProfileImageURL     *string    `json:"profile_image_url"`
	DisplayName         *string    `json:"display_name"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:                  u.ID.String(),
		Email:               u.Email,
		IsActivatedReferral: u.IsActivatedReferral,
		LatestLoginAt:       u.LatestLoginAt,
		FirebaseUID:         u.FirebaseUID,
		ProfileImageURL:     u.ProfileImageURL,
		DisplayName:         u.DisplayName,
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}
}

type CreateUserResponse struct {
	shared.ResponseBody
}

type GetUserResponse struct {
	shared.ResponseBody
}

type UpdateUserResponse struct {
	shared.ResponseBody
}

type DeleteUserResponse struct {
	shared.ResponseBody
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
}

type RefreshTokenResponse struct {
	shared.ResponseBody
}

type ActivateReferralResponse struct {
	shared.ResponseBody
}

type GetTotalUsersResponse struct {
	TotalUsers int64 `json:"total_users"`
}
