package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type SecurityConfig struct {
	CSRFSecret    string
	CookieSecret  string
	SessionSecret string
	Domain        string
}

// GenerateSecretKey generates a random 32-byte secret key encoded as base64
func GenerateSecretKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a base64-encoded 32-byte default secret if random generation fails
		return base64.StdEncoding.EncodeToString([]byte("astroneko-default-secret-32byte!"))
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// SetupHelmetMiddleware configures security headers
func SetupHelmetMiddleware() fiber.Handler {
	return helmet.New(helmet.Config{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",
		HSTSMaxAge:         31536000,
		ReferrerPolicy:     "no-referrer",
	})
}

// SetupRateLimitMiddleware configures rate limiting
func SetupRateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        500,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if forwardedFor := c.Get("x-forwarded-for"); forwardedFor != "" {
				return forwardedFor
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, please try again later",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

// SetupEncryptCookieMiddleware configures cookie encryption
func SetupEncryptCookieMiddleware(config SecurityConfig) fiber.Handler {
	return encryptcookie.New(encryptcookie.Config{
		Key:    config.CookieSecret,
		Except: []string{"csrf_"}, // Don't encrypt CSRF tokens
	})
}

// SetupSessionMiddleware configures session management
func SetupSessionMiddleware(config SecurityConfig) *session.Store {
	return session.New(session.Config{
		KeyLookup:      "cookie:session_id",
		CookieDomain:   config.Domain,
		CookiePath:     "/",
		CookieSecure:   true, // Should be true in production with HTTPS
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		Expiration:     24 * time.Hour,
		KeyGenerator: func() string {
			return GenerateSecretKey()
		},
	})
}

// SetupCSRFMiddleware configures CSRF protection
func SetupCSRFMiddleware(config SecurityConfig, sessionStore *session.Store) fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_token",
		CookieSameSite: "Lax",
		CookieSecure:   true,  // Should be true in production with HTTPS
		CookieHTTPOnly: false, // Must be false so JavaScript can read it
		CookieDomain:   config.Domain,
		Expiration:     1 * time.Hour,
		KeyGenerator: func() string {
			return GenerateSecretKey()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "CSRF token validation failed",
			})
		},
		Extractor: func(c *fiber.Ctx) (string, error) {
			// Try header first
			token := c.Get("X-Csrf-Token")
			if token != "" {
				return token, nil
			}
			// Fall back to form field
			return c.FormValue("csrf_token"), nil
		},
		Session:    sessionStore,
		SessionKey: "csrf_token",
	})
}

// SetupAgentReplyRateLimitMiddleware configures rate limiting for agent reply endpoint
func SetupAgentReplyRateLimitMiddleware() fiber.Handler {
	limit := 500
	expiration := 1 * time.Minute
	return limiter.New(limiter.Config{
		Max:        limit,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			if forwardedFor := c.Get("x-forwarded-for"); forwardedFor != "" {
				return forwardedFor
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": "Now too many people use our Astroneko fortune teller, please wait a few minutes and use it again",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

// SetupXSSProtectionMiddleware adds additional XSS protection headers
func SetupXSSProtectionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline' 'unsafe-eval' *.googleapis.com *.gstatic.com; "+
				"style-src 'self' 'unsafe-inline' *.googleapis.com; "+
				"img-src 'self' data: *.googleusercontent.com *.googleapis.com; "+
				"connect-src 'self' *.googleapis.com *.firebase.com; "+
				"font-src 'self' *.gstatic.com; "+
				"object-src 'none'; "+
				"media-src 'self'; "+
				"frame-src 'self' *.firebase.com; "+
				"base-uri 'self';")
		return c.Next()
	}
}

// SetupSecureHeadersMiddleware adds additional security headers
func SetupSecureHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Feature-Policy", "geolocation 'none'; microphone 'none'; camera 'none'")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		return c.Next()
	}
}
