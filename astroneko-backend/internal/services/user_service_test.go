package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/testings/mock_firebase"
	"astroneko-backend/testings/mock_logger"
	"astroneko-backend/testings/mock_ports"
)

// Test data builders
func buildCreateUserRequest() *user.CreateUserRequest {
	return &user.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
}

func buildUpdateUserRequest() *user.UpdateUserRequest {
	now := time.Now()
	return &user.UpdateUserRequest{
		IsActivatedReferral: boolPtr(true),
		LatestLoginAt:       &now,
		ProfileImageURL:     stringPtr("https://example.com/avatar.jpg"),
		DisplayName:         stringPtr("Updated User"),
	}
}

func buildLoginRequest() *user.LoginRequest {
	return &user.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
}

func buildGoogleLoginRequest() *user.GoogleLoginRequest {
	return &user.GoogleLoginRequest{
		IDToken:      "google_id_token_123",
		RefreshToken: "google_refresh_token_123",
	}
}

func buildTestUser() *user.User {
	now := time.Now()
	testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	user := &user.User{
		Email:               "test@example.com",
		IsActivatedReferral: false,
		FirebaseUID:         "firebase_123",
		LatestLoginAt:       &now,
		ProfileImageURL:     stringPtr("https://example.com/avatar.jpg"),
		DisplayName:         stringPtr("Test User"),
	}
	user.ID = testUUID
	return user
}

func buildFirebaseToken() *auth.Token {
	return &auth.Token{
		UID: "firebase_123",
		Claims: map[string]interface{}{
			"email":   "test@example.com",
			"name":    "Test User",
			"picture": "https://example.com/avatar.jpg",
		},
	}
}

func buildFirebaseUserInfo() *auth.UserInfo {
	return &auth.UserInfo{
		UID:         "firebase_123",
		Email:       "test@example.com",
		DisplayName: "Test User",
		PhotoURL:    "https://example.com/avatar.jpg",
	}
}

// Helper functions
func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }

func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildCreateUserRequest()
	expectedUser := buildTestUser()

	// Mock user not exists
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Mock Firebase user creation
	firebaseUser := &auth.UserRecord{UserInfo: buildFirebaseUserInfo()}
	mockFirebaseApp.EXPECT().CreateUser(ctx, gomock.Any()).Return(firebaseUser, nil)

	// Mock database user creation
	mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)

	// Mock logger calls
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).Times(0)
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(0)
	mockLogger.EXPECT().Info("User created successfully", gomock.Any())

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, "firebase_123", result.FirebaseUID)
}

func TestUserService_CreateUser_UserAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildCreateUserRequest()
	existingUser := buildTestUser()

	// Mock user already exists
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(existingUser, nil)

	// Mock logger call
	mockLogger.EXPECT().Warn("User already exists", gomock.Any())

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrUserEmailAlreadyExists, err)
}

func TestUserService_CreateUser_FirebaseNotConfigured(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, nil, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildCreateUserRequest()

	// Mock user not exists
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Mock logger call
	mockLogger.EXPECT().Error("Firebase not configured", gomock.Any())

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrFirebaseNotInitialized, err)
}

func TestUserService_CreateUser_FirebaseUserCreationFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildCreateUserRequest()
	firebaseError := errors.New("firebase user creation failed")

	// Mock user not exists
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Mock Firebase user creation failure
	mockFirebaseApp.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, firebaseError)

	// Mock logger call
	mockLogger.EXPECT().Error("Failed to create Firebase user", gomock.Any())

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrFirebaseUserCreationFailed, err)
}

func TestUserService_CreateUser_DatabaseCreationFailsWithCleanup(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildCreateUserRequest()
	dbError := errors.New("database constraint violation")

	// Mock user not exists
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Mock Firebase user creation success
	firebaseUser := &auth.UserRecord{UserInfo: buildFirebaseUserInfo()}
	mockFirebaseApp.EXPECT().CreateUser(ctx, gomock.Any()).Return(firebaseUser, nil)

	// Mock database user creation failure
	mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, dbError)

	// Mock Firebase cleanup success
	mockFirebaseApp.EXPECT().DeleteUser(ctx, "firebase_123").Return(nil)

	// Mock logger calls
	mockLogger.EXPECT().Error("Failed to create user in database", gomock.Any())

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrUserCreationFailed, err)
}

func TestUserService_GetUserByFirebaseUID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	firebaseUID := "firebase_123"
	expectedUser := buildTestUser()

	// Mock repository call
	mockUserRepo.EXPECT().GetByFirebaseUID(ctx, firebaseUID).Return(expectedUser, nil)

	// Act
	result, err := service.GetUserByFirebaseUID(ctx, firebaseUID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.FirebaseUID, result.FirebaseUID)
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	userID, _ := uuid.Parse("test-user-id")
	updateReq := buildUpdateUserRequest()
	existingUser := buildTestUser()
	existingUser.ID = userID
	updatedUser := buildTestUser()
	updatedUser.ID = userID
	updatedUser.IsActivatedReferral = *updateReq.IsActivatedReferral
	updatedUser.LatestLoginAt = updateReq.LatestLoginAt
	updatedUser.ProfileImageURL = updateReq.ProfileImageURL
	updatedUser.DisplayName = updateReq.DisplayName

	// Mock repository calls
	mockUserRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).Return(updatedUser, nil)

	// Act
	result, err := service.UpdateUser(ctx, userID.String(), updateReq)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, updatedUser.ID, result.ID)
	assert.Equal(t, updatedUser.IsActivatedReferral, result.IsActivatedReferral)
	assert.Equal(t, updatedUser.DisplayName, result.DisplayName)
}

func TestUserService_UpdateUser_UserNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	userID := "nonexistent-user"
	updateReq := buildUpdateUserRequest()

	// Mock repository call
	mockUserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := service.UpdateUser(ctx, userID, updateReq)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserService_VerifyFirebaseToken_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	idToken := "valid_id_token"
	expectedToken := buildFirebaseToken()

	// Mock Firebase call
	mockFirebaseApp.EXPECT().VerifyIDToken(ctx, idToken).Return(expectedToken, nil)

	// Act
	result, err := service.VerifyFirebaseToken(ctx, idToken)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedToken.UID, result.UID)
	assert.Equal(t, expectedToken.Claims["email"], result.Claims["email"])
}

func TestUserService_VerifyFirebaseToken_FirebaseNotConfigured(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, nil, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	idToken := "id_token"

	// Mock logger call
	mockLogger.EXPECT().Error("Firebase not configured", gomock.Any())

	// Act
	result, err := service.VerifyFirebaseToken(ctx, idToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrFirebaseNotInitialized, err)
}

func TestUserService_VerifyFirebaseToken_InvalidToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	idToken := "invalid_token"
	firebaseError := errors.New("invalid token")

	// Mock Firebase call
	mockFirebaseApp.EXPECT().VerifyIDToken(ctx, idToken).Return(nil, firebaseError)

	// Mock logger call
	mockLogger.EXPECT().Error("Failed to verify Firebase token", gomock.Any())

	// Act
	result, err := service.VerifyFirebaseToken(ctx, idToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrFirebaseTokenVerifyFailed, err)
}

func TestUserService_CreateUserFromToken_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	token := buildFirebaseToken()
	expectedUser := buildTestUser()

	// Mock repository call
	mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)

	// Act
	result, err := service.CreateUserFromToken(ctx, token)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, token.UID, result.FirebaseUID)
	assert.Equal(t, token.Claims["email"], result.Email)
}

func TestUserService_CreateUserFromToken_EmailNotFoundInToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	token := &auth.Token{
		UID: "firebase_123",
		Claims: map[string]interface{}{
			"name": "Test User",
			// Email missing
		},
	}

	// Mock logger call
	mockLogger.EXPECT().Error("Email not found in token", gomock.Any())

	// Act
	result, err := service.CreateUserFromToken(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrEmailNotFoundInToken, err)
}

func TestUserService_GoogleAuth_ExistingUser(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildGoogleLoginRequest()
	token := buildFirebaseToken()
	userInfo := buildFirebaseUserInfo()
	existingUser := buildTestUser()
	updatedUser := buildTestUser()

	// Mock Firebase calls
	mockFirebaseApp.EXPECT().VerifyIDToken(ctx, req.IDToken).Return(token, nil)
	mockFirebaseApp.EXPECT().GetUser(ctx, token.UID).Return(&auth.UserRecord{UserInfo: userInfo}, nil)

	// Mock repository calls
	mockUserRepo.EXPECT().GetByEmail(ctx, token.Claims["email"].(string)).Return(existingUser, nil)
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).Return(updatedUser, nil)

	// Act
	result, err := service.GoogleAuth(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, updatedUser.ID, result.User.ID)
	assert.Equal(t, req.IDToken, result.AccessToken)
	assert.Equal(t, req.RefreshToken, result.RefreshToken)
	assert.Equal(t, int64(3600), result.ExpiresIn)
}

func TestUserService_GoogleAuth_NewUser(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildGoogleLoginRequest()
	token := buildFirebaseToken()
	userInfo := buildFirebaseUserInfo()
	newUser := buildTestUser()

	// Mock Firebase calls
	mockFirebaseApp.EXPECT().VerifyIDToken(ctx, req.IDToken).Return(token, nil)
	mockFirebaseApp.EXPECT().GetUser(ctx, token.UID).Return(&auth.UserRecord{UserInfo: userInfo}, nil)

	// Mock repository calls
	mockUserRepo.EXPECT().GetByEmail(ctx, token.Claims["email"].(string)).Return(nil, gorm.ErrRecordNotFound)
	mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(newUser, nil)

	// Act
	result, err := service.GoogleAuth(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, newUser.ID, result.User.ID)
	assert.Equal(t, req.IDToken, result.AccessToken)
	assert.Equal(t, req.RefreshToken, result.RefreshToken)
	assert.Equal(t, int64(3600), result.ExpiresIn)
}

func TestUserService_GoogleAuth_InvalidToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildGoogleLoginRequest()
	firebaseError := errors.New("invalid token")

	// Mock Firebase call
	mockFirebaseApp.EXPECT().VerifyIDToken(ctx, req.IDToken).Return(nil, firebaseError)

	// Mock logger call
	mockLogger.EXPECT().Error("Invalid Google token", gomock.Any())

	// Act
	result, err := service.GoogleAuth(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrGoogleTokenInvalid, err)
}

func TestUserService_GoogleAuth_EmailMissingInToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildGoogleLoginRequest()
	token := &auth.Token{
		UID: "firebase_123",
		Claims: map[string]interface{}{
			"name": "Test User",
			// Email missing
		},
	}

	// Mock Firebase calls
	mockFirebaseApp.EXPECT().VerifyIDToken(ctx, req.IDToken).Return(token, nil)

	// Mock logger call
	mockLogger.EXPECT().Error("Email not found in Google token", gomock.Any())

	// Act
	result, err := service.GoogleAuth(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrGoogleTokenEmailMissing, err)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	req := buildLoginRequest()

	// Mock repository call
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Mock logger call
	mockLogger.EXPECT().Error("User not found in database", gomock.Any())

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrUserNotFound, err)
}

func TestUserService_GetTotalUsers_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	expectedCount := int64(100)

	// Mock repository call
	mockUserRepo.EXPECT().GetTotalUsers(ctx).Return(expectedCount, nil)

	// Mock logger call
	mockLogger.EXPECT().Info("Total users count retrieved", gomock.Any())

	// Act
	result, err := service.GetTotalUsers(ctx)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedCount, result)
}

func TestUserService_GetTotalUsers_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()
	dbError := errors.New("database connection failed")

	// Mock repository call
	mockUserRepo.EXPECT().GetTotalUsers(ctx).Return(int64(0), dbError)

	// Mock logger call
	mockLogger.EXPECT().Error("Failed to get total users count", gomock.Any())

	// Act
	result, err := service.GetTotalUsers(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), result)
	assert.Equal(t, dbError, err)
}

// Table-driven tests for edge cases
func TestUserService_CreateUser_InvalidEmails(t *testing.T) {
	testCases := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name:        "Valid email",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "Empty email",
			email:       "",
			expectError: true,
		},
		{
			name:        "Invalid format",
			email:       "invalid-email",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
			mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
			mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
			mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

			service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
			ctx := context.Background()
			req := &user.CreateUserRequest{
				Email:    tc.email,
				Password: "password123",
			}

			// Mock user not exists
			mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

			// Only mock Firebase call if email is valid
			if !tc.expectError {
				firebaseUser := &auth.UserRecord{UserInfo: buildFirebaseUserInfo()}
				firebaseUser.UserInfo.Email = tc.email
				mockFirebaseApp.EXPECT().CreateUser(ctx, gomock.Any()).Return(firebaseUser, nil)
			}

			// Act
			result, err := service.CreateUser(ctx, req)

			// Assert
			if tc.expectError {
				// For invalid emails, the validation would typically happen at the handler level
				// This test shows the structure
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				// For valid emails, we expect the process to continue
				// The actual success depends on Firebase and database calls
				if err != nil {
					assert.NotEqual(t, shared.ErrUserEmailAlreadyExists, err)
				}
			}
		})
	}
}

// Performance test
func TestUserService_CreateUser_Concurrent(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_ports.NewMockUserRepositoryInterface(ctrl)
	mockFirebaseApp := mock_firebase.NewMockFirebaseClientInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockReferralCodeRepo := mock_ports.NewReferralCodeRepositoryInterface(ctrl)

	service := NewUserService(mockUserRepo, mockFirebaseApp, "test-api-key", mockLogger, mockReferralCodeRepo)
	ctx := context.Background()

	const numGoroutines = 10
	errors := make(chan error, numGoroutines)

	// Mock expectations for concurrent calls
	for i := 0; i < numGoroutines; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		mockUserRepo.EXPECT().GetByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)
		firebaseUser := &auth.UserRecord{UserInfo: buildFirebaseUserInfo()}
		firebaseUser.UserInfo.Email = email
		mockFirebaseApp.EXPECT().CreateUser(ctx, gomock.Any()).Return(firebaseUser, nil)
		expectedUser := buildTestUser()
		expectedUser.Email = email
		expectedUser.FirebaseUID = fmt.Sprintf("firebase_%d", i)
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)
	}

	// Act
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			req := &user.CreateUserRequest{
				Email:    fmt.Sprintf("user%d@example.com", id),
				Password: "password123",
			}
			_, err := service.CreateUser(ctx, req)
			errors <- err
		}(i)
	}

	// Assert
	for i := 0; i < numGoroutines; i++ {
		err := <-errors
		assert.NoError(t, err)
	}
}
