package crm_user

import (
	"astroneko-backend/internal/core/domain/shared"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CRMUser struct {
	shared.NoDeletedModel
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"-" gorm:"not null"`
}

func (CRMUser) TableName() string {
	return "astroneko_crm_users"
}

func (u *CRMUser) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *CRMUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type CreateCRMUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type CRMLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CRMUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (u *CRMUser) ToResponse() *CRMUserResponse {
	return &CRMUserResponse{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type CRMLoginResponse struct {
	User  *CRMUserResponse `json:"user"`
	Token string           `json:"token"`
}

type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}
