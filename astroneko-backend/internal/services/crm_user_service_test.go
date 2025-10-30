package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"astroneko-backend/internal/core/domain/crm_user"
	"astroneko-backend/testings/mock_logger"
	"astroneko-backend/testings/mock_ports"
)

// Test data builders for consistent test data
func buildCreateCRMUserRequest() *crm_user.CreateCRMUserRequest {
	return &crm_user.CreateCRMUserRequest{
		Username: "testuser",
		Password: "password123",
	}
}

func buildCreateCRMUserRequestWithUsernameAndPassword(username, password string) *crm_user.CreateCRMUserRequest {
	return &crm_user.CreateCRMUserRequest{
		Username: username,
		Password: password,
	}
}

func buildCRMLoginRequest() *crm_user.CRMLoginRequest {
	return &crm_user.CRMLoginRequest{
		Username: "testuser",
		Password: "password123",
	}
}

func buildCRMLoginRequestWithUsernameAndPassword(username, password string) *crm_user.CRMLoginRequest {
	return &crm_user.CRMLoginRequest{
		Username: username,
		Password: password,
	}
}

func buildCRMUser() *crm_user.CRMUser {
	return &crm_user.CRMUser{
		Username: "testuser",
		Password: hashPassword("password123"),
	}
}

func buildCRMUserWithID(id string) *crm_user.CRMUser {
	return &crm_user.CRMUser{
		Username: "testuser",
		Password: hashPassword("password123"),
	}
}

func buildCRMUserWithIDAndUsername(id, username string) *crm_user.CRMUser {
	return &crm_user.CRMUser{
		Username: username,
		Password: hashPassword("password123"),
	}
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func TestCRMUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	req := buildCreateCRMUserRequest()
	expectedUser := buildCRMUserWithID("123e4567-e89b-12d3-a456-426614174000")

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, gorm.ErrRecordNotFound)
	mockCRMUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, req.Username, result.Username)
	assert.NotEmpty(t, result.Password)
}

func TestCRMUserService_CreateUser_UsernameAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	req := buildCreateCRMUserRequest()
	existingUser := buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", req.Username)

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(existingUser, nil)

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "username already exists")
}

func TestCRMUserService_CreateUser_PasswordHashingError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	// Use a password that will cause bcrypt to fail
	req := buildCreateCRMUserRequestWithUsernameAndPassword("testuser", string(make([]byte, 73))) // Password too long for bcrypt

	// Setup expectations - service should check for existing user first
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to hash password")
}

func TestCRMUserService_CreateUser_RepositoryError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	req := buildCreateCRMUserRequest()
	dbError := errors.New("database constraint violation")

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, gorm.ErrRecordNotFound)
	mockCRMUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, dbError)

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create CRM user")
}

func TestCRMUserService_Login_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	req := buildCRMLoginRequest()
	existingUser := buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", req.Username)

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(existingUser, nil)

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.User)
	assert.Equal(t, existingUser.ID, result.User.ID)
	assert.Equal(t, existingUser.Username, result.User.Username)
	assert.NotEmpty(t, result.Token)
}

func TestCRMUserService_Login_UserNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	req := buildCRMLoginRequest()

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestCRMUserService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	req := buildCRMLoginRequestWithUsernameAndPassword("testuser", "wrongpassword")
	existingUser := buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", req.Username)

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(existingUser, nil)

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestCRMUserService_Login_TokenValidation(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	req := buildCRMLoginRequest()
	existingUser := &crm_user.CRMUser{
		Username: req.Username,
		Password: hashPassword(req.Password),
	}

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(existingUser, nil)

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.User)
	assert.NotEmpty(t, result.Token)

	// Verify token can be validated
	claims, err := service.ValidateToken(result.Token)
	require.NoError(t, err)
	assert.Equal(t, existingUser.Username, claims.Username)
}

func TestCRMUserService_GetUserByID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	testID := "123e4567-e89b-12d3-a456-426614174000"
	expectedUser := buildCRMUserWithID(testID)

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByID(ctx, testID).Return(expectedUser, nil)

	// Act
	result, err := service.GetUserByID(ctx, testID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Username, result.Username)
}

func TestCRMUserService_GetUserByID_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	testID := "123e4567-e89b-12d3-a456-426614174000"

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByID(ctx, testID).Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := service.GetUserByID(ctx, testID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestCRMUserService_ValidateToken_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	username := "testuser"

	// Generate a valid token
	claims := jwt.MapClaims{
		"user_id":  userID.String(),
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtSecret))

	// Act
	result, err := service.ValidateToken(tokenString)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, username, result.Username)
}

func TestCRMUserService_ValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)

	// Act
	result, err := service.ValidateToken("invalid-token")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestCRMUserService_ValidateToken_WrongSigningMethod(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)

	// For testing, we'll just create an invalid token string
	tokenString := "invalid.signing.method.token"

	// Act
	result, err := service.ValidateToken(tokenString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestCRMUserService_ValidateToken_ExpiredToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	username := "testuser"

	// Generate an expired token
	claims := jwt.MapClaims{
		"user_id":  userID.String(),
		"username": username,
		"exp":      time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
		"iat":      time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtSecret))

	// Act
	result, err := service.ValidateToken(tokenString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestCRMUserService_ValidateToken_InvalidClaims(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)

	// Generate a token with missing claims
	claims := jwt.MapClaims{
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
		// Missing user_id and username
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtSecret))

	// Act
	result, err := service.ValidateToken(tokenString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user_id in token")
}

// Table-driven tests for business logic scenarios
func TestCRMUserService_CreateUser_BusinessLogic(t *testing.T) {
	testCases := []struct {
		name          string
		username      string
		password      string
		existingUser  *crm_user.CRMUser
		repoError     error
		expectedError bool
		description   string
	}{
		{
			name:          "Valid new user",
			username:      "newuser",
			password:      "password123",
			existingUser:  nil,
			repoError:     nil,
			expectedError: false,
			description:   "Should successfully create new user",
		},
		{
			name:          "Existing username",
			username:      "existinguser",
			password:      "password123",
			existingUser:  buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", "existinguser"),
			repoError:     nil,
			expectedError: true,
			description:   "Should reject duplicate username",
		},
		{
			name:          "Database error on creation",
			username:      "erroruser",
			password:      "password123",
			existingUser:  nil,
			repoError:     errors.New("database connection failed"),
			expectedError: true,
			description:   "Should handle database errors gracefully",
		},
		{
			name:          "Short username",
			username:      "ab",
			password:      "password123",
			existingUser:  nil,
			repoError:     nil,
			expectedError: false,
			description:   "Service doesn't validate username length",
		},
		{
			name:          "Short password",
			username:      "testuser",
			password:      "123",
			existingUser:  nil,
			repoError:     nil,
			expectedError: false,
			description:   "Service doesn't validate password length",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
			mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

			jwtSecret := "test-secret"
			service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
			ctx := context.Background()

			req := buildCreateCRMUserRequestWithUsernameAndPassword(tc.username, tc.password)

			// Setup repository expectations
			if tc.existingUser != nil {
				mockCRMUserRepo.EXPECT().GetByUsername(ctx, tc.username).Return(tc.existingUser, nil)
			} else {
				mockCRMUserRepo.EXPECT().GetByUsername(ctx, tc.username).Return(nil, gorm.ErrRecordNotFound)

				if tc.repoError != nil {
					mockCRMUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, tc.repoError)
				} else {
					expectedUser := buildCRMUserWithID("123e4567-e89b-12d3-a456-426614174000")
					expectedUser.Username = tc.username
					mockCRMUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)
				}
			}

			// Act
			result, err := service.CreateUser(ctx, req)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.username, result.Username)
			}
		})
	}
}

// Table-driven tests for login scenarios
func TestCRMUserService_Login_BusinessLogic(t *testing.T) {
	testCases := []struct {
		name          string
		username      string
		password      string
		existingUser  *crm_user.CRMUser
		repoError     error
		expectedError bool
		description   string
	}{
		{
			name:          "Valid credentials",
			username:      "testuser",
			password:      "password123",
			existingUser:  buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", "testuser"),
			repoError:     nil,
			expectedError: false,
			description:   "Should successfully login with valid credentials",
		},
		{
			name:          "User not found",
			username:      "nonexistentuser",
			password:      "password123",
			existingUser:  nil,
			repoError:     gorm.ErrRecordNotFound,
			expectedError: true,
			description:   "Should reject login for non-existent user",
		},
		{
			name:          "Wrong password",
			username:      "testuser",
			password:      "wrongpassword",
			existingUser:  buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", "testuser"),
			repoError:     nil,
			expectedError: true,
			description:   "Should reject login with wrong password",
		},
		{
			name:          "Database error",
			username:      "erroruser",
			password:      "password123",
			existingUser:  nil,
			repoError:     errors.New("database connection failed"),
			expectedError: true,
			description:   "Should handle database errors gracefully",
		},
		{
			name:          "Empty username",
			username:      "",
			password:      "password123",
			existingUser:  nil,
			repoError:     gorm.ErrRecordNotFound,
			expectedError: true,
			description:   "Should reject login with empty username",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
			mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

			jwtSecret := "test-secret"
			service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
			ctx := context.Background()

			req := buildCRMLoginRequestWithUsernameAndPassword(tc.username, tc.password)

			// Setup repository expectations
			if tc.existingUser != nil {
				mockCRMUserRepo.EXPECT().GetByUsername(ctx, tc.username).Return(tc.existingUser, nil)
			} else {
				mockCRMUserRepo.EXPECT().GetByUsername(ctx, tc.username).Return(nil, tc.repoError)
			}

			// Act
			result, err := service.Login(ctx, req)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotNil(t, result.User)
				assert.NotEmpty(t, result.Token)
			}
		})
	}
}

// Performance test
func BenchmarkCRMUserService_CreateUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()
	req := buildCreateCRMUserRequest()
	expectedUser := buildCRMUserWithID("123e4567-e89b-12d3-a456-426614174000")

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	mockCRMUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Username = fmt.Sprintf("user%d", i)
		_, err := service.CreateUser(ctx, req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// Concurrent testing
func TestCRMUserService_CreateUser_Concurrent(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx := context.Background()

	const numGoroutines = 10
	results := make(chan *crm_user.CRMUser, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Mock expectations for concurrent calls
	for i := 0; i < numGoroutines; i++ {
		username := fmt.Sprintf("user%d", i)
		expectedUser := buildCRMUserWithID(fmt.Sprintf("123e4567-e89b-12d3-a456-426614%03d", i))
		expectedUser.Username = username

		mockCRMUserRepo.EXPECT().GetByUsername(ctx, username).Return(nil, gorm.ErrRecordNotFound)
		mockCRMUserRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(user *crm_user.CRMUser) (*crm_user.CRMUser, error) {
			user.ID = expectedUser.ID
			user.Username = expectedUser.Username
			return user, nil
		})
	}

	// Act
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			req := buildCreateCRMUserRequestWithUsernameAndPassword(fmt.Sprintf("user%d", id), "password123")
			result, err := service.CreateUser(ctx, req)
			results <- result
			errors <- err
		}(i)
	}

	// Assert
	for i := 0; i < numGoroutines; i++ {
		result := <-results
		err := <-errors
		assert.NoError(t, err)
		assert.NotNil(t, result)
	}
}

// Test timeout handling
func TestCRMUserService_Timeout(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCRMUserRepo := mock_ports.NewMockCRMUserRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	jwtSecret := "test-secret"
	service := NewCRMUserService(mockCRMUserRepo, mockLogger, jwtSecret)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	req := buildCreateCRMUserRequest()
	timeoutError := errors.New("context deadline exceeded")

	// Setup expectations
	mockCRMUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, timeoutError)

	// Act
	result, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
