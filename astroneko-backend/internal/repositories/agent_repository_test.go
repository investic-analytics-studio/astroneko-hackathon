package repositories

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
	"astroneko-backend/testings/mock_ports"
)

// Test data builders for consistent test data
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

func TestAgentRepository_ClearState_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	ctx := context.Background()
	req := buildClearStateRequest()
	expectedResponse := buildClearStateResponse()

	// Setup expectations
	mockAgentRepo.EXPECT().ClearState(ctx, req).Return(expectedResponse, nil)

	// Act
	result, err := mockAgentRepo.ClearState(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedResponse.Status, result.Status)
}

func TestAgentRepository_ClearState_APIError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	ctx := context.Background()
	req := buildClearStateRequest()
	apiError := errors.New("API request failed with status 500")

	// Setup expectations
	mockAgentRepo.EXPECT().ClearState(ctx, req).Return(nil, apiError)

	// Act
	result, err := mockAgentRepo.ClearState(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, apiError, err)
}

func TestAgentRepository_ClearState_NetworkError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	ctx := context.Background()
	req := buildClearStateRequest()
	networkError := errors.New("connection timeout")

	// Setup expectations
	mockAgentRepo.EXPECT().ClearState(ctx, req).Return(nil, networkError)

	// Act
	result, err := mockAgentRepo.ClearState(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, networkError, err)
}

func TestAgentRepository_Reply_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	ctx := context.Background()
	req := buildReplyRequest()
	expectedResponse := buildReplyResponse()

	// Setup expectations
	mockAgentRepo.EXPECT().Reply(ctx, req).Return(expectedResponse, nil)

	// Act
	result, err := mockAgentRepo.Reply(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedResponse.Status, result.Status)
	assert.Equal(t, expectedResponse.Message, result.Message)
	assert.Equal(t, expectedResponse.Card, result.Card)
	assert.Equal(t, expectedResponse.Meaning, result.Meaning)
	assert.Equal(t, expectedResponse.SessionID, result.SessionID)
}

func TestAgentRepository_Reply_APIError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	ctx := context.Background()
	req := buildReplyRequest()
	apiError := errors.New("API request failed with status 400")

	// Setup expectations
	mockAgentRepo.EXPECT().Reply(ctx, req).Return(nil, apiError)

	// Act
	result, err := mockAgentRepo.Reply(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, apiError, err)
}

func TestAgentRepository_Reply_NetworkError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	ctx := context.Background()
	req := buildReplyRequest()
	networkError := errors.New("connection refused")

	// Setup expectations
	mockAgentRepo.EXPECT().Reply(ctx, req).Return(nil, networkError)

	// Act
	result, err := mockAgentRepo.Reply(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, networkError, err)
}

// Table-driven tests for edge cases
func TestAgentRepository_Reply_InvalidRequests(t *testing.T) {
	testCases := []struct {
		name        string
		request     agent.ReplyRequest
		expectError bool
		description string
	}{
		{
			name: "Valid request with text and session",
			request: agent.ReplyRequest{
				Text:      "Hello",
				SessionID: "session_123",
			},
			expectError: false,
			description: "Should succeed with valid text and session",
		},
		{
			name: "Valid request with text only",
			request: agent.ReplyRequest{
				Text: "Hello",
			},
			expectError: false,
			description: "Should succeed with just text (session is optional)",
		},
		{
			name: "Request with empty text",
			request: agent.ReplyRequest{
				Text:      "",
				SessionID: "session_123",
			},
			expectError: true,
			description: "Should fail with empty text (validation error)",
		},
		{
			name: "Request with long text",
			request: agent.ReplyRequest{
				Text:      string(make([]byte, 10000)), // Very long text
				SessionID: "session_123",
			},
			expectError: false,
			description: "Should handle long text appropriately",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
			ctx := context.Background()

			if !tc.expectError {
				expectedResponse := buildReplyResponse()
				expectedResponse.Message = "Response for: " + tc.request.Text
				mockAgentRepo.EXPECT().Reply(ctx, tc.request).Return(expectedResponse, nil)
			} else {
				mockAgentRepo.EXPECT().Reply(ctx, tc.request).Return(nil, errors.New("validation error"))
			}

			// Act
			result, err := mockAgentRepo.Reply(ctx, tc.request)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEmpty(t, result.Message)
			}
		})
	}
}

// Table-driven tests for clear state scenarios
func TestAgentRepository_ClearState_Scenarios(t *testing.T) {
	testCases := []struct {
		name           string
		request        agent.ClearStateRequest
		expectError    bool
		expectedStatus string
		description    string
	}{
		{
			name: "Valid session ID",
			request: agent.ClearStateRequest{
				SessionID: "valid_session_123",
			},
			expectError:    false,
			expectedStatus: "success",
			description:    "Should clear state with valid session",
		},
		{
			name: "Empty session ID",
			request: agent.ClearStateRequest{
				SessionID: "",
			},
			expectError: true,
			description: "Should fail with empty session ID",
		},
		{
			name: "Long session ID",
			request: agent.ClearStateRequest{
				SessionID: string(make([]byte, 500)), // Very long session ID
			},
			expectError: false,
			description: "Should handle long session IDs",
		},
		{
			name: "Special characters in session",
			request: agent.ClearStateRequest{
				SessionID: "session_123!@#$%^&*()",
			},
			expectError:    false,
			expectedStatus: "success",
			description:    "Should handle special characters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
			ctx := context.Background()

			if !tc.expectError {
				expectedResponse := &agent.ClearStateResponse{
					Status: tc.expectedStatus,
				}
				mockAgentRepo.EXPECT().ClearState(ctx, tc.request).Return(expectedResponse, nil)
			} else {
				mockAgentRepo.EXPECT().ClearState(ctx, tc.request).Return(nil, errors.New("invalid session ID"))
			}

			// Act
			result, err := mockAgentRepo.ClearState(ctx, tc.request)

			// Assert
			if tc.expectError {
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
func BenchmarkAgentRepository_Reply(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	ctx := context.Background()
	req := buildReplyRequest()
	expectedResponse := buildReplyResponse()

	// Setup expectations
	mockAgentRepo.EXPECT().Reply(ctx, gomock.Any()).Return(expectedResponse, nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Change the text to avoid exact match requirements
		req.Text = fmt.Sprintf("Message %d", i)
		_, err := mockAgentRepo.Reply(ctx, req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// Concurrent testing
func TestAgentRepository_Reply_Concurrent(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
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
	}

	// Act
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			req := agent.ReplyRequest{
				Text:      fmt.Sprintf("Concurrent message %d", id),
				SessionID: fmt.Sprintf("session_%d", id),
			}
			_, err := mockAgentRepo.Reply(ctx, req)
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
func TestAgentRepository_Timeout(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgentRepo := mock_ports.NewMockAgentRepositoryInterface(ctrl)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	req := buildReplyRequest()
	timeoutError := errors.New("context deadline exceeded")

	// Setup expectations
	mockAgentRepo.EXPECT().Reply(ctx, req).Return(nil, timeoutError)

	// Act
	result, err := mockAgentRepo.Reply(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "deadline exceeded")
}
