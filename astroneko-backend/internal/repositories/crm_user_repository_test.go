package repositories

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"astroneko-backend/internal/core/domain/crm_user"
	"astroneko-backend/testings/mock_ports"
)

// Test data builders for consistent test data
func buildCRMUser() *crm_user.CRMUser {
	return &crm_user.CRMUser{
		Username: "testuser",
		Password: "hashedpassword",
	}
}

func buildCRMUserWithID(id string) *crm_user.CRMUser {
	return &crm_user.CRMUser{
		Username: "testuser",
		Password: "hashedpassword",
	}
}

func buildCRMUserWithIDAndUsername(id, username string) *crm_user.CRMUser {
	return &crm_user.CRMUser{
		Username: username,
		Password: "hashedpassword",
	}
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func TestCRMUserRepository_Create_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	testUser := buildCRMUser()
	expectedUser := buildCRMUserWithID("123e4567-e89b-12d3-a456-426614174000")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(value interface{}) error {
		if user, ok := value.(*crm_user.CRMUser); ok {
			user.ID = expectedUser.ID
		}
		return nil
	})

	// Act
	result, err := repository.Create(ctx, testUser)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, testUser.Username, result.Username)
	assert.Equal(t, testUser.Password, result.Password)
}

func TestCRMUserRepository_Create_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	testUser := buildCRMUser()
	dbError := errors.New("database constraint violation")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Create(gomock.Any()).Return(dbError)

	// Act
	result, err := repository.Create(ctx, testUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)
}

func TestCRMUserRepository_GetByUsername_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	username := "testuser"
	expectedUser := buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", username)

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Where("username = ?", username).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
		if user, ok := dest.(*crm_user.CRMUser); ok {
			*user = *expectedUser
		}
		return nil
	})

	// Act
	result, err := repository.GetByUsername(ctx, username)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Username, result.Username)
}

func TestCRMUserRepository_GetByUsername_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	username := "nonexistentuser"

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Where("username = ?", username).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)

	// Act
	result, err := repository.GetByUsername(ctx, username)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestCRMUserRepository_GetByID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	testID := "123e4567-e89b-12d3-a456-426614174000"
	expectedUser := buildCRMUserWithID(testID)

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Where("id = ?", gomock.Any()).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
		if user, ok := dest.(*crm_user.CRMUser); ok {
			*user = *expectedUser
		}
		return nil
	})

	// Act
	result, err := repository.GetByID(ctx, testID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Username, result.Username)
}

func TestCRMUserRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	testID := "123e4567-e89b-12d3-a456-426614174000"

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Where("id = ?", gomock.Any()).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)

	// Act
	result, err := repository.GetByID(ctx, testID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestCRMUserRepository_GetByID_InvalidUUID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	invalidID := "invalid-uuid"

	// Act
	result, err := repository.GetByID(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid UUID")
}

func TestCRMUserRepository_Update_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	testUser := buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", "updateduser")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Save(testUser).Return(nil)

	// Act
	result, err := repository.Update(ctx, testUser)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, testUser.ID, result.ID)
	assert.Equal(t, testUser.Username, result.Username)
}

func TestCRMUserRepository_Update_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	testUser := buildCRMUserWithID("123e4567-e89b-12d3-a456-426614174000")
	dbError := errors.New("database constraint violation")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Save(testUser).Return(dbError)

	// Act
	result, err := repository.Update(ctx, testUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)
}

func TestCRMUserRepository_Delete_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	testID := "123e4567-e89b-12d3-a456-426614174000"

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Where("id = ?", gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

	// Act
	err := repository.Delete(ctx, testID)

	// Assert
	assert.NoError(t, err)
}

func TestCRMUserRepository_Delete_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	testID := "123e4567-e89b-12d3-a456-426614174000"
	dbError := errors.New("foreign key constraint violation")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Where("id = ?", gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(dbError)

	// Act
	err := repository.Delete(ctx, testID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dbError, err)
}

func TestCRMUserRepository_Delete_InvalidUUID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()
	invalidID := "invalid-uuid"

	// Act
	err := repository.Delete(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}

// Table-driven tests for edge cases
func TestCRMUserRepository_GetByUsername_InvalidUsernames(t *testing.T) {
	testCases := []struct {
		name     string
		username string
		error    string
	}{
		{
			name:     "Valid username",
			username: "testuser",
			error:    "",
		},
		{
			name:     "Empty string",
			username: "",
			error:    "", // Repository doesn't validate username, just passes through
		},
		{
			name:     "Very long username",
			username: string(make([]byte, 500)),
			error:    "",
		},
		{
			name:     "Username with special characters",
			username: "user@#$%^&*()",
			error:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
			repository := NewCRMUserRepository(mockDB)
			ctx := context.Background()

			if tc.username == "nonexistentuser" {
				// Setup expectations for not found case
				mockDB.EXPECT().WithContext(ctx).Return(mockDB)
				mockDB.EXPECT().Where("username = ?", tc.username).Return(mockDB)
				mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)
			} else {
				// Setup expectations for success case
				expectedUser := buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", tc.username)
				mockDB.EXPECT().WithContext(ctx).Return(mockDB)
				mockDB.EXPECT().Where("username = ?", tc.username).Return(mockDB)
				mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
					if user, ok := dest.(*crm_user.CRMUser); ok {
						*user = *expectedUser
					}
					return nil
				})
			}

			// Act
			result, err := repository.GetByUsername(ctx, tc.username)

			// Assert
			if tc.username == "nonexistentuser" {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				if err == nil {
					require.NotNil(t, result)
					assert.Equal(t, tc.username, result.Username)
				}
			}
		})
	}
}

// Concurrent testing
func TestCRMUserRepository_Create_Concurrent(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewCRMUserRepository(mockDB)
	ctx := context.Background()

	const numGoroutines = 10
	errors := make(chan error, numGoroutines)

	// Mock expectations for concurrent calls
	for i := 0; i < numGoroutines; i++ {
		username := fmt.Sprintf("user%d", i)
		expectedUser := buildCRMUserWithIDAndUsername(fmt.Sprintf("123e4567-e89b-12d3-a456-426614%03d", i), username)

		mockDB.EXPECT().WithContext(ctx).Return(mockDB)
		mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(value interface{}) error {
			if user, ok := value.(*crm_user.CRMUser); ok {
				user.ID = expectedUser.ID
				user.Username = expectedUser.Username
			}
			return nil
		})
	}

	// Act
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			testUser := &crm_user.CRMUser{
				Username: fmt.Sprintf("user%d", id),
				Password: hashPassword("password123"),
			}
			_, err := repository.Create(ctx, testUser)
			errors <- err
		}(i)
	}

	// Assert
	for i := 0; i < numGoroutines; i++ {
		err := <-errors
		assert.NoError(t, err)
	}
}

// Test with different username formats
func TestCRMUserRepository_UsernameFormats(t *testing.T) {
	testCases := []struct {
		name        string
		username    string
		shouldExist bool
		description string
	}{
		{
			name:        "Standard lowercase username",
			username:    "testuser",
			shouldExist: true,
			description: "Standard username format should work",
		},
		{
			name:        "Standard uppercase username",
			username:    "TESTUSER",
			shouldExist: true,
			description: "Uppercase username should work",
		},
		{
			name:        "Mixed case username",
			username:    "UsErNaMe",
			shouldExist: true,
			description: "Mixed case username should work",
		},
		{
			name:        "Username with numbers",
			username:    "user123",
			shouldExist: true,
			description: "Username with numbers should work",
		},
		{
			name:        "Username with underscores",
			username:    "test_user",
			shouldExist: true,
			description: "Username with underscores should work",
		},
		{
			name:        "Empty username",
			username:    "",
			shouldExist: false,
			description: "Empty username should not find user",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
			repository := NewCRMUserRepository(mockDB)
			ctx := context.Background()

			if tc.shouldExist {
				// Setup expectations for existing user
				expectedUser := buildCRMUserWithIDAndUsername("123e4567-e89b-12d3-a456-426614174000", tc.username)
				mockDB.EXPECT().WithContext(ctx).Return(mockDB)
				mockDB.EXPECT().Where("username = ?", tc.username).Return(mockDB)
				mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
					if user, ok := dest.(*crm_user.CRMUser); ok {
						*user = *expectedUser
					}
					return nil
				})
			} else {
				// Setup expectations for not found
				mockDB.EXPECT().WithContext(ctx).Return(mockDB)
				mockDB.EXPECT().Where("username = ?", tc.username).Return(mockDB)
				mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)
			}

			// Act
			result, err := repository.GetByUsername(ctx, tc.username)

			// Assert
			if tc.shouldExist {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.username, result.Username)
			} else {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
			}
		})
	}
}
