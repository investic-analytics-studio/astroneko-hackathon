package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user"
	referralCodePorts "astroneko-backend/internal/core/ports/referral_code"
	userPorts "astroneko-backend/internal/core/ports/user"
	"astroneko-backend/pkg/firebase"
	"astroneko-backend/pkg/logger"

	"firebase.google.com/go/v4/auth"
)

type UserService struct {
	userRepo         userPorts.RepositoryInterface
	firebaseApp      firebase.FirebaseClientInterface
	firebaseAPIKey   string
	logger           logger.Logger
	referralCodeRepo referralCodePorts.RepositoryInterface
}

func NewUserService(userRepo userPorts.RepositoryInterface, firebaseApp firebase.FirebaseClientInterface, firebaseAPIKey string, log logger.Logger, referralCodeRepo referralCodePorts.RepositoryInterface) *UserService {
	return &UserService{
		userRepo:         userRepo,
		firebaseApp:      firebaseApp,
		firebaseAPIKey:   firebaseAPIKey,
		logger:           log,
		referralCodeRepo: referralCodeRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	// Creating new user - no logging needed for routine operation

	// Check if user already exists by email
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.logger.Warn("User already exists",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "email", Value: req.Email})
		return nil, shared.ErrUserEmailAlreadyExists
	}

	// Create user in Firebase first
	if s.firebaseApp == nil {
		s.logger.Error("Firebase not configured", logger.Field{Key: "module", Value: "user_service"})
		return nil, shared.ErrFirebaseNotInitialized
	}

	firebaseUserParams := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password).
		EmailVerified(false).
		Disabled(false)

	firebaseUser, err := s.firebaseApp.CreateUser(ctx, firebaseUserParams)
	if err != nil {
		s.logger.Error("Failed to create Firebase user",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "email", Value: req.Email},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrFirebaseUserCreationFailed
	}

	// Firebase user created - intermediate step, no logging needed

	// Create user in our database with Firebase UID
	newUser := &user.User{
		Email:               req.Email,
		IsActivatedReferral: false,
		FirebaseUID:         firebaseUser.UID,
	}

	dbUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		s.logger.Error("Failed to create user in database",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "firebase_uid", Value: firebaseUser.UID},
			logger.Field{Key: "error", Value: err.Error()})

		// If database creation fails, try to delete the Firebase user to maintain consistency
		if deleteErr := s.firebaseApp.DeleteUser(ctx, firebaseUser.UID); deleteErr != nil {
			s.logger.Error("Failed to cleanup Firebase user after database error",
				logger.Field{Key: "module", Value: "user_service"},
				logger.Field{Key: "firebase_uid", Value: firebaseUser.UID},
				logger.Field{Key: "error", Value: deleteErr.Error()})
			return nil, fmt.Errorf("failed to create user in database and failed to cleanup firebase user: %v (original error: %w)", deleteErr, err)
		}
		return nil, shared.ErrUserCreationFailed
	}

	s.logger.Info("User created successfully",
		logger.Field{Key: "module", Value: "user_service"},
		logger.Field{Key: "user_id", Value: dbUser.ID.String()},
		logger.Field{Key: "firebase_uid", Value: firebaseUser.UID},
		logger.Field{Key: "email", Value: req.Email})

	return dbUser, nil
}

func (s *UserService) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*user.User, error) {
	return s.userRepo.GetByFirebaseUID(ctx, firebaseUID)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req *user.UpdateUserRequest) (*user.User, error) {
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.IsActivatedReferral != nil {
		existingUser.IsActivatedReferral = *req.IsActivatedReferral
	}
	if req.LatestLoginAt != nil {
		existingUser.LatestLoginAt = req.LatestLoginAt
	}
	if req.ProfileImageURL != nil {
		existingUser.ProfileImageURL = req.ProfileImageURL
	}
	if req.DisplayName != nil {
		existingUser.DisplayName = req.DisplayName
	}

	return s.userRepo.Update(ctx, existingUser)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) VerifyFirebaseToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if s.firebaseApp == nil {
		s.logger.Error("Firebase not configured", logger.Field{Key: "module", Value: "user_service"})
		return nil, shared.ErrFirebaseNotInitialized
	}

	// Verifying token - no logging needed for routine operation

	token, err := s.firebaseApp.VerifyIDToken(ctx, idToken)
	if err != nil {
		s.logger.Error("Failed to verify Firebase token",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrFirebaseTokenVerifyFailed
	}

	// Token verified - no logging needed for routine success

	return token, nil
}

func (s *UserService) CreateUserFromToken(ctx context.Context, token *auth.Token) (*user.User, error) {
	email, ok := token.Claims["email"].(string)
	if !ok {
		s.logger.Error("Email not found in token",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "firebase_uid", Value: token.UID})
		return nil, shared.ErrEmailNotFoundInToken
	}

	// Creating user from token - no logging needed

	// Create user directly in database (user already exists in Firebase)
	newUser := &user.User{
		Email:               email,
		IsActivatedReferral: false,
		FirebaseUID:         token.UID,
	}

	dbUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		s.logger.Error("Failed to create user from token",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "firebase_uid", Value: token.UID},
			logger.Field{Key: "email", Value: email},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrUserCreationFailed
	}

	// User created from token - no logging needed

	return dbUser, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*user.AuthResponse, error) {
	// Refreshing token - no logging needed for routine operation

	// Call Firebase to refresh the token
	tokenResp, err := firebase.RefreshFirebaseToken(refreshToken)
	if err != nil {
		s.logger.Error("Failed to refresh Firebase token",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrFirebaseTokenRefreshFailed
	}

	// Verify the new ID token
	token, err := s.VerifyFirebaseToken(ctx, tokenResp.IDToken)
	if err != nil {
		s.logger.Error("Failed to verify refreshed token",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrFirebaseTokenVerifyFailed
	}

	// Get or create user
	existingUser, err := s.GetUserByFirebaseUID(ctx, token.UID)
	if err != nil {
		// User not found, creating from token - no logging needed

		// User doesn't exist, create new one
		newUser, createErr := s.CreateUserFromToken(ctx, token)
		if createErr != nil {
			s.logger.Error("Failed to create user from refreshed token",
				logger.Field{Key: "module", Value: "user_service"},
				logger.Field{Key: "firebase_uid", Value: token.UID},
				logger.Field{Key: "error", Value: createErr.Error()})
			return nil, shared.ErrUserCreationFailed
		}
		existingUser = newUser
	}

	// Convert expires_in string to int64
	expiresIn, _ := strconv.ParseInt(tokenResp.ExpiresIn, 10, 64)

	// Token refreshed - no logging needed

	return &user.AuthResponse{
		User:         existingUser.ToResponse(),
		AccessToken:  tokenResp.IDToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *UserService) GoogleAuth(ctx context.Context, req *user.GoogleLoginRequest) (*user.AuthResponse, error) {
	// Verify the Google ID token with Firebase (user already exists in Firebase from frontend)
	token, err := s.VerifyFirebaseToken(ctx, req.IDToken)
	if err != nil {
		s.logger.Error("Invalid Google token",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrGoogleTokenInvalid
	}

	// Extract email from token
	email, ok := token.Claims["email"].(string)
	if !ok || email == "" {
		s.logger.Error("Email not found in Google token",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "firebase_uid", Value: token.UID})
		return nil, shared.ErrGoogleTokenEmailMissing
	}

	// Get user info from Firebase
	userInfo, err := s.firebaseApp.GetUser(ctx, token.UID)
	if err != nil {
		s.logger.Error("Failed to get user info from Firebase",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "email", Value: email},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrFirebaseUserNotFound
	}

	// Extract profile information from Firebase user info
	var profileImageURL *string
	var displayName *string

	// Get profile image URL from Firebase user info
	if userInfo.PhotoURL != "" {
		profileImageURL = &userInfo.PhotoURL
	}

	// Get display name from Firebase user info
	if userInfo.DisplayName != "" {
		displayName = &userInfo.DisplayName
	}

	// Check if user already exists in our database by email
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		// Existing user found - no logging needed

		// User already exists, update login time and profile info
		now := time.Now()
		updateReq := &user.UpdateUserRequest{
			LatestLoginAt:   &now,
			ProfileImageURL: profileImageURL,
			DisplayName:     displayName,
		}
		updatedUser, updateErr := s.UpdateUser(ctx, existingUser.ID.String(), updateReq)
		if updateErr != nil {
			s.logger.Warn("Failed to update user profile for existing user",
				logger.Field{Key: "module", Value: "user_service"},
				logger.Field{Key: "user_id", Value: existingUser.ID.String()},
				logger.Field{Key: "error", Value: updateErr.Error()})
			updatedUser = existingUser
		}

		refreshToken := req.RefreshToken
		return &user.AuthResponse{
			User:         updatedUser.ToResponse(),
			AccessToken:  req.IDToken,
			RefreshToken: refreshToken,
			ExpiresIn:    3600, // Default 1 hour for Firebase tokens
		}, nil
	}

	// User doesn't exist in our database, create new one
	// Creating user from Google auth - no logging needed

	newUser := &user.User{
		Email:               email,
		IsActivatedReferral: false,
		FirebaseUID:         token.UID,
		ProfileImageURL:     profileImageURL,
		DisplayName:         displayName,
	}

	dbUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		s.logger.Error("Failed to create user from Google authentication",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "firebase_uid", Value: token.UID},
			logger.Field{Key: "email", Value: email},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrUserCreationFailed
	}

	// User created from Google auth - no logging needed

	refreshToken := req.RefreshToken
	if refreshToken == "" {
		refreshToken = "not_provided" // Placeholder for optional refresh token
	}
	return &user.AuthResponse{
		User:         dbUser.ToResponse(),
		AccessToken:  req.IDToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // Default 1 hour for Firebase tokens
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *user.LoginRequest) (*user.AuthResponse, error) {
	// Processing login - no logging needed for routine operation

	// Check if user exists in our database
	dbUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("User not found in database",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "email", Value: req.Email},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrUserNotFound
	}

	// Authenticate with Firebase using Admin SDK
	firebaseResp, err := firebase.SignInWithEmailPassword(req.Email, req.Password)
	if err != nil {
		s.logger.Error("Firebase authentication failed",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "email", Value: req.Email},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, shared.ErrFirebaseAuthFailed
	}

	// Verify that the Firebase user matches our database record
	if firebaseResp.LocalID != dbUser.FirebaseUID {
		s.logger.Error("Firebase UID mismatch",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "user_id", Value: dbUser.ID.String()},
			logger.Field{Key: "database_firebase_uid", Value: dbUser.FirebaseUID},
			logger.Field{Key: "firebase_response_uid", Value: firebaseResp.LocalID})
		return nil, shared.ErrFirebaseUIDMismatch
	}

	// Firebase auth successful - no logging needed

	// Update latest login time
	now := time.Now()
	updateReq := &user.UpdateUserRequest{
		LatestLoginAt: &now,
	}
	updatedUser, err := s.UpdateUser(ctx, dbUser.ID.String(), updateReq)
	if err != nil {
		// Log error but don't fail login
		s.logger.Warn("Failed to update login time for user",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "user_id", Value: dbUser.ID.String()},
			logger.Field{Key: "error", Value: err.Error()})
		updatedUser = dbUser
	}

	// Convert expires_in string to int64
	expiresIn, _ := strconv.ParseInt(firebaseResp.ExpiresIn, 10, 64)

	// Login successful - no logging needed for routine success

	return &user.AuthResponse{
		User:         updatedUser.ToResponse(),
		AccessToken:  firebaseResp.IDToken,
		RefreshToken: firebaseResp.RefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *UserService) GetTotalUsers(ctx context.Context) (int64, error) {
	totalUsers, err := s.userRepo.GetTotalUsers(ctx)
	if err != nil {
		s.logger.Error("Failed to get total users count",
			logger.Field{Key: "module", Value: "user_service"},
			logger.Field{Key: "error", Value: err.Error()})
		return 0, err
	}

	s.logger.Info("Total users count retrieved",
		logger.Field{Key: "module", Value: "user_service"},
		logger.Field{Key: "total_users", Value: totalUsers})

	return totalUsers, nil
}
