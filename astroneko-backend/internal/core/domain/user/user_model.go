package user

import (
	"time"

	"astroneko-backend/internal/core/domain/shared"
)

type User struct {
	shared.NoDeletedModel
	Email               string     `json:"email" gorm:"uniqueIndex;not null"`
	IsActivatedReferral bool       `json:"is_activated_referral" gorm:"default:false;not null"`
	LatestLoginAt       *time.Time `json:"latest_login_at"`
	FirebaseUID         string     `json:"firebase_uid" gorm:"not null"`
	ProfileImageURL     *string    `json:"profile_image_url"`
	DisplayName         *string    `json:"display_name"`
}

func (User) TableName() string {
	return "astroneko_auth_users"
}
