package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/testings/mock_ports"
)

// Test data builders for consistent test data
func buildTestUser() *user.User {
	now := time.Now()
	return &user.User{
		Email:               "test@example.com",
		IsActivatedReferral: false,
		FirebaseUID:         "firebase_123",
		LatestLoginAt:       &now,
		ProfileImageURL:     stringPtr("https://example.com/avatar.jpg"),
		DisplayName:         stringPtr("Test User"),
	}
}

func buildUpdatedTestUser(original *user.User) *user.User {
	now := time.Now()
	return &user.User{
		Email:               original.Email,
		IsActivatedReferral: true,
		FirebaseUID:         original.FirebaseUID,
		LatestLoginAt:       &now,
		ProfileImageURL:     stringPtr("https://example.com/new-avatar.jpg"),
		DisplayName:         stringPtr("Updated User"),
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

func TestUserRepository_Create_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	testUser := buildTestUser()

	expectedUser := buildTestUser()
	expectedUser.ID = uuid.New()

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(value interface{}) error {
		if user, ok := value.(*user.User); ok {
			user.ID = expectedUser.ID
		}
		return nil
	})

	// Act
	result, err := repository.Create(ctx, testUser)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, testUser.Email, result.Email)
	assert.Equal(t, testUser.FirebaseUID, result.FirebaseUID)
}

func TestUserRepository_Create_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	testUser := buildTestUser()

	dbError := errors.New("database connection failed")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Create(gomock.Any()).Return(dbError)

	// Act
	result, err := repository.Create(ctx, testUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), dbError.Error())
}

func TestUserRepository_GetByID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	testID := uuid.New().String()
	expectedUser := buildTestUser()
	expectedUser.ID = uuid.New()

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Where("id = ?", gomock.Any()).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
		if user, ok := dest.(*user.User); ok {
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
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.FirebaseUID, result.FirebaseUID)
}

func TestUserRepository_GetByID_InvalidUUID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	invalidID := "invalid-uuid"

	// Act
	result, err := repository.GetByID(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid UUID")
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	testID := uuid.New().String()

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Where("id = ?", gomock.Any()).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)

	// Act
	result, err := repository.GetByID(ctx, testID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "record not found")
}

func TestUserRepository_GetByFirebaseUID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	firebaseUID := "firebase_123"
	expectedUser := buildTestUser()
	expectedUser.FirebaseUID = firebaseUID

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Where("firebase_uid = ?", firebaseUID).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
		if user, ok := dest.(*user.User); ok {
			*user = *expectedUser
		}
		return nil
	})

	// Act
	result, err := repository.GetByFirebaseUID(ctx, firebaseUID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, firebaseUID, result.FirebaseUID)
	assert.Equal(t, expectedUser.Email, result.Email)
}

func TestUserRepository_GetByFirebaseUID_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	firebaseUID := "nonexistent_uid"

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Where("firebase_uid = ?", firebaseUID).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)

	// Act
	result, err := repository.GetByFirebaseUID(ctx, firebaseUID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "record not found")
}

func TestUserRepository_GetByEmail_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	email := "test@example.com"
	expectedUser := buildTestUser()
	expectedUser.Email = email

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Where("email = ?", email).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
		if user, ok := dest.(*user.User); ok {
			*user = *expectedUser
		}
		return nil
	})

	// Act
	result, err := repository.GetByEmail(ctx, email)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, email, result.Email)
	assert.Equal(t, expectedUser.FirebaseUID, result.FirebaseUID)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	email := "nonexistent@example.com"

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Where("email = ?", email).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)

	// Act
	result, err := repository.GetByEmail(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "record not found")
}

func TestUserRepository_Update_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	originalUser := buildTestUser()
	originalUser.ID = uuid.New()
	updatedUser := buildUpdatedTestUser(originalUser)

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Save(gomock.Any()).DoAndReturn(func(value interface{}) error {
		if user, ok := value.(*user.User); ok {
			user.IsActivatedReferral = updatedUser.IsActivatedReferral
			user.LatestLoginAt = updatedUser.LatestLoginAt
			user.ProfileImageURL = updatedUser.ProfileImageURL
			user.DisplayName = updatedUser.DisplayName
		}
		return nil
	})

	// Act
	result, err := repository.Update(ctx, updatedUser)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, updatedUser.IsActivatedReferral, result.IsActivatedReferral)
	assert.Equal(t, updatedUser.DisplayName, result.DisplayName)
}

func TestUserRepository_Update_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	testUser := buildTestUser()
	testUser.ID = uuid.New()

	dbError := errors.New("constraint violation")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Save(gomock.Any()).Return(dbError)

	// Act
	result, err := repository.Update(ctx, testUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), dbError.Error())
}

func TestUserRepository_Delete_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	testID := uuid.New().String()

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Where("id = ?", gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Delete(gomock.Any()).Return(nil)

	// Act
	err := repository.Delete(ctx, testID)

	// Assert
	assert.NoError(t, err)
}

func TestUserRepository_Delete_InvalidUUID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	invalidID := "invalid-uuid"

	// Act
	err := repository.Delete(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}

func TestUserRepository_Delete_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	testID := uuid.New().String()
	dbError := errors.New("foreign key constraint violation")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Where("id = ?", gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Delete(gomock.Any()).Return(dbError)

	// Act
	err := repository.Delete(ctx, testID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), dbError.Error())
}

func TestUserRepository_List_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	limit := 10
	offset := 0
	expectedUsers := []*user.User{
		buildTestUser(),
		buildTestUser(),
	}
	expectedCount := int64(2)

	// Set unique IDs for expected users
	for i, u := range expectedUsers {
		u.ID = uuid.New()
		u.Email = "user" + string(rune(i)) + "@example.com"
	}

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Model(gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Count(gomock.Any()).DoAndReturn(func(count *int64) error {
		*count = expectedCount
		return nil
	})
	mockDB.EXPECT().Limit(limit).Return(mockDB)
	mockDB.EXPECT().Offset(offset).Return(mockDB)
	mockDB.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
		if users, ok := dest.(*[]*user.User); ok {
			*users = expectedUsers
		}
		return nil
	})

	// Act
	users, count, err := repository.List(ctx, limit, offset)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.Len(t, users, len(expectedUsers))
}

func TestUserRepository_List_CountError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	limit := 10
	offset := 0
	dbError := errors.New("count query failed")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Model(gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Count(gomock.Any()).Return(dbError)

	// Act
	users, count, err := repository.List(ctx, limit, offset)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), dbError.Error())
}

func TestUserRepository_List_FindError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	limit := 10
	offset := 0
	expectedCount := int64(5)
	dbError := errors.New("find query failed")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Model(gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Count(gomock.Any()).DoAndReturn(func(count *int64) error {
		*count = expectedCount
		return nil
	})
	mockDB.EXPECT().Limit(limit).Return(mockDB)
	mockDB.EXPECT().Offset(offset).Return(mockDB)
	mockDB.EXPECT().Find(gomock.Any()).Return(dbError)

	// Act
	users, count, err := repository.List(ctx, limit, offset)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), dbError.Error())
}

func TestUserRepository_GetTotalUsers_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	expectedCount := int64(42)

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Model(gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Count(gomock.Any()).DoAndReturn(func(count *int64) error {
		*count = expectedCount
		return nil
	})

	// Act
	count, err := repository.GetTotalUsers(ctx)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestUserRepository_GetTotalUsers_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewUserRepository(mockDB)
	ctx := context.Background()
	dbError := errors.New("database connection failed")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Model(gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Count(gomock.Any()).Return(dbError)

	// Act
	count, err := repository.GetTotalUsers(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.Contains(t, err.Error(), dbError.Error())
}

// Table-driven tests for edge cases
func TestUserRepository_GetByID_InvalidUUIDs(t *testing.T) {
	testCases := []struct {
		name  string
		id    string
		error string
	}{
		{
			name:  "Empty string",
			id:    "",
			error: "invalid UUID",
		},
		{
			name:  "Invalid format",
			id:    "not-a-uuid",
			error: "invalid UUID",
		},
		{
			name:  "Partial UUID",
			id:    "12345678-1234-1234-1234",
			error: "invalid UUID",
		},
		{
			name:  "Too long",
			id:    "12345678-1234-1234-1234-1234567890123",
			error: "invalid UUID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
			repository := NewUserRepository(mockDB)
			ctx := context.Background()

			// Act
			result, err := repository.GetByID(ctx, tc.id)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tc.error)
		})
	}
}
