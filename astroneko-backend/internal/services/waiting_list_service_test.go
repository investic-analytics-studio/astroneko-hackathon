package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/waiting_list"
	"astroneko-backend/testings/mock_logger"
	"astroneko-backend/testings/mock_ports"
)

// Test data builders for consistent test data
func buildWaitingListUser() *waiting_list.WaitingListUser {
	return &waiting_list.WaitingListUser{
		Email: "test@example.com",
	}
}

func buildWaitingListUserWithID(id string) *waiting_list.WaitingListUser {
	return &waiting_list.WaitingListUser{
		Email: "test@example.com",
	}
}

func buildWaitingListUserWithIDAndEmail(id, email string) *waiting_list.WaitingListUser {
	return &waiting_list.WaitingListUser{
		Email: email,
	}
}

func TestWaitingListService_JoinWaitingList_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	email := "test@example.com"
	expectedUser := buildWaitingListUserWithID("123e4567-e89b-12d3-a456-426614174000")

	// Setup expectations - user not found, then creation succeeds
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)
	mockWaitingListRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)

	// Mock logger calls
	mockLogger.EXPECT().Warn("User already exists in waiting list", gomock.Any()).Times(0) // Should not be called
	mockLogger.EXPECT().Error("Failed to add user to waiting list", gomock.Any()).Times(0) // Should not be called
	mockLogger.EXPECT().Info("User successfully added to waiting list", gomock.Any())

	// Act
	result, err := service.JoinWaitingList(ctx, email)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, email, result.Email)
}

func TestWaitingListService_JoinWaitingList_UserAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	email := "existing@example.com"
	existingUser := buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174000", email)

	// Setup expectations - user already exists
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(existingUser, nil)

	// Mock logger calls
	mockLogger.EXPECT().Warn("User already exists in waiting list", gomock.Any())
	mockLogger.EXPECT().Error("Failed to add user to waiting list", gomock.Any()).Times(0)     // Should not be called
	mockLogger.EXPECT().Info("User successfully added to waiting list", gomock.Any()).Times(0) // Should not be called

	// Act
	result, err := service.JoinWaitingList(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrWaitingListUserAlreadyExists, err)
}

func TestWaitingListService_JoinWaitingList_CreationFailed(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	email := "fail@example.com"
	dbError := errors.New("database constraint violation")

	// Setup expectations - user not found, but creation fails
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)
	mockWaitingListRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, dbError)

	// Mock logger calls
	mockLogger.EXPECT().Warn("User already exists in waiting list", gomock.Any()).Times(0) // Should not be called
	mockLogger.EXPECT().Error("Failed to add user to waiting list", gomock.Any())
	mockLogger.EXPECT().Info("User successfully added to waiting list", gomock.Any()).Times(0) // Should not be called

	// Act
	result, err := service.JoinWaitingList(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrWaitingListUserCreationFailed, err)
}

func TestWaitingListService_JoinWaitingList_EmptyEmail(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	email := ""
	expectedUser := &waiting_list.WaitingListUser{Email: email}

	// Setup expectations - empty email should be processed normally
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)
	mockWaitingListRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)

	// Mock logger calls
	mockLogger.EXPECT().Info("User successfully added to waiting list", gomock.Any())

	// Act
	result, err := service.JoinWaitingList(ctx, email)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, email, result.Email)
}

func TestWaitingListService_GetWaitingListUsers_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	limit := 10
	offset := 0
	expectedUsers := []*waiting_list.WaitingListUser{
		buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174000", "user1@example.com"),
		buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174001", "user2@example.com"),
	}
	expectedCount := int64(2)

	// Setup expectations
	mockWaitingListRepo.EXPECT().List(ctx, limit, offset).Return(expectedUsers, expectedCount, nil)

	// Act
	users, count, err := service.GetWaitingListUsers(ctx, limit, offset)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.Len(t, users, len(expectedUsers))
	assert.Equal(t, expectedUsers[0].Email, users[0].Email)
	assert.Equal(t, expectedUsers[1].Email, users[1].Email)
}

func TestWaitingListService_GetWaitingListUsers_RepositoryError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	limit := 10
	offset := 0
	dbError := errors.New("database connection failed")

	// Setup expectations
	mockWaitingListRepo.EXPECT().List(ctx, limit, offset).Return(nil, int64(0), dbError)

	// Act
	users, count, err := service.GetWaitingListUsers(ctx, limit, offset)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.Nil(t, users)
	assert.Equal(t, dbError, err)
}

func TestWaitingListService_GetWaitingListUsers_Pagination(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()

	testCases := []struct {
		name          string
		limit         int
		offset        int
		expectedUsers []*waiting_list.WaitingListUser
		expectedCount int64
		description   string
	}{
		{
			name:          "First page",
			limit:         5,
			offset:        0,
			expectedUsers: []*waiting_list.WaitingListUser{buildWaitingListUserWithID("123e4567-e89b-12d3-a456-426614174000")},
			expectedCount: int64(15),
			description:   "Should return first 5 users with total count",
		},
		{
			name:          "Second page",
			limit:         5,
			offset:        5,
			expectedUsers: []*waiting_list.WaitingListUser{buildWaitingListUserWithID("123e4567-e89b-12d3-a456-426614174001")},
			expectedCount: int64(15),
			description:   "Should return second 5 users with total count",
		},
		{
			name:          "Empty page",
			limit:         10,
			offset:        100,
			expectedUsers: []*waiting_list.WaitingListUser{},
			expectedCount: int64(15),
			description:   "Should return empty list for offset beyond data",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup expectations
			mockWaitingListRepo.EXPECT().List(ctx, tc.limit, tc.offset).Return(tc.expectedUsers, tc.expectedCount, nil)

			// Act
			users, count, err := service.GetWaitingListUsers(ctx, tc.limit, tc.offset)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCount, count)
			assert.Len(t, users, len(tc.expectedUsers))
		})
	}
}

func TestWaitingListService_GetWaitingListUserByEmail_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	email := "test@example.com"
	expectedUser := buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174000", email)

	// Setup expectations
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(expectedUser, nil)

	// Act
	result, err := service.GetWaitingListUserByEmail(ctx, email)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)
}

func TestWaitingListService_GetWaitingListUserByEmail_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	email := "nonexistent@example.com"

	// Setup expectations
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := service.GetWaitingListUserByEmail(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestWaitingListService_IsInWaitingListByEmail_UserExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	email := "existing@example.com"
	existingUser := buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174000", email)

	// Setup expectations
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(existingUser, nil)

	// Act
	isInList, err := service.IsInWaitingListByEmail(ctx, email)

	// Assert
	require.NoError(t, err)
	assert.True(t, isInList)
}

func TestWaitingListService_IsInWaitingListByEmail_UserNotExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	email := "nonexistent@example.com"

	// Setup expectations
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)

	// Act
	isInList, err := service.IsInWaitingListByEmail(ctx, email)

	// Assert
	require.NoError(t, err)
	assert.False(t, isInList)
}

func TestWaitingListService_IsInWaitingListByEmail_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	email := "error@example.com"
	dbError := errors.New("database connection failed")

	// Setup expectations
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(nil, dbError)

	// Act
	isInList, err := service.IsInWaitingListByEmail(ctx, email)

	// Assert
	require.NoError(t, err) // Service should handle error and return false
	assert.False(t, isInList)
}

// Table-driven tests for business logic scenarios
func TestWaitingListService_JoinWaitingList_BusinessLogic(t *testing.T) {
	testCases := []struct {
		name           string
		email          string
		existingUser   *waiting_list.WaitingListUser
		repoError      error
		expectedError  error
		expectedStatus string
		description    string
	}{
		{
			name:           "Valid new user",
			email:          "newuser@example.com",
			existingUser:   nil,
			repoError:      nil,
			expectedError:  nil,
			expectedStatus: "success",
			description:    "Should successfully add new user to waiting list",
		},
		{
			name:           "Existing user",
			email:          "existing@example.com",
			existingUser:   buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174000", "existing@example.com"),
			repoError:      nil,
			expectedError:  shared.ErrWaitingListUserAlreadyExists,
			expectedStatus: "",
			description:    "Should reject duplicate email",
		},
		{
			name:           "Database error on check",
			email:          "error@example.com",
			existingUser:   nil,
			repoError:      errors.New("connection failed"),
			expectedError:  shared.ErrWaitingListUserCreationFailed,
			expectedStatus: "",
			description:    "Should handle database errors gracefully",
		},
		{
			name:           "Empty email",
			email:          "",
			existingUser:   nil,
			repoError:      nil,
			expectedError:  nil,
			expectedStatus: "success",
			description:    "Should handle empty email (validation not in service)",
		},
		{
			name:           "Invalid email format",
			email:          "invalid-email",
			existingUser:   nil,
			repoError:      nil,
			expectedError:  nil,
			expectedStatus: "success",
			description:    "Should handle invalid email format (validation not in service)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
			mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

			service := NewWaitingListService(mockWaitingListRepo, mockLogger)
			ctx := context.Background()

			// Setup repository expectations
			if tc.existingUser != nil {
				mockWaitingListRepo.EXPECT().GetByEmail(ctx, tc.email).Return(tc.existingUser, nil)
				mockLogger.EXPECT().Warn("User already exists in waiting list", gomock.Any())
			} else if tc.repoError != nil {
				if tc.repoError.Error() == "connection failed" {
					mockWaitingListRepo.EXPECT().GetByEmail(ctx, tc.email).Return(nil, tc.repoError)
					mockWaitingListRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, tc.repoError)
					mockLogger.EXPECT().Error("Failed to add user to waiting list", gomock.Any())
				}
			} else {
				mockWaitingListRepo.EXPECT().GetByEmail(ctx, tc.email).Return(nil, gorm.ErrRecordNotFound)
				expectedUser := &waiting_list.WaitingListUser{Email: tc.email}
				mockWaitingListRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)
				mockLogger.EXPECT().Info("User successfully added to waiting list", gomock.Any())
			}

			// Act
			result, err := service.JoinWaitingList(ctx, tc.email)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.email, result.Email)
			}
		})
	}
}

// Performance test
func BenchmarkWaitingListService_JoinWaitingList(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()
	expectedUser := buildWaitingListUserWithID("123e4567-e89b-12d3-a456-426614174000")

	// Setup expectations
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	mockWaitingListRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil).AnyTimes()
	mockLogger.EXPECT().Info("User successfully added to waiting list", gomock.Any()).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.JoinWaitingList(ctx, fmt.Sprintf("user%d@example.com", i))
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// Concurrent testing
func TestWaitingListService_JoinWaitingList_Concurrent(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx := context.Background()

	const numGoroutines = 10
	results := make(chan *waiting_list.WaitingListUser, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Mock expectations for concurrent calls
	for i := 0; i < numGoroutines; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		expectedUser := buildWaitingListUserWithID(fmt.Sprintf("123e4567-e89b-12d3-a456-426614%03d", i))

		mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)
		mockWaitingListRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)
		mockLogger.EXPECT().Info("User successfully added to waiting list", gomock.Any())
	}

	// Act
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			email := fmt.Sprintf("user%d@example.com", id)
			result, err := service.JoinWaitingList(ctx, email)
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
func TestWaitingListService_Timeout(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewWaitingListService(mockWaitingListRepo, mockLogger)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	email := "timeout@example.com"
	timeoutError := errors.New("context deadline exceeded")

	// Setup expectations
	mockWaitingListRepo.EXPECT().GetByEmail(ctx, email).Return(nil, timeoutError)
	mockWaitingListRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, timeoutError)
	mockLogger.EXPECT().Error("Failed to add user to waiting list", gomock.Any())

	// Act
	result, err := service.JoinWaitingList(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, shared.ErrWaitingListUserCreationFailed, err)
}

// Test with different email formats
func TestWaitingListService_EmailFormats(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		shouldSucceed bool
		description   string
	}{
		{
			name:          "Standard lowercase email",
			email:         "user@example.com",
			shouldSucceed: true,
			description:   "Standard email format should work",
		},
		{
			name:          "Standard uppercase email",
			email:         "USER@EXAMPLE.COM",
			shouldSucceed: true,
			description:   "Uppercase email should work",
		},
		{
			name:          "Mixed case email",
			email:         "UsEr@ExAmPlE.cOm",
			shouldSucceed: true,
			description:   "Mixed case email should work",
		},
		{
			name:          "Email with dots",
			email:         "user.name@example.com",
			shouldSucceed: true,
			description:   "Email with dots should work",
		},
		{
			name:          "Email with plus",
			email:         "user+tag@example.com",
			shouldSucceed: true,
			description:   "Email with plus should work",
		},
		{
			name:          "Empty email",
			email:         "",
			shouldSucceed: true,
			description:   "Empty email should be processed (no validation in service)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWaitingListRepo := mock_ports.NewMockWaitingListRepositoryInterface(ctrl)
			mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

			service := NewWaitingListService(mockWaitingListRepo, mockLogger)
			ctx := context.Background()

			if tc.shouldSucceed {
				expectedUser := &waiting_list.WaitingListUser{Email: tc.email}

				mockWaitingListRepo.EXPECT().GetByEmail(ctx, tc.email).Return(nil, gorm.ErrRecordNotFound)
				mockWaitingListRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)
				mockLogger.EXPECT().Info("User successfully added to waiting list", gomock.Any())
			}

			// Act
			result, err := service.JoinWaitingList(ctx, tc.email)

			// Assert
			if tc.shouldSucceed {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.email, result.Email)
			}
		})
	}
}
