package middleware

import (
	"context"
	"fmt"
	"time"

	"astroneko-backend/internal/core/domain/guest_usage"
	guestUsagePort "astroneko-backend/internal/core/ports/guest_usage"
	"astroneko-backend/pkg/logger"
	"astroneko-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	DefaultGuestDailyLimit = 3
	AgentReplyEndpoint     = "/api/v1/agent/reply"
)

type GuestRateLimitMiddleware struct {
	guestRepo guestUsagePort.Repository
	logger    logger.Logger
}

func NewGuestRateLimitMiddleware(repo guestUsagePort.Repository, log logger.Logger) *GuestRateLimitMiddleware {
	return &GuestRateLimitMiddleware{
		guestRepo: repo,
		logger:    log,
	}
}

// GuestOrAuthRateLimit applies rate limiting based on user type:
// - logged_in_with_referral: Unlimited access
// - logged_in_no_referral: 3 requests per day (daily reset)
// - guest: 3 requests lifetime (no reset)
func (m *GuestRateLimitMiddleware) GuestOrAuthRateLimit(endpoint string, guestLimit int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType, _ := c.Locals("user_type").(string)

		switch userType {
		case "logged_in_with_referral":
			// Unlimited access for users with activated referral
			m.logger.Info("User with activated referral accessing endpoint",
				logger.Field{Key: "endpoint", Value: endpoint},
				logger.Field{Key: "user_type", Value: userType})
			return c.Next()

		case "logged_in_no_referral":
			// Logged-in users without referral: 3 requests per day
			m.logger.Info("Logged-in user without referral accessing endpoint",
				logger.Field{Key: "endpoint", Value: endpoint},
				logger.Field{Key: "user_type", Value: userType})
			return m.handleLoggedInUserRequest(c, endpoint, guestLimit)

		case "guest", "":
			// Guest users: 3 requests lifetime (no daily reset)
			m.logger.Info("Guest user accessing endpoint",
				logger.Field{Key: "endpoint", Value: endpoint},
				logger.Field{Key: "user_type", Value: userType})
			return m.handleGuestLifetimeLimit(c, endpoint, guestLimit)

		default:
			// Fallback: treat as guest
			m.logger.Warn("Unknown user type, treating as guest",
				logger.Field{Key: "endpoint", Value: endpoint},
				logger.Field{Key: "user_type", Value: userType})
			return m.handleGuestLifetimeLimit(c, endpoint, guestLimit)
		}
	}
}

// handleLoggedInUserRequest handles rate limiting for logged-in users without referral (3 per day with daily reset)
func (m *GuestRateLimitMiddleware) handleLoggedInUserRequest(c *fiber.Ctx, endpoint string, limit int) error {
	ctx := context.Background()

	// Get user from context
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context not found",
		})
	}

	// Use user ID as the composite key for logged-in users
	userID := c.Locals("firebase_uid").(string)
	compositeKey := "user_" + userID

	// CRITICAL: Use fixed daily window (next midnight UTC) so the same windowResetStr is used throughout the day
	// This ensures GetByCompositeKey finds existing records instead of creating new ones every request
	windowResetAt := guest_usage.GetNextResetTime() // Next midnight UTC
	windowResetStr := utils.GetWindowResetString(windowResetAt)

	usage, err := m.guestRepo.GetByCompositeKey(ctx, compositeKey, endpoint, windowResetStr)
	if err != nil {
		m.logger.Error("Failed to get logged-in user usage",
			logger.Field{Key: "composite_key", Value: compositeKey},
			logger.Field{Key: "error", Value: err.Error()})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check rate limit",
		})
	}

	// First request from this user today
	if usage == nil {
		usage = &guest_usage.GuestAPIUsage{
			IPAddress:     c.IP(),
			UserAgentHash: "logged_in_user",
			CompositeKey:  compositeKey,
			Endpoint:      endpoint,
			UsageCount:    1,
			DailyLimit:    limit,
			WindowResetAt: windowResetAt,
			LastRequestAt: time.Now(),
			IsBlocked:     false,
		}

		if err := m.guestRepo.Create(ctx, usage); err != nil {
			m.logger.Error("Failed to create logged-in user usage",
				logger.Field{Key: "composite_key", Value: compositeKey},
				logger.Field{Key: "error", Value: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to track usage",
			})
		}

		m.setRateLimitHeaders(c, usage)
		m.logger.Info("New logged-in user request",
			logger.Field{Key: "user_id", Value: userID},
			logger.Field{Key: "endpoint", Value: endpoint},
			logger.Field{Key: "remaining", Value: usage.RemainingRequests()})

		return c.Next()
	}

	// No need to check window expiration here because:
	// - We use GetNextResetTime() which always returns next midnight UTC
	// - If we're in a new day, GetByCompositeKey won't find the old record (different windowResetStr)
	// - A new record will be created automatically above (usage == nil case)

	// Check if limit exceeded
	if !usage.CanMakeRequest() {
		m.setRateLimitHeaders(c, usage)

		resetIn := time.Until(usage.WindowResetAt)
		hours := int(resetIn.Hours())
		minutes := int(resetIn.Minutes()) % 60

		m.logger.Warn("Logged-in user rate limit exceeded",
			logger.Field{Key: "user_id", Value: userID},
			logger.Field{Key: "endpoint", Value: endpoint},
			logger.Field{Key: "usage_count", Value: usage.UsageCount},
			logger.Field{Key: "daily_limit", Value: usage.DailyLimit})

		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":       "Daily limit exceeded",
			"message":     "You've used all your free daily requests. Please activate a referral code for unlimited access or try again tomorrow.",
			"used":        usage.UsageCount,
			"limit":       usage.DailyLimit,
			"reset_in":    resetIn.String(),
			"reset_hours": hours,
			"reset_mins":  minutes,
		})
	}

	// Increment usage count
	if err := m.guestRepo.IncrementUsage(ctx, usage.ID); err != nil {
		m.logger.Error("Failed to increment logged-in user usage",
			logger.Field{Key: "id", Value: usage.ID},
			logger.Field{Key: "error", Value: err.Error()})
	}

	usage.IncrementUsage()
	m.setRateLimitHeaders(c, usage)

	m.logger.Info("Logged-in user request allowed",
		logger.Field{Key: "user_id", Value: userID},
		logger.Field{Key: "usage_count", Value: usage.UsageCount},
		logger.Field{Key: "remaining", Value: usage.RemainingRequests()})

	return c.Next()
}

// handleGuestLifetimeLimit handles rate limiting for guest users (3 lifetime, NO daily reset)
func (m *GuestRateLimitMiddleware) handleGuestLifetimeLimit(c *fiber.Ctx, endpoint string, limit int) error {
	ctx := context.Background()

	// Generate multi-factor fingerprint
	fingerprint := utils.GenerateGuestFingerprint(c)

	// For guests, use a fixed window that never resets (year 9999)
	lifetimeWindow := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	windowResetStr := utils.GetWindowResetString(lifetimeWindow)

	usage, err := m.guestRepo.GetByCompositeKey(ctx, fingerprint.CompositeKey, endpoint, windowResetStr)
	if err != nil {
		m.logger.Error("Failed to get guest lifetime usage",
			logger.Field{Key: "composite_key", Value: fingerprint.CompositeKey},
			logger.Field{Key: "error", Value: err.Error()})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check rate limit",
		})
	}

	// First request from this guest (ever)
	if usage == nil {
		usage = &guest_usage.GuestAPIUsage{
			IPAddress:     fingerprint.IPAddress,
			UserAgentHash: fingerprint.UserAgentHash,
			CompositeKey:  fingerprint.CompositeKey,
			Endpoint:      endpoint,
			UsageCount:    1,
			DailyLimit:    limit,
			WindowResetAt: lifetimeWindow, // Never resets
			LastRequestAt: time.Now(),
			IsBlocked:     false,
		}

		if err := m.guestRepo.Create(ctx, usage); err != nil {
			m.logger.Error("Failed to create guest lifetime usage",
				logger.Field{Key: "composite_key", Value: fingerprint.CompositeKey},
				logger.Field{Key: "error", Value: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to track usage",
			})
		}

		m.setRateLimitHeaders(c, usage)
		m.logger.Info("New guest lifetime request",
			logger.Field{Key: "ip", Value: fingerprint.IPAddress},
			logger.Field{Key: "endpoint", Value: endpoint},
			logger.Field{Key: "remaining", Value: usage.RemainingRequests()})

		return c.Next()
	}

	// Check if guest is blocked
	if usage.IsBlocked {
		m.logger.Warn("Blocked guest attempted request",
			logger.Field{Key: "composite_key", Value: fingerprint.CompositeKey},
			logger.Field{Key: "reason", Value: usage.BlockedReason})

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "Access denied",
			"reason": "Your access has been blocked due to suspicious activity",
		})
	}

	// Check if lifetime limit exceeded (NO RESET for guests)
	if !usage.CanMakeRequest() {
		m.setRateLimitHeaders(c, usage)

		m.logger.Warn("Guest lifetime limit exceeded",
			logger.Field{Key: "ip", Value: fingerprint.IPAddress},
			logger.Field{Key: "endpoint", Value: endpoint},
			logger.Field{Key: "usage_count", Value: usage.UsageCount},
			logger.Field{Key: "limit", Value: usage.DailyLimit})

		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":   "Free trial limit exceeded",
			"message": "You've used all 3 free requests. Please sign in for unlimited access.",
			"used":    usage.UsageCount,
			"limit":   usage.DailyLimit,
		})
	}

	// Increment usage count (no reset)
	if err := m.guestRepo.IncrementUsage(ctx, usage.ID); err != nil {
		m.logger.Error("Failed to increment guest lifetime usage",
			logger.Field{Key: "id", Value: usage.ID},
			logger.Field{Key: "error", Value: err.Error()})
	}

	usage.IncrementUsage()
	m.setRateLimitHeaders(c, usage)

	m.logger.Info("Guest lifetime request allowed",
		logger.Field{Key: "ip", Value: fingerprint.IPAddress},
		logger.Field{Key: "usage_count", Value: usage.UsageCount},
		logger.Field{Key: "remaining", Value: usage.RemainingRequests()})

	return c.Next()
}

// Legacy function - kept for backward compatibility but no longer used in agent routes
func (m *GuestRateLimitMiddleware) handleGuestRequest(c *fiber.Ctx, endpoint string, limit int) error {
	ctx := context.Background()

	// Generate multi-factor fingerprint
	fingerprint := utils.GenerateGuestFingerprint(c)

	// Get or create usage record
	windowResetStr := utils.GetWindowResetString(fingerprint.WindowResetAt)
	usage, err := m.guestRepo.GetByCompositeKey(ctx, fingerprint.CompositeKey, endpoint, windowResetStr)
	if err != nil {
		m.logger.Error("Failed to get guest usage",
			logger.Field{Key: "composite_key", Value: fingerprint.CompositeKey},
			logger.Field{Key: "error", Value: err.Error()})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check rate limit",
		})
	}

	// First request from this guest today
	if usage == nil {
		usage = &guest_usage.GuestAPIUsage{
			IPAddress:     fingerprint.IPAddress,
			UserAgentHash: fingerprint.UserAgentHash,
			CompositeKey:  fingerprint.CompositeKey,
			Endpoint:      endpoint,
			UsageCount:    1,
			DailyLimit:    limit,
			WindowResetAt: fingerprint.WindowResetAt,
			LastRequestAt: time.Now(),
			IsBlocked:     false,
		}

		if err := m.guestRepo.Create(ctx, usage); err != nil {
			m.logger.Error("Failed to create guest usage",
				logger.Field{Key: "composite_key", Value: fingerprint.CompositeKey},
				logger.Field{Key: "error", Value: err.Error()})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to track usage",
			})
		}

		// Add rate limit info to response headers
		m.setRateLimitHeaders(c, usage)

		m.logger.Info("New guest request",
			logger.Field{Key: "ip", Value: fingerprint.IPAddress},
			logger.Field{Key: "endpoint", Value: endpoint},
			logger.Field{Key: "remaining", Value: usage.RemainingRequests()})

		return c.Next()
	}

	// Check if guest is blocked
	if usage.IsBlocked {
		m.logger.Warn("Blocked guest attempted request",
			logger.Field{Key: "composite_key", Value: fingerprint.CompositeKey},
			logger.Field{Key: "reason", Value: usage.BlockedReason})

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "Access denied",
			"reason": "Your access has been blocked due to suspicious activity",
		})
	}

	// Check if window has expired (new day)
	if time.Now().After(usage.WindowResetAt) {
		// Reset the usage count
		usage.ResetWindow()
		usage.IncrementUsage()

		if err := m.guestRepo.Create(ctx, usage); err != nil {
			m.logger.Error("Failed to reset guest usage window",
				logger.Field{Key: "composite_key", Value: fingerprint.CompositeKey},
				logger.Field{Key: "error", Value: err.Error()})
		}

		m.setRateLimitHeaders(c, usage)
		return c.Next()
	}

	// Check if limit exceeded
	if !usage.CanMakeRequest() {
		m.setRateLimitHeaders(c, usage)

		resetIn := time.Until(usage.WindowResetAt)
		hours := int(resetIn.Hours())
		minutes := int(resetIn.Minutes()) % 60

		m.logger.Warn("Guest rate limit exceeded",
			logger.Field{Key: "ip", Value: fingerprint.IPAddress},
			logger.Field{Key: "endpoint", Value: endpoint},
			logger.Field{Key: "usage_count", Value: usage.UsageCount},
			logger.Field{Key: "daily_limit", Value: usage.DailyLimit})

		// Check if user is logged in but without activated referral
		authHeader := c.Get("Authorization")
		hasToken := authHeader != "" && len(authHeader) > 7 // "Bearer "

		var message string
		if hasToken {
			// Logged in user without activated referral
			message = "You are logged in. Please activate a referral code for unlimited access."
		} else {
			// Pure guest user
			message = "You've used all your free daily requests. Please sign in for unlimited access or try again tomorrow."
		}

		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":       "Daily limit exceeded",
			"message":     message,
			"used":        usage.UsageCount,
			"limit":       usage.DailyLimit,
			"reset_in":    resetIn.String(),
			"reset_hours": hours,
			"reset_mins":  minutes,
		})
	}

	// Increment usage count
	if err := m.guestRepo.IncrementUsage(ctx, usage.ID); err != nil {
		m.logger.Error("Failed to increment guest usage",
			logger.Field{Key: "id", Value: usage.ID},
			logger.Field{Key: "error", Value: err.Error()})
		// Don't fail the request on increment error
	}

	usage.IncrementUsage() // Update local copy for header
	m.setRateLimitHeaders(c, usage)

	m.logger.Info("Guest request allowed",
		logger.Field{Key: "ip", Value: fingerprint.IPAddress},
		logger.Field{Key: "usage_count", Value: usage.UsageCount},
		logger.Field{Key: "remaining", Value: usage.RemainingRequests()})

	return c.Next()
}

// setRateLimitHeaders adds standard rate limit headers to response
func (m *GuestRateLimitMiddleware) setRateLimitHeaders(c *fiber.Ctx, usage *guest_usage.GuestAPIUsage) {
	c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", usage.DailyLimit))
	c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", usage.RemainingRequests()))
	c.Set("X-RateLimit-Reset", usage.WindowResetAt.Format(time.RFC3339))
}

// SetupGuestAgentReplyRateLimit creates middleware specifically for agent reply endpoint
func SetupGuestAgentReplyRateLimit(repo guestUsagePort.Repository, log logger.Logger) fiber.Handler {
	middleware := NewGuestRateLimitMiddleware(repo, log)
	return middleware.GuestOrAuthRateLimit(AgentReplyEndpoint, DefaultGuestDailyLimit)
}

// AbuseDetectionMiddleware detects and blocks suspicious patterns
func (m *GuestRateLimitMiddleware) AbuseDetectionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user != nil {
			return c.Next() // Skip for authenticated users
		}

		ctx := context.Background()
		fingerprint := utils.GenerateGuestFingerprint(c)

		// Check if same IP has too many different composite keys (incognito bypass attempt)
		since := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		records, err := m.guestRepo.GetByIPAddress(ctx, fingerprint.IPAddress, since)
		if err == nil && len(records) > 10 {
			// Same IP with >10 different fingerprints in 24h = suspicious
			m.logger.Warn("Potential abuse detected - multiple fingerprints from same IP",
				logger.Field{Key: "ip", Value: fingerprint.IPAddress},
				logger.Field{Key: "fingerprint_count", Value: len(records)})

			// Block this composite key
			_ = m.guestRepo.BlockGuest(ctx, fingerprint.CompositeKey, "Multiple fingerprints from same IP")

			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Suspicious activity detected",
			})
		}

		return c.Next()
	}
}
