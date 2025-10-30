package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"astroneko-backend/internal/core/domain/history"
	"astroneko-backend/testings/mock_logger"
	"astroneko-backend/testings/mock_ports"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Test data builders
func buildTestHistorySession(userID uuid.UUID) *history.Session {
	return &history.Session{
		ID:          uuid.New(),
		UserID:      userID,
		HistoryName: "Test Conversation",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}
}

func buildTestHistoryMessage(sessionID uuid.UUID) *history.Message {
	return &history.Message{
		ID:         uuid.New(),
		SessionID:  sessionID,
		Message:    "What is the weather today?",
		Role:       "user",
		UsedTokens: 50,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// GetUserSessions Tests
func TestHistoryService_GetUserSessions_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	sessions := []history.Session{
		*buildTestHistorySession(userID),
		*buildTestHistorySession(userID),
	}

	// Expect repository call
	mockHistoryRepo.EXPECT().
		GetSessionsByUserID(ctx, userID, gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sessions, nil)

	// Act
	response, err := service.GetUserSessions(ctx, userID, "updated_at", "desc", "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, response.Total)
	assert.Len(t, response.Sessions, 2)
	assert.Equal(t, sessions[0].ID, response.Sessions[0].SessionID)
	assert.Equal(t, sessions[0].HistoryName, response.Sessions[0].HistoryName)
}

func TestHistoryService_GetUserSessions_EmptyResult(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Expect repository call
	mockHistoryRepo.EXPECT().
		GetSessionsByUserID(ctx, userID, gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]history.Session{}, nil)

	// Act
	response, err := service.GetUserSessions(ctx, userID, "updated_at", "desc", "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, response.Total)
	assert.Empty(t, response.Sessions)
}

func TestHistoryService_GetUserSessions_RepositoryError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	repoError := errors.New("database connection error")

	// Expect repository call and logger calls
	mockHistoryRepo.EXPECT().
		GetSessionsByUserID(ctx, userID, gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, repoError)

	mockLogger.EXPECT().
		Error(gomock.Eq("Failed to retrieve user sessions"), gomock.Any()).
		Times(1)

	// Act
	response, err := service.GetUserSessions(ctx, userID, "updated_at", "desc", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to retrieve sessions")
}

func TestHistoryService_GetUserSessions_WithSearchQuery(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	sessions := []history.Session{
		*buildTestHistorySession(userID),
	}

	// Expect repository call with search query
	mockHistoryRepo.EXPECT().
		GetSessionsByUserID(ctx, userID, gomock.Any(), gomock.Any(), "test").
		Return(sessions, nil)

	// Act
	response, err := service.GetUserSessions(ctx, userID, "updated_at", "desc", "test")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, response.Total)
	assert.Len(t, response.Sessions, 1)
}

// GetSessionMessages Tests
func TestHistoryService_GetSessionMessages_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()

	session := buildTestHistorySession(userID)
	session.ID = sessionID

	messages := []history.Message{
		*buildTestHistoryMessage(sessionID),
		*buildTestHistoryMessage(sessionID),
	}

	// Expect repository calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(true, nil)

	mockHistoryRepo.EXPECT().
		GetSessionByID(ctx, sessionID).
		Return(session, nil)

	mockHistoryRepo.EXPECT().
		GetMessagesBySessionID(ctx, sessionID, gomock.Any()).
		Return(messages, nil)

	// Act
	response, err := service.GetSessionMessages(ctx, userID, sessionID, "asc")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, sessionID, response.SessionID)
	assert.Equal(t, session.HistoryName, response.HistoryName)
	assert.Equal(t, 2, response.Total)
	assert.Len(t, response.Messages, 2)
	assert.Equal(t, messages[0].Message, response.Messages[0].Message)
	assert.Equal(t, messages[0].Role, response.Messages[0].Role)
	assert.Equal(t, messages[0].UsedTokens, response.Messages[0].UsedTokens)
}

func TestHistoryService_GetSessionMessages_UnauthorizedAccess(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()

	// Expect repository call and logger calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(false, nil)

	mockLogger.EXPECT().
		Warn(gomock.Eq("Unauthorized session access attempt"), gomock.Any()).
		Times(1)

	// Act
	response, err := service.GetSessionMessages(ctx, userID, sessionID, "asc")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "access denied")
}

func TestHistoryService_GetSessionMessages_ValidationError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()
	validationError := errors.New("validation query failed")

	// Expect repository call and logger calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(false, validationError)

	mockLogger.EXPECT().
		Error(gomock.Eq("Failed to validate session ownership"), gomock.Any()).
		Times(1)

	// Act
	response, err := service.GetSessionMessages(ctx, userID, sessionID, "asc")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to validate session ownership")
}

func TestHistoryService_GetSessionMessages_SessionNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()
	sessionError := errors.New("session not found")

	// Expect repository calls and logger calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(true, nil)

	mockHistoryRepo.EXPECT().
		GetSessionByID(ctx, sessionID).
		Return(nil, sessionError)

	mockLogger.EXPECT().
		Error(gomock.Eq("Failed to retrieve session"), gomock.Any()).
		Times(1)

	// Act
	response, err := service.GetSessionMessages(ctx, userID, sessionID, "asc")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to retrieve session")
}

func TestHistoryService_GetSessionMessages_MessagesRetrievalError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()

	session := buildTestHistorySession(userID)
	session.ID = sessionID

	messagesError := errors.New("messages retrieval failed")

	// Expect repository calls and logger calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(true, nil)

	mockHistoryRepo.EXPECT().
		GetSessionByID(ctx, sessionID).
		Return(session, nil)

	mockHistoryRepo.EXPECT().
		GetMessagesBySessionID(ctx, sessionID, gomock.Any()).
		Return(nil, messagesError)

	mockLogger.EXPECT().
		Error(gomock.Eq("Failed to retrieve session messages"), gomock.Any()).
		Times(1)

	// Act
	response, err := service.GetSessionMessages(ctx, userID, sessionID, "asc")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to retrieve messages")
}

func TestHistoryService_GetSessionMessages_EmptyMessages(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()

	session := buildTestHistorySession(userID)
	session.ID = sessionID

	// Expect repository calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(true, nil)

	mockHistoryRepo.EXPECT().
		GetSessionByID(ctx, sessionID).
		Return(session, nil)

	mockHistoryRepo.EXPECT().
		GetMessagesBySessionID(ctx, sessionID, gomock.Any()).
		Return([]history.Message{}, nil)

	// Act
	response, err := service.GetSessionMessages(ctx, userID, sessionID, "asc")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, sessionID, response.SessionID)
	assert.Equal(t, 0, response.Total)
	assert.Empty(t, response.Messages)
}

func TestHistoryService_DeleteSession_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()

	// Expect repository calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(true, nil)

	mockHistoryRepo.EXPECT().
		DeleteSession(ctx, sessionID).
		Return(nil)

	mockLogger.EXPECT().
		Info(gomock.Eq("Session deleted successfully"), gomock.Any()).
		Times(1)

	// Act
	err := service.DeleteSession(ctx, userID, sessionID)

	// Assert
	assert.NoError(t, err)
}

func TestHistoryService_DeleteSession_UnauthorizedAccess(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()

	// Expect repository call and logger calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(false, nil)

	mockLogger.EXPECT().
		Warn(gomock.Eq("Unauthorized session deletion attempt"), gomock.Any()).
		Times(1)

	// Act
	err := service.DeleteSession(ctx, userID, sessionID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestHistoryService_DeleteSession_ValidationError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()
	validationError := errors.New("validation query failed")

	// Expect repository call and logger calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(false, validationError)

	mockLogger.EXPECT().
		Error(gomock.Eq("Failed to validate session ownership for deletion"), gomock.Any()).
		Times(1)

	// Act
	err := service.DeleteSession(ctx, userID, sessionID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate session ownership")
}

func TestHistoryService_DeleteSession_DeleteError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHistoryRepo := mock_ports.NewHistoryRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewHistoryService(mockHistoryRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()
	deleteError := errors.New("delete operation failed")

	// Expect repository calls and logger calls
	mockHistoryRepo.EXPECT().
		ValidateSessionOwnership(ctx, sessionID, userID).
		Return(true, nil)

	mockHistoryRepo.EXPECT().
		DeleteSession(ctx, sessionID).
		Return(deleteError)

	mockLogger.EXPECT().
		Error(gomock.Eq("Failed to delete session"), gomock.Any()).
		Times(1)

	// Act
	err := service.DeleteSession(ctx, userID, sessionID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete session")
}
