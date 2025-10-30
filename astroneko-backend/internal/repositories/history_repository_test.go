package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"astroneko-backend/internal/core/domain/history"
	"astroneko-backend/testings/mock_ports"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Test data builders
func buildTestSession(userID uuid.UUID) *history.Session {
	return &history.Session{
		ID:          uuid.New(),
		UserID:      userID,
		HistoryName: "Test Session",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}
}

func buildTestMessage(sessionID uuid.UUID) *history.Message {
	return &history.Message{
		ID:         uuid.New(),
		SessionID:  sessionID,
		Message:    "Test message content",
		Role:       "user",
		UsedTokens: 100,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// GetSessionsByUserID Tests
func TestHistoryRepository_GetSessionsByUserID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	userID := uuid.New()

	expectedSessions := []history.Session{
		*buildTestSession(userID),
		*buildTestSession(userID),
	}

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Where("user_id = ?", userID).
		Return(mockDB)

	mockDB.EXPECT().
		Order("updated_at DESC").
		Return(mockDB)

	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) error {
			sessions := dest.(*[]history.Session)
			*sessions = expectedSessions
			return nil
		})

	// Act
	sessions, err := repo.GetSessionsByUserID(ctx, userID, "updated_at", "desc", "")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sessions, 2)
	assert.Equal(t, userID, sessions[0].UserID)
	assert.Equal(t, userID, sessions[1].UserID)
}

func TestHistoryRepository_GetSessionsByUserID_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	userID := uuid.New()
	dbError := errors.New("database connection error")

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Where("user_id = ?", userID).
		Return(mockDB)

	mockDB.EXPECT().
		Order("updated_at DESC").
		Return(mockDB)

	mockDB.EXPECT().
		Find(gomock.Any()).
		Return(dbError)

	// Act
	sessions, err := repo.GetSessionsByUserID(ctx, userID, "updated_at", "desc", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, sessions)
	assert.Contains(t, err.Error(), "failed to get sessions")
}

func TestHistoryRepository_GetSessionsByUserID_EmptyResult(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	userID := uuid.New()

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Where("user_id = ?", userID).
		Return(mockDB)

	mockDB.EXPECT().
		Order("updated_at DESC").
		Return(mockDB)

	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) error {
			sessions := dest.(*[]history.Session)
			*sessions = []history.Session{}
			return nil
		})

	// Act
	sessions, err := repo.GetSessionsByUserID(ctx, userID, "updated_at", "desc", "")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, sessions)
}

func TestHistoryRepository_GetSessionsByUserID_WithSearchQuery(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	userID := uuid.New()
	searchQuery := "test"

	expectedSessions := []history.Session{
		*buildTestSession(userID),
	}

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Where("user_id = ?", userID).
		Return(mockDB)

	mockDB.EXPECT().
		Where("history_name ILIKE ?", "%"+searchQuery+"%").
		Return(mockDB)

	mockDB.EXPECT().
		Order("updated_at DESC").
		Return(mockDB)

	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) error {
			sessions := dest.(*[]history.Session)
			*sessions = expectedSessions
			return nil
		})

	// Act
	sessions, err := repo.GetSessionsByUserID(ctx, userID, "updated_at", "desc", searchQuery)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sessions, 1)
	assert.Equal(t, userID, sessions[0].UserID)
}

// GetSessionByID Tests
func TestHistoryRepository_GetSessionByID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()

	expectedSession := buildTestSession(userID)
	expectedSession.ID = sessionID

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Where("id = ?", sessionID).
		Return(mockDB)

	mockDB.EXPECT().
		First(gomock.Any(), gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) error {
			session := dest.(*history.Session)
			*session = *expectedSession
			return nil
		})

	// Act
	session, err := repo.GetSessionByID(ctx, sessionID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, sessionID, session.ID)
	assert.Equal(t, userID, session.UserID)
}

func TestHistoryRepository_GetSessionByID_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()
	notFoundError := errors.New("record not found")

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Where("id = ?", sessionID).
		Return(mockDB)

	mockDB.EXPECT().
		First(gomock.Any()).
		Return(notFoundError)

	// Act
	session, err := repo.GetSessionByID(ctx, sessionID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "failed to get session")
}

// ValidateSessionOwnership Tests
func TestHistoryRepository_ValidateSessionOwnership_Success_IsOwner(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Model(gomock.Any()).
		Return(mockDB)

	mockDB.EXPECT().
		Where("id = ? AND user_id = ?", sessionID, userID).
		Return(mockDB)

	mockDB.EXPECT().
		Count(gomock.Any()).
		DoAndReturn(func(count *int64) error {
			*count = 1
			return nil
		})

	// Act
	isOwner, err := repo.ValidateSessionOwnership(ctx, sessionID, userID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, isOwner)
}

func TestHistoryRepository_ValidateSessionOwnership_Success_NotOwner(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Model(gomock.Any()).
		Return(mockDB)

	mockDB.EXPECT().
		Where("id = ? AND user_id = ?", sessionID, userID).
		Return(mockDB)

	mockDB.EXPECT().
		Count(gomock.Any()).
		DoAndReturn(func(count *int64) error {
			*count = 0
			return nil
		})

	// Act
	isOwner, err := repo.ValidateSessionOwnership(ctx, sessionID, userID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, isOwner)
}

func TestHistoryRepository_ValidateSessionOwnership_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()
	dbError := errors.New("database connection error")

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Model(gomock.Any()).
		Return(mockDB)

	mockDB.EXPECT().
		Where("id = ? AND user_id = ?", sessionID, userID).
		Return(mockDB)

	mockDB.EXPECT().
		Count(gomock.Any()).
		Return(dbError)

	// Act
	isOwner, err := repo.ValidateSessionOwnership(ctx, sessionID, userID)

	// Assert
	assert.Error(t, err)
	assert.False(t, isOwner)
	assert.Contains(t, err.Error(), "failed to validate session ownership")
}

// GetMessagesBySessionID Tests
func TestHistoryRepository_GetMessagesBySessionID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()

	expectedMessages := []history.Message{
		*buildTestMessage(sessionID),
		*buildTestMessage(sessionID),
		*buildTestMessage(sessionID),
	}

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Where("session_id = ?", sessionID).
		Return(mockDB)

	mockDB.EXPECT().
		Order("created_at ASC").
		Return(mockDB)

	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) error {
			messages := dest.(*[]history.Message)
			*messages = expectedMessages
			return nil
		})

	// Act
	messages, err := repo.GetMessagesBySessionID(ctx, sessionID, "asc")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, messages, 3)
	assert.Equal(t, sessionID, messages[0].SessionID)
	assert.Equal(t, sessionID, messages[1].SessionID)
	assert.Equal(t, sessionID, messages[2].SessionID)
}

func TestHistoryRepository_GetMessagesBySessionID_EmptyResult(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Where("session_id = ?", sessionID).
		Return(mockDB)

	mockDB.EXPECT().
		Order("created_at ASC").
		Return(mockDB)

	mockDB.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) error {
			messages := dest.(*[]history.Message)
			*messages = []history.Message{}
			return nil
		})

	// Act
	messages, err := repo.GetMessagesBySessionID(ctx, sessionID, "asc")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, messages)
}

func TestHistoryRepository_GetMessagesBySessionID_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()
	dbError := errors.New("database connection error")

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Where("session_id = ?", sessionID).
		Return(mockDB)

	mockDB.EXPECT().
		Order("created_at ASC").
		Return(mockDB)

	mockDB.EXPECT().
		Find(gomock.Any()).
		Return(dbError)

	// Act
	messages, err := repo.GetMessagesBySessionID(ctx, sessionID, "asc")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.Contains(t, err.Error(), "failed to get messages")
}

func TestHistoryRepository_DeleteSession_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Delete(gomock.Any(), sessionID).
		Return(nil)

	// Act
	err := repo.DeleteSession(ctx, sessionID)

	// Assert
	assert.NoError(t, err)
}

func TestHistoryRepository_DeleteSession_DatabaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockDatabaseInterface(ctrl)
	repo := NewHistoryRepository(mockDB)

	ctx := context.Background()
	sessionID := uuid.New()
	dbError := errors.New("database connection error")

	// Expect DB calls
	mockDB.EXPECT().
		WithContext(ctx).
		Return(mockDB)

	mockDB.EXPECT().
		Delete(gomock.Any(), sessionID).
		Return(dbError)

	// Act
	err := repo.DeleteSession(ctx, sessionID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete session")
}
