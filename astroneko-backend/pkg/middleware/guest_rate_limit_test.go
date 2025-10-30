package middleware

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"astroneko-backend/internal/core/domain/guest_usage"
	"astroneko-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger is a simple mock implementation of logger.Logger
type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...logger.Field)  {}
func (m *MockLogger) Warn(msg string, fields ...logger.Field)  {}
func (m *MockLogger) Error(msg string, fields ...logger.Field) {}

// MockGuestUsageRepository is a mock implementation of guest_usage.Repository
type MockGuestUsageRepository struct {
	mock.Mock
}

func (m *MockGuestUsageRepository) Create(ctx context.Context, usage *guest_usage.GuestAPIUsage) error {
	args := m.Called(ctx, usage)
	return args.Error(0)
}

func (m *MockGuestUsageRepository) GetByCompositeKey(ctx context.Context, compositeKey, endpoint, windowResetStr string) (*guest_usage.GuestAPIUsage, error) {
	args := m.Called(ctx, compositeKey, endpoint, windowResetStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guest_usage.GuestAPIUsage), args.Error(1)
}

func (m *MockGuestUsageRepository) IncrementUsage(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGuestUsageRepository) GetByIPAddress(ctx context.Context, ipAddress, since string) ([]*guest_usage.GuestAPIUsage, error) {
	args := m.Called(ctx, ipAddress, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*guest_usage.GuestAPIUsage), args.Error(1)
}

func (m *MockGuestUsageRepository) BlockGuest(ctx context.Context, compositeKey, reason string) error {
	args := m.Called(ctx, compositeKey, reason)
	return args.Error(0)
}

func (m *MockGuestUsageRepository) GetAll(ctx context.Context) ([]*guest_usage.GuestAPIUsage, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*guest_usage.GuestAPIUsage), args.Error(1)
}

func (m *MockGuestUsageRepository) DeleteByID(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGuestUsageRepository) DeleteOldRecords(ctx context.Context, olderThan string) error {
	args := m.Called(ctx, olderThan)
	return args.Error(0)
}

func (m *MockGuestUsageRepository) ResetExpiredWindows(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestGuestOrAuthRateLimit_LoggedInWithReferral tests unlimited access for users with referral
func TestGuestOrAuthRateLimit_LoggedInWithReferral(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGuestUsageRepository)
	log := &MockLogger{}

	middleware := NewGuestRateLimitMiddleware(mockRepo, log)
	handler := middleware.GuestOrAuthRateLimit("/api/v1/agent/reply", 3)

	app.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("user_type", "logged_in_with_referral")
		return c.Next()
	}, handler, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	// Test: User with activated referral should have unlimited access
	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify no repository calls were made (no rate limiting)
	mockRepo.AssertNotCalled(t, "GetByCompositeKey")
	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertNotCalled(t, "IncrementUsage")
}

// TestGuestOrAuthRateLimit_LoggedInNoReferral tests daily limit for logged-in users without referral
func TestGuestOrAuthRateLimit_LoggedInNoReferral_FirstRequest(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGuestUsageRepository)
	log := &MockLogger{}

	middleware := NewGuestRateLimitMiddleware(mockRepo, log)
	handler := middleware.GuestOrAuthRateLimit("/api/v1/agent/reply", 3)

	// Mock repository to return nil (first request)
	mockRepo.On("GetByCompositeKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(usage *guest_usage.GuestAPIUsage) bool {
		return usage.UsageCount == 1 && usage.DailyLimit == 3
	})).Return(nil)

	app.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("user_type", "logged_in_no_referral")
		c.Locals("user", "mock_user")
		c.Locals("firebase_uid", "test_uid_123")
		return c.Next()
	}, handler, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify Create was called once
	mockRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
}

// TestGuestOrAuthRateLimit_LoggedInNoReferral_DailyLimitExceeded tests limit enforcement
func TestGuestOrAuthRateLimit_LoggedInNoReferral_DailyLimitExceeded(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGuestUsageRepository)
	log := &MockLogger{}

	middleware := NewGuestRateLimitMiddleware(mockRepo, log)
	handler := middleware.GuestOrAuthRateLimit("/api/v1/agent/reply", 3)

	// Mock repository to return usage at limit
	existingUsage := &guest_usage.GuestAPIUsage{
		ID:            "1",
		CompositeKey:  "user_test_uid_123",
		Endpoint:      "/api/v1/agent/reply",
		UsageCount:    3,
		DailyLimit:    3,
		WindowResetAt: time.Now().Add(12 * time.Hour),
		IsBlocked:     false,
	}
	mockRepo.On("GetByCompositeKey", mock.Anything, "user_test_uid_123", mock.Anything, mock.Anything).Return(existingUsage, nil)

	app.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("user_type", "logged_in_no_referral")
		c.Locals("user", "mock_user")
		c.Locals("firebase_uid", "test_uid_123")
		return c.Next()
	}, handler, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode) // Too Many Requests

	// Parse response body
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "Daily limit exceeded", body["error"])
	assert.Contains(t, body["message"], "activate a referral code")
}

// TestGuestOrAuthRateLimit_Guest_FirstRequest tests guest lifetime limit
func TestGuestOrAuthRateLimit_Guest_FirstRequest(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGuestUsageRepository)
	log := &MockLogger{}

	middleware := NewGuestRateLimitMiddleware(mockRepo, log)
	handler := middleware.GuestOrAuthRateLimit("/api/v1/agent/reply", 3)

	// Mock repository to return nil (first request)
	mockRepo.On("GetByCompositeKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(usage *guest_usage.GuestAPIUsage) bool {
		// Verify lifetime window is set to year 9999
		return usage.UsageCount == 1 &&
			usage.DailyLimit == 3 &&
			usage.WindowResetAt.Year() == 9999
	})).Return(nil)

	app.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("user_type", "guest")
		return c.Next()
	}, handler, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("User-Agent", "TestBrowser/1.0")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify Create was called with lifetime window
	mockRepo.AssertCalled(t, "Create", mock.Anything, mock.MatchedBy(func(usage *guest_usage.GuestAPIUsage) bool {
		return usage.WindowResetAt.Year() == 9999
	}))
}

// TestGuestOrAuthRateLimit_Guest_LifetimeLimitExceeded tests guest cannot exceed 3 lifetime
func TestGuestOrAuthRateLimit_Guest_LifetimeLimitExceeded(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGuestUsageRepository)
	log := &MockLogger{}

	middleware := NewGuestRateLimitMiddleware(mockRepo, log)
	handler := middleware.GuestOrAuthRateLimit("/api/v1/agent/reply", 3)

	// Mock repository to return usage at lifetime limit
	lifetimeWindow := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	existingUsage := &guest_usage.GuestAPIUsage{
		ID:            "1",
		CompositeKey:  "test_composite_key",
		Endpoint:      "/api/v1/agent/reply",
		UsageCount:    3,
		DailyLimit:    3,
		WindowResetAt: lifetimeWindow, // Never resets
		IsBlocked:     false,
	}
	mockRepo.On("GetByCompositeKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(existingUsage, nil)

	app.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("user_type", "guest")
		return c.Next()
	}, handler, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("User-Agent", "TestBrowser/1.0")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode) // Too Many Requests

	// Parse response body
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "Free trial limit exceeded", body["error"])
	assert.Contains(t, body["message"], "sign in")
	assert.NotContains(t, body, "reset_in") // No reset for guests
}

// TestGuestOrAuthRateLimit_Guest_SecondRequest tests guest can make multiple requests until limit
func TestGuestOrAuthRateLimit_Guest_SecondRequest(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGuestUsageRepository)
	log := &MockLogger{}

	middleware := NewGuestRateLimitMiddleware(mockRepo, log)
	handler := middleware.GuestOrAuthRateLimit("/api/v1/agent/reply", 3)

	// Mock repository to return existing usage (2/3)
	lifetimeWindow := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	existingUsage := &guest_usage.GuestAPIUsage{
		ID:            "1",
		CompositeKey:  "test_composite_key",
		Endpoint:      "/api/v1/agent/reply",
		UsageCount:    2,
		DailyLimit:    3,
		WindowResetAt: lifetimeWindow,
		IsBlocked:     false,
	}
	mockRepo.On("GetByCompositeKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(existingUsage, nil)
	mockRepo.On("IncrementUsage", mock.Anything, "1").Return(nil)

	app.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("user_type", "guest")
		return c.Next()
	}, handler, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("User-Agent", "TestBrowser/1.0")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify IncrementUsage was called
	mockRepo.AssertCalled(t, "IncrementUsage", mock.Anything, "1")
}

// TestGuestOrAuthRateLimit_LoggedInNoReferral_MultipleRequests tests rate limiting across multiple requests
func TestGuestOrAuthRateLimit_LoggedInNoReferral_MultipleRequests(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGuestUsageRepository)
	log := &MockLogger{}

	middleware := NewGuestRateLimitMiddleware(mockRepo, log)
	handler := middleware.GuestOrAuthRateLimit("/api/v1/agent/reply", 3)

	// Simulate 4 requests from the same logged-in user
	// Requests 1-3 should succeed, request 4 should fail

	// Request 1: First request - no existing record
	mockRepo.On("GetByCompositeKey", mock.Anything, "user_test_uid_456", mock.Anything, mock.Anything).Return(nil, nil).Once()
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(usage *guest_usage.GuestAPIUsage) bool {
		return usage.UsageCount == 1 && usage.CompositeKey == "user_test_uid_456"
	})).Return(nil).Once()

	app.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("user_type", "logged_in_no_referral")
		c.Locals("user", "mock_user")
		c.Locals("firebase_uid", "test_uid_456")
		return c.Next()
	}, handler, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req1 := httptest.NewRequest("POST", "/test", nil)
	resp1, _ := app.Test(req1)
	assert.Equal(t, 200, resp1.StatusCode, "Request 1 should succeed")

	// Request 2: Second request - existing record with count=1
	existingUsage2 := &guest_usage.GuestAPIUsage{
		ID:            "test_id_1",
		CompositeKey:  "user_test_uid_456",
		UsageCount:    1,
		DailyLimit:    3,
		WindowResetAt: guest_usage.GetNextResetTime(),
	}
	mockRepo.On("GetByCompositeKey", mock.Anything, "user_test_uid_456", mock.Anything, mock.Anything).Return(existingUsage2, nil).Once()
	mockRepo.On("IncrementUsage", mock.Anything, "test_id_1").Return(nil).Once()

	req2 := httptest.NewRequest("POST", "/test", nil)
	resp2, _ := app.Test(req2)
	assert.Equal(t, 200, resp2.StatusCode, "Request 2 should succeed")

	// Request 3: Third request - existing record with count=2
	existingUsage3 := &guest_usage.GuestAPIUsage{
		ID:            "test_id_1",
		CompositeKey:  "user_test_uid_456",
		UsageCount:    2,
		DailyLimit:    3,
		WindowResetAt: guest_usage.GetNextResetTime(),
	}
	mockRepo.On("GetByCompositeKey", mock.Anything, "user_test_uid_456", mock.Anything, mock.Anything).Return(existingUsage3, nil).Once()
	mockRepo.On("IncrementUsage", mock.Anything, "test_id_1").Return(nil).Once()

	req3 := httptest.NewRequest("POST", "/test", nil)
	resp3, _ := app.Test(req3)
	assert.Equal(t, 200, resp3.StatusCode, "Request 3 should succeed")

	// Request 4: Fourth request - existing record with count=3 (limit reached)
	existingUsage4 := &guest_usage.GuestAPIUsage{
		ID:            "test_id_1",
		CompositeKey:  "user_test_uid_456",
		UsageCount:    3,
		DailyLimit:    3,
		WindowResetAt: guest_usage.GetNextResetTime(),
	}
	mockRepo.On("GetByCompositeKey", mock.Anything, "user_test_uid_456", mock.Anything, mock.Anything).Return(existingUsage4, nil).Once()

	req4 := httptest.NewRequest("POST", "/test", nil)
	resp4, _ := app.Test(req4)
	assert.Equal(t, 429, resp4.StatusCode, "Request 4 should be blocked")

	var body map[string]interface{}
	json.NewDecoder(resp4.Body).Decode(&body)
	assert.Equal(t, "Daily limit exceeded", body["error"])
	assert.Contains(t, body["message"], "activate a referral code")
}

// TestGuestOrAuthRateLimit_UnknownUserType tests fallback to guest behavior
func TestGuestOrAuthRateLimit_UnknownUserType(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockGuestUsageRepository)
	log := &MockLogger{}

	middleware := NewGuestRateLimitMiddleware(mockRepo, log)
	handler := middleware.GuestOrAuthRateLimit("/api/v1/agent/reply", 3)

	// Mock repository
	mockRepo.On("GetByCompositeKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	app.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("user_type", "unknown_type")
		return c.Next()
	}, handler, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("User-Agent", "TestBrowser/1.0")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Should fall back to guest lifetime behavior
	mockRepo.AssertCalled(t, "Create", mock.Anything, mock.MatchedBy(func(usage *guest_usage.GuestAPIUsage) bool {
		return usage.WindowResetAt.Year() == 9999
	}))
}
