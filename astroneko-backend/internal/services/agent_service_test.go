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

	"astroneko-backend/internal/core/domain/agent"
	"astroneko-backend/testings/mock_logger"
	"astroneko-backend/testings/mock_ports"
)

// Test data builders
func buildClearStateRequest() agent.ClearStateRequest {
	return agent.ClearStateRequest{
		SessionID: "session_123",
	}
}

func buildClearStateResponse() *agent.ClearStateResponse {
	return &agent.ClearStateResponse{
		Status: "success",
	}
}

func buildReplyRequest() agent.ReplyRequest {
	return agent.ReplyRequest{
		Text:      "Hello, agent!",
		SessionID: "session_123",
	}
}

func buildReplyResponse() *agent.ReplyResponse {
	return &agent.ReplyResponse{
		Status:    "success",
		Message:   "Hello! I'm the cat fortune agent. Nice to meet you!",
		Card:      "🐱",
		Meaning:   "The cat brings good fortune and wisdom",
		SessionID: "session_123",
	}
}

// Helper functions
func agentServiceStringPtr(s string) *string { return &s }

func TestAgentService_ClearState_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx := context.Background()
	userID := "user_123"
	req := buildClearStateRequest()
	expectedResponse := buildClearStateResponse()

	// Mock repository call
	mockAgentRepo.EXPECT().ClearState(ctx, req).Return(expectedResponse, nil)

	// Mock logger calls
	mockLogger.EXPECT().Info("Clearing agent state for user", gomock.Any())
	mockLogger.EXPECT().Info("Agent state cleared successfully", gomock.Any())

	// Act
	result, err := service.ClearState(ctx, userID, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedResponse.Status, result.Status)
}

func TestAgentService_ClearState_RepositoryError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx := context.Background()
	userID := "user_123"
	req := buildClearStateRequest()
	repoError := errors.New("API request failed")

	// Mock repository call
	mockAgentRepo.EXPECT().ClearState(ctx, req).Return(nil, repoError)

	// Mock logger calls
	mockLogger.EXPECT().Info("Clearing agent state for user", gomock.Any())
	mockLogger.EXPECT().Error("Failed to clear agent state", gomock.Any())

	// Act
	result, err := service.ClearState(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)
}

func TestAgentService_ClearState_ContextTimeout(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	userID := "user_123"
	req := buildClearStateRequest()
	timeoutError := errors.New("context deadline exceeded")

	// Mock repository call
	mockAgentRepo.EXPECT().ClearState(ctx, req).Return(nil, timeoutError)

	// Mock logger calls
	mockLogger.EXPECT().Info("Clearing agent state for user", gomock.Any())
	mockLogger.EXPECT().Error("Failed to clear agent state", gomock.Any())

	// Act
	result, err := service.ClearState(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "deadline exceeded")
}

func TestAgentService_Reply_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx := context.Background()
	userID := "user_123"
	req := buildReplyRequest()
	expectedResponse := buildReplyResponse()

	// Mock repository call
	mockAgentRepo.EXPECT().Reply(ctx, req).Return(expectedResponse, nil)

	// Mock logger calls
	mockLogger.EXPECT().Info("Sending message to agent", gomock.Any())
	mockLogger.EXPECT().Info("Agent reply received successfully", gomock.Any())

	// Act
	result, err := service.Reply(ctx, userID, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedResponse.Status, result.Status)
	assert.Equal(t, expectedResponse.Message, result.Message)
	assert.Equal(t, expectedResponse.Card, result.Card)
	assert.Equal(t, expectedResponse.Meaning, result.Meaning)
	assert.Equal(t, expectedResponse.SessionID, result.SessionID)
}

func TestAgentService_Reply_RepositoryError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx := context.Background()
	userID := "user_123"
	req := buildReplyRequest()
	repoError := errors.New("API request failed with status 500")

	// Mock repository call
	mockAgentRepo.EXPECT().Reply(ctx, req).Return(nil, repoError)

	// Mock logger calls
	mockLogger.EXPECT().Info("Sending message to agent", gomock.Any())
	mockLogger.EXPECT().Error("Failed to get agent reply", gomock.Any())

	// Act
	result, err := service.Reply(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)
}

func TestAgentService_Reply_NetworkError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx := context.Background()
	userID := "user_123"
	req := buildReplyRequest()
	networkError := errors.New("connection timeout")

	// Mock repository call
	mockAgentRepo.EXPECT().Reply(ctx, req).Return(nil, networkError)

	// Mock logger calls
	mockLogger.EXPECT().Info("Sending message to agent", gomock.Any())
	mockLogger.EXPECT().Error("Failed to get agent reply", gomock.Any())

	// Act
	result, err := service.Reply(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, networkError, err)
}

func TestAgentService_Reply_EmptyText(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx := context.Background()
	userID := "user_123"
	req := agent.ReplyRequest{
		Text:      "", // Empty text
		SessionID: "session_123",
	}
	validationError := errors.New("text cannot be empty")

	// Mock repository call
	mockAgentRepo.EXPECT().Reply(ctx, req).Return(nil, validationError)

	// Mock logger calls
	mockLogger.EXPECT().Info("Sending message to agent", gomock.Any())
	mockLogger.EXPECT().Error("Failed to get agent reply", gomock.Any())

	// Act
	result, err := service.Reply(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, validationError, err)
}

func TestAgentService_Reply_LongText(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx := context.Background()
	userID := "user_123"

	// Create a very long text message
	longText := string(make([]byte, 5000))
	for i := range longText {
		longText = longText[:i] + "a" + longText[i+1:]
	}

	req := agent.ReplyRequest{
		Text:      longText,
		SessionID: "session_123",
	}
	expectedResponse := &agent.ReplyResponse{
		Status:    "success",
		Message:   "Response received for long message",
		SessionID: "session_123",
	}

	// Mock repository call
	mockAgentRepo.EXPECT().Reply(ctx, req).Return(expectedResponse, nil)

	// Mock logger calls
	mockLogger.EXPECT().Info("Sending message to agent", gomock.Any())
	mockLogger.EXPECT().Info("Agent reply received successfully", gomock.Any())

	// Act
	result, err := service.Reply(ctx, userID, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedResponse.Status, result.Status)
	assert.Equal(t, expectedResponse.Message, result.Message)
}

// Table-driven tests for business logic scenarios
func TestAgentService_Reply_BusinessLogic(t *testing.T) {
	testCases := []struct {
		name           string
		userID         string
		request        agent.ReplyRequest
		repoError      error
		expectedError  bool
		expectedStatus string
		description    string
	}{
		{
			name:   "Valid user and request",
			userID: "user_123",
			request: agent.ReplyRequest{
				Text:      "Hello",
				SessionID: "session_123",
			},
			repoError:      nil,
			expectedError:  false,
			expectedStatus: "success",
			description:    "Should succeed with valid user and request",
		},
		{
			name:   "User with special characters",
			userID: "user_special-123!@#",
			request: agent.ReplyRequest{
				Text:      "Hello",
				SessionID: "session_123",
			},
			repoError:      nil,
			expectedError:  false,
			expectedStatus: "success",
			description:    "Should handle special characters in user ID",
		},
		{
			name:   "Request without session ID",
			userID: "user_123",
			request: agent.ReplyRequest{
				Text: "Hello",
			},
			repoError:      nil,
			expectedError:  false,
			expectedStatus: "success",
			description:    "Should handle requests without session ID",
		},
		{
			name:   "Repository returns error",
			userID: "user_123",
			request: agent.ReplyRequest{
				Text:      "Hello",
				SessionID: "session_123",
			},
			repoError:     errors.New("API rate limit exceeded"),
			expectedError: true,
			description:   "Should propagate repository errors",
		},
		{
			name:   "Empty user ID",
			userID: "",
			request: agent.ReplyRequest{
				Text:      "Hello",
				SessionID: "session_123",
			},
			repoError:      nil,
			expectedError:  false,
			expectedStatus: "success",
			description:    "Should handle empty user ID (no validation in service)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
			mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

			service := NewAgentService(mockAgentRepo, mockLogger)
			ctx := context.Background()

			// Setup repository expectation
			if tc.repoError != nil {
				mockAgentRepo.EXPECT().Reply(ctx, tc.request).Return(nil, tc.repoError)
			} else {
				expectedResponse := &agent.ReplyResponse{
					Status:    tc.expectedStatus,
					Message:   "Response for: " + tc.request.Text,
					SessionID: tc.request.SessionID,
				}
				mockAgentRepo.EXPECT().Reply(ctx, tc.request).Return(expectedResponse, nil)
			}

			// Setup logger expectations
			mockLogger.EXPECT().Info("Sending message to agent", gomock.Any())

			if tc.repoError != nil {
				mockLogger.EXPECT().Error("Failed to get agent reply", gomock.Any())
			} else {
				mockLogger.EXPECT().Info("Agent reply received successfully", gomock.Any())
			}

			// Act
			result, err := service.Reply(ctx, tc.userID, tc.request)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.expectedStatus, result.Status)
			}
		})
	}
}

// Table-driven tests for clear state business logic
func TestAgentService_ClearState_BusinessLogic(t *testing.T) {
	testCases := []struct {
		name           string
		userID         string
		request        agent.ClearStateRequest
		repoError      error
		expectedError  bool
		expectedStatus string
		description    string
	}{
		{
			name:   "Valid user and session",
			userID: "user_123",
			request: agent.ClearStateRequest{
				SessionID: "session_123",
			},
			repoError:      nil,
			expectedError:  false,
			expectedStatus: "success",
			description:    "Should clear state with valid inputs",
		},
		{
			name:   "Empty session ID",
			userID: "user_123",
			request: agent.ClearStateRequest{
				SessionID: "",
			},
			repoError:     errors.New("session ID cannot be empty"),
			expectedError: true,
			description:   "Should fail with empty session ID",
		},
		{
			name:   "Repository error",
			userID: "user_123",
			request: agent.ClearStateRequest{
				SessionID: "session_123",
			},
			repoError:     errors.New("database connection failed"),
			expectedError: true,
			description:   "Should propagate repository errors",
		},
		{
			name:   "Long session ID",
			userID: "user_123",
			request: agent.ClearStateRequest{
				SessionID: string(make([]byte, 1000)),
			},
			repoError:      nil,
			expectedError:  false,
			expectedStatus: "success",
			description:    "Should handle long session IDs",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
			mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

			service := NewAgentService(mockAgentRepo, mockLogger)
			ctx := context.Background()

			// Setup repository expectation
			if tc.repoError != nil {
				mockAgentRepo.EXPECT().ClearState(ctx, tc.request).Return(nil, tc.repoError)
			} else {
				expectedResponse := &agent.ClearStateResponse{
					Status: tc.expectedStatus,
				}
				mockAgentRepo.EXPECT().ClearState(ctx, tc.request).Return(expectedResponse, nil)
			}

			// Setup logger expectations
			mockLogger.EXPECT().Info("Clearing agent state for user", gomock.Any())

			if tc.repoError != nil {
				mockLogger.EXPECT().Error("Failed to clear agent state", gomock.Any())
			} else {
				mockLogger.EXPECT().Info("Agent state cleared successfully", gomock.Any())
			}

			// Act
			result, err := service.ClearState(ctx, tc.userID, tc.request)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.expectedStatus, result.Status)
			}
		})
	}
}

// Performance test
func BenchmarkAgentService_Reply(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx := context.Background()
	userID := "user_123"
	req := buildReplyRequest()
	expectedResponse := buildReplyResponse()

	// Setup expectations
	mockAgentRepo.EXPECT().Reply(ctx, gomock.Any()).Return(expectedResponse, nil).AnyTimes()
	mockLogger.EXPECT().Info("Sending message to agent", gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info("Agent reply received successfully", gomock.Any()).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Text = fmt.Sprintf("Message %d", i)
		_, err := service.Reply(ctx, userID, req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// Concurrent testing
func TestAgentService_Concurrent(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	service := NewAgentService(mockAgentRepo, mockLogger)
	ctx := context.Background()

	const numGoroutines = 10
	errors := make(chan error, numGoroutines)

	// Mock expectations for concurrent calls
	for i := 0; i < numGoroutines; i++ {
		req := agent.ReplyRequest{
			Text:      fmt.Sprintf("Concurrent message %d", i),
			SessionID: fmt.Sprintf("session_%d", i),
		}
		expectedResponse := &agent.ReplyResponse{
			Status:    "success",
			Message:   fmt.Sprintf("Response for message %d", i),
			SessionID: fmt.Sprintf("session_%d", i),
		}

		mockAgentRepo.EXPECT().Reply(ctx, req).Return(expectedResponse, nil)
		mockLogger.EXPECT().Info("Sending message to agent", gomock.Any())
		mockLogger.EXPECT().Info("Agent reply received successfully", gomock.Any())
	}

	// Act
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			userID := fmt.Sprintf("user_%d", id)
			req := agent.ReplyRequest{
				Text:      fmt.Sprintf("Concurrent message %d", id),
				SessionID: fmt.Sprintf("session_%d", id),
			}
			_, err := service.Reply(ctx, userID, req)
			errors <- err
		}(i)
	}

	// Assert
	for i := 0; i < numGoroutines; i++ {
		err := <-errors
		assert.NoError(t, err)
	}
}

// Test service with different user behaviors
func TestAgentService_UserBehavior(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		request     agent.ReplyRequest
		description string
	}{
		{
			name:   "Active user with session",
			userID: "active_user_123",
			request: agent.ReplyRequest{
				Text:      "Hello from active user",
				SessionID: "persistent_session_456",
			},
			description: "Active user with persistent session",
		},
		{
			name:   "Anonymous user",
			userID: "anonymous",
			request: agent.ReplyRequest{
				Text: "Hello from anonymous user",
			},
			description: "Anonymous user without session",
		},
		{
			name:   "Bot user",
			userID: "bot_user_789",
			request: agent.ReplyRequest{
				Text:      "Bot command",
				SessionID: "bot_session_789",
			},
			description: "Bot user with specific session",
		},
		{
			name:   "Premium user",
			userID: "premium_user_999",
			request: agent.ReplyRequest{
				Text:      "Premium user request",
				SessionID: "premium_session_999",
			},
			description: "Premium user with enhanced session",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
			mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

			service := NewAgentService(mockAgentRepo, mockLogger)
			ctx := context.Background()

			expectedResponse := &agent.ReplyResponse{
				Status:    "success",
				Message:   "Response for " + tc.request.Text,
				SessionID: tc.request.SessionID,
			}

			// Mock repository call
			mockAgentRepo.EXPECT().Reply(ctx, tc.request).Return(expectedResponse, nil)

			// Mock logger calls
			mockLogger.EXPECT().Info("Sending message to agent", gomock.Any())
			mockLogger.EXPECT().Info("Agent reply received successfully", gomock.Any())

			// Act
			result, err := service.Reply(ctx, tc.userID, tc.request)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, "success", result.Status)
			assert.Contains(t, result.Message, tc.request.Text)
		})
	}
}
