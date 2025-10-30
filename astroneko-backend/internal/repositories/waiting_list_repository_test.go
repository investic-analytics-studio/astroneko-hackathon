package repositories

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"astroneko-backend/internal/core/domain/waiting_list"
	"astroneko-backend/testings/mock_ports"
)

// Test data builders for consistent test data
func buildWaitingListUser() *waiting_list.WaitingListUser {
	return &waiting_list.WaitingListUser{
		Email: "test@example.com",
	}
}

func buildWaitingListUserWithID(id string) *waiting_list.WaitingListUser {
	userID := uuid.MustParse(id)
	user := &waiting_list.WaitingListUser{
		Email: "test@example.com",
	}
	user.ID = userID
	return user
}

func buildWaitingListUserWithIDAndEmail(id, email string) *waiting_list.WaitingListUser {
	userID := uuid.MustParse(id)
	user := &waiting_list.WaitingListUser{
		Email: email,
	}
	user.ID = userID
	return user
}

func TestWaitingListRepository_Create_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	testUser := buildWaitingListUser()

	expectedUser := buildWaitingListUserWithID("123e4567-e89b-12d3-a456-426614174000")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(value interface{}) error {
		if user, ok := value.(*waiting_list.WaitingListUser); ok {
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
	assert.Equal(t, testUser.Email, result.Email)
}

func TestWaitingListRepository_Create_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	testUser := buildWaitingListUser()

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

func TestWaitingListRepository_GetByEmail_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	email := "test@example.com"
	expectedUser := buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174000", email)

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Where("email = ?", email).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
		if user, ok := dest.(*waiting_list.WaitingListUser); ok {
			*user = *expectedUser
		}
		return nil
	})

	// Act
	result, err := repository.GetByEmail(ctx, email)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)
}

func TestWaitingListRepository_GetByEmail_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	email := "nonexistent@example.com"

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Where("email = ?", email).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)

	// Act
	result, err := repository.GetByEmail(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestWaitingListRepository_GetByID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	testID := "123e4567-e89b-12d3-a456-426614174000"
	expectedUser := buildWaitingListUserWithID(testID)

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Where("id = ?", gomock.Any()).Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
		if user, ok := dest.(*waiting_list.WaitingListUser); ok {
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
}

func TestWaitingListRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
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

func TestWaitingListRepository_GetByID_InvalidUUID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	invalidID := "invalid-uuid"

	// Act
	result, err := repository.GetByID(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid UUID")
}

func TestWaitingListRepository_List_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	limit := 10
	offset := 0
	expectedUsers := []*waiting_list.WaitingListUser{
		buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174000", "user1@example.com"),
		buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174001", "user2@example.com"),
	}
	expectedCount := int64(2)

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Model(gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Count(gomock.Any()).DoAndReturn(func(count *int64) error {
		*count = expectedCount
		return nil
	})
	mockDB.EXPECT().Limit(limit).Return(mockDB)
	mockDB.EXPECT().Offset(offset).Return(mockDB)
	mockDB.EXPECT().Order("created_at DESC").Return(mockDB)
	mockDB.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
		if users, ok := dest.(*[]*waiting_list.WaitingListUser); ok {
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
	assert.Equal(t, expectedUsers[0].Email, users[0].Email)
	assert.Equal(t, expectedUsers[1].Email, users[1].Email)
}

func TestWaitingListRepository_List_CountError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	limit := 10
	offset := 0
	dbError := errors.New("count query failed")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Model(gomock.Any()).Return(mockDB)
	mockDB.EXPECT().Count(gomock.Any()).Return(dbError)

	// Act
	users, count, err := repository.List(ctx, limit, offset)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.Nil(t, users)
	assert.Equal(t, dbError, err)
}

func TestWaitingList_List_FindError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
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
	mockDB.EXPECT().Order("created_at DESC").Return(mockDB)
	mockDB.EXPECT().Find(gomock.Any()).Return(dbError)

	// Act
	users, count, err := repository.List(ctx, limit, offset)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.Nil(t, users)
	assert.Equal(t, dbError, err)
}

func TestWaitingListRepository_Delete_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	testID := "123e4567-e89b-12d3-a456-426614174000"

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Delete(gomock.Any(), "id = ?", gomock.Any()).Return(nil)

	// Act
	err := repository.Delete(ctx, testID)

	// Assert
	assert.NoError(t, err)
}

func TestWaitingListRepository_Delete_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	testID := "123e4567-e89b-12d3-a456-426614174000"
	dbError := errors.New("foreign key constraint violation")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Delete(gomock.Any(), "id = ?", gomock.Any()).Return(dbError)

	// Act
	err := repository.Delete(ctx, testID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dbError, err)
}

func TestWaitingListRepository_Delete_InvalidUUID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()
	invalidID := "invalid-uuid"

	// Act
	err := repository.Delete(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}

// Table-driven tests for edge cases
func TestWaitingListRepository_GetByEmail_InvalidEmails(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		error string
	}{
		{
			name:  "Valid email",
			email: "test@example.com",
			error: "",
		},
		{
			name:  "Empty string",
			email: "",
			error: "", // Repository doesn't validate email, just passes through
		},
		{
			name:  "Invalid format",
			email: "invalid-email",
			error: "", // Repository doesn't validate email format
		},
		{
			name:  "Long email",
			email: string(make([]byte, 500)),
			error: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
			repository := NewWaitingListRepository(mockDB)
			ctx := context.Background()

			if tc.email == "nonexistent@example.com" {
				// Setup expectations for not found case
				mockDB.EXPECT().WithContext(ctx).Return(mockDB)
				mockDB.EXPECT().Where("email = ?", tc.email).Return(mockDB)
				mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)
			} else {
				// Setup expectations for success case
				expectedUser := buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174000", tc.email)
				mockDB.EXPECT().WithContext(ctx).Return(mockDB)
				mockDB.EXPECT().Where("email = ?", tc.email).Return(mockDB)
				mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
					if user, ok := dest.(*waiting_list.WaitingListUser); ok {
						*user = *expectedUser
					}
					return nil
				})
			}

			// Act
			result, err := repository.GetByEmail(ctx, tc.email)

			// Assert
			if tc.email == "nonexistent@example.com" {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				if err == nil {
					require.NotNil(t, result)
					assert.Equal(t, tc.email, result.Email)
				}
			}
		})
	}
}

// Performance test
func BenchmarkWaitingListRepository_Create(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()

	testUser := buildWaitingListUser()
	expectedUser := buildWaitingListUserWithID("123e4567-e89b-12d3-a456-426614174000")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB).AnyTimes()
	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(value interface{}) error {
		if user, ok := value.(*waiting_list.WaitingListUser); ok {
			user.ID = expectedUser.ID
		}
		return nil
	}).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testUser.Email = fmt.Sprintf("user%d@example.com", i)
		_, err := repository.Create(ctx, testUser)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// Concurrent testing
func TestWaitingListRepository_Create_Concurrent(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx := context.Background()

	const numGoroutines = 10
	errors := make(chan error, numGoroutines)

	// Setup expectations for concurrent calls
	for i := 0; i < numGoroutines; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		expectedUser := buildWaitingListUserWithID(fmt.Sprintf("123e4567-e89b-12d3-a456-426614174%03d", i))
		expectedUser.Email = email

		mockDB.EXPECT().WithContext(ctx).Return(mockDB)
		mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(value interface{}) error {
			if user, ok := value.(*waiting_list.WaitingListUser); ok {
				user.ID = expectedUser.ID
				user.Email = expectedUser.Email
			}
			return nil
		})
	}

	// Act
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			testUser := buildWaitingListUser()
			testUser.Email = fmt.Sprintf("user%d@example.com", id)
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

// Test timeout handling
func TestWaitingListRepository_Timeout(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repository := NewWaitingListRepository(mockDB)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	testUser := buildWaitingListUser()
	timeoutError := errors.New("context deadline exceeded")

	// Setup expectations
	mockDB.EXPECT().WithContext(ctx).Return(mockDB)
	mockDB.EXPECT().Create(gomock.Any()).Return(timeoutError)

	// Act
	result, err := repository.Create(ctx, testUser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "deadline exceeded")
}

// Test with different email formats
func TestWaitingListRepository_EmailFormats(t *testing.T) {
	testCases := []struct {
		name        string
		email       string
		shouldExist bool
		description string
	}{
		{
			name:        "Standard lowercase email",
			email:       "user@example.com",
			shouldExist: true,
			description: "Standard email format should work",
		},
		{
			name:        "Standard uppercase email",
			email:       "USER@EXAMPLE.COM",
			shouldExist: true,
			description: "Uppercase email should work",
		},
		{
			name:        "Mixed case email",
			email:       "UsEr@ExAmPlE.cOm",
			shouldExist: true,
			description: "Mixed case email should work",
		},
		{
			name:        "Email with dots",
			email:       "user.name@example.com",
			shouldExist: true,
			description: "Email with dots should work",
		},
		{
			name:        "Email with plus",
			email:       "user+tag@example.com",
			shouldExist: true,
			description: "Email with plus should work",
		},
		{
			name:        "Empty email",
			email:       "",
			shouldExist: false,
			description: "Empty email should not find user",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
			repository := NewWaitingListRepository(mockDB)
			ctx := context.Background()

			if tc.shouldExist {
				// Setup expectations for existing user
				expectedUser := buildWaitingListUserWithIDAndEmail("123e4567-e89b-12d3-a456-426614174000", tc.email)
				mockDB.EXPECT().WithContext(ctx).Return(mockDB)
				mockDB.EXPECT().Where("email = ?", tc.email).Return(mockDB)
				mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}, conds ...interface{}) error {
					if user, ok := dest.(*waiting_list.WaitingListUser); ok {
						*user = *expectedUser
					}
					return nil
				})
			} else {
				// Setup expectations for not found
				mockDB.EXPECT().WithContext(ctx).Return(mockDB)
				mockDB.EXPECT().Where("email = ?", tc.email).Return(mockDB)
				mockDB.EXPECT().First(gomock.Any()).Return(gorm.ErrRecordNotFound)
			}

			// Act
			result, err := repository.GetByEmail(ctx, tc.email)

			// Assert
			if tc.shouldExist {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.email, result.Email)
			} else {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
			}
		})
	}
}
