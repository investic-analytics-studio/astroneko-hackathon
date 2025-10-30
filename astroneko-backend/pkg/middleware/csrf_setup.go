package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// CSRFManager manages CSRF protection for the application
type CSRFManager struct {
	sessionStore *session.Store
	config       SecurityConfig
}

// NewCSRFManager creates a new CSRF manager
func NewCSRFManager(config SecurityConfig) *CSRFManager {
	sessionStore := SetupSessionMiddleware(config)
	return &CSRFManager{
		sessionStore: sessionStore,
		config:       config,
	}
}

// GetCSRFMiddleware returns CSRF middleware for protecting routes
func (cm *CSRFManager) GetCSRFMiddleware() fiber.Handler {
	return SetupCSRFMiddleware(cm.config, cm.sessionStore)
}

// GetCSRFToken endpoint to get CSRF token for frontend
func (cm *CSRFManager) GetCSRFToken(c *fiber.Ctx) error {
	// Get session
	sess, err := cm.sessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get session",
		})
	}
	defer sess.Save()

	// Generate CSRF token
	token := GenerateSecretKey()
	sess.Set("csrf_token", token)

	// Set CSRF cookie
	c.Cookie(&fiber.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Domain:   cm.config.Domain,
		Path:     "/",
		MaxAge:   3600,  // 1 hour
		Secure:   true,  // Should be true in production with HTTPS
		HTTPOnly: false, // Must be false so JavaScript can read it
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"csrf_token": token,
	})
}

// SetupCSRFRoutes sets up CSRF-related routes
func SetupCSRFRoutes(app *fiber.App, csrfManager *CSRFManager) {
	// CSRF token endpoint
	app.Get("/csrf-token", csrfManager.GetCSRFToken)
}
