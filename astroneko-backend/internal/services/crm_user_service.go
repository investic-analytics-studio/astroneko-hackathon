package services

import (
	"context"
	"fmt"
	"time"

	"astroneko-backend/internal/core/domain/crm_user"
	crmUserPorts "astroneko-backend/internal/core/ports/crm_user"
	"astroneko-backend/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CRMUserService struct {
	crmUserRepo crmUserPorts.RepositoryInterface
	logger      logger.Logger
	jwtSecret   string
}

func NewCRMUserService(crmUserRepo crmUserPorts.RepositoryInterface, logger logger.Logger, jwtSecret string) *CRMUserService {
	return &CRMUserService{
		crmUserRepo: crmUserRepo,
		logger:      logger,
		jwtSecret:   jwtSecret,
	}
}

func (s *CRMUserService) CreateUser(ctx context.Context, req *crm_user.CreateCRMUserRequest) (*crm_user.CRMUser, error) {
	existingUser, err := s.crmUserRepo.GetByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	newUser := &crm_user.CRMUser{
		Username: req.Username,
	}

	if err := newUser.HashPassword(req.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	createdUser, err := s.crmUserRepo.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create CRM user: %w", err)
	}

	return createdUser, nil
}

func (s *CRMUserService) Login(ctx context.Context, req *crm_user.CRMLoginRequest) (*crm_user.CRMLoginResponse, error) {
	user, err := s.crmUserRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.CheckPassword(req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &crm_user.CRMLoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *CRMUserService) GetUserByID(ctx context.Context, id string) (*crm_user.CRMUser, error) {
	return s.crmUserRepo.GetByID(ctx, id)
}

func (s *CRMUserService) ValidateToken(tokenString string) (*crm_user.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format in token")
	}

	username, ok := (*claims)["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username in token")
	}

	return &crm_user.JWTClaims{
		UserID:   userID,
		Username: username,
	}, nil
}

func (s *CRMUserService) generateJWT(user *crm_user.CRMUser) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
