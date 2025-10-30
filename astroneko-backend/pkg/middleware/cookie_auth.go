package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// CookieAuthConfig holds configuration for cookie-based authentication
type CookieAuthConfig struct {
	Domain   string
	Secure   bool // Should be true in production with HTTPS
	HTTPOnly bool
	SameSite string
	Path     string
	MaxAge   int // in seconds
}

// AuthCookies contains the names for authentication cookies
type AuthCookies struct {
	AccessToken  string
	RefreshToken string
	UserID       string
}

// DefaultCookieAuthConfig returns default configuration for cookie authentication
func DefaultCookieAuthConfig(domain string) CookieAuthConfig {
	return CookieAuthConfig{
		Domain:   domain,
		Secure:   true,  // Should be true in production with HTTPS
		HTTPOnly: true,  // Prevent XSS attacks
		SameSite: "Lax", // CSRF protection while allowing some cross-site usage
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 days in seconds
	}
}

// DefaultAuthCookies returns default cookie names
func DefaultAuthCookies() AuthCookies {
	return AuthCookies{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		UserID:       "user_id",
	}
}

// SetAuthCookies sets authentication cookies on the response
func SetAuthCookies(c *fiber.Ctx, accessToken, refreshToken, userID string, config CookieAuthConfig, cookieNames AuthCookies) {
	// Set access token cookie (shorter expiry)
	c.Cookie(&fiber.Cookie{
		Name:     cookieNames.AccessToken,
		Value:    accessToken,
		Domain:   config.Domain,
		Path:     config.Path,
		MaxAge:   60 * 60, // 1 hour
		Secure:   config.Secure,
		HTTPOnly: config.HTTPOnly,
		SameSite: config.SameSite,
	})

	// Set refresh token cookie (longer expiry)
	c.Cookie(&fiber.Cookie{
		Name:     cookieNames.RefreshToken,
		Value:    refreshToken,
		Domain:   config.Domain,
		Path:     config.Path,
		MaxAge:   config.MaxAge, // 7 days
		Secure:   config.Secure,
		HTTPOnly: config.HTTPOnly,
		SameSite: config.SameSite,
	})
}

// GetAuthCookies retrieves authentication cookies from the request
func GetAuthCookies(c *fiber.Ctx, cookieNames AuthCookies) (accessToken, refreshToken, userID string) {
	accessToken = c.Cookies(cookieNames.AccessToken)
	refreshToken = c.Cookies(cookieNames.RefreshToken)
	userID = c.Cookies(cookieNames.UserID)
	return
}

// ClearAuthCookies removes authentication cookies
func ClearAuthCookies(c *fiber.Ctx, config CookieAuthConfig, cookieNames AuthCookies) {
	// Clear access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     cookieNames.AccessToken,
		Value:    "",
		Domain:   config.Domain,
		Path:     config.Path,
		MaxAge:   -1, // Expire immediately
		Secure:   config.Secure,
		HTTPOnly: config.HTTPOnly,
		SameSite: config.SameSite,
	})

	// Clear refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     cookieNames.RefreshToken,
		Value:    "",
		Domain:   config.Domain,
		Path:     config.Path,
		MaxAge:   -1, // Expire immediately
		Secure:   config.Secure,
		HTTPOnly: config.HTTPOnly,
		SameSite: config.SameSite,
	})
}

// ExtendedAuthMiddleware extends the existing Firebase auth middleware with cookie support
func (m *AuthMiddleware) WithCookieAuth(config CookieAuthConfig, cookieNames AuthCookies) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First try to get token from Authorization header (existing behavior)
		authHeader := c.Get("Authorization")
		var token string

		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			// Fall back to cookie-based authentication
			token, _, _ = GetAuthCookies(c, cookieNames)
		}

		// If no token found in either header or cookies, return unauthorized
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		// Set the token in the Authorization header for the existing middleware to process
		c.Set("Authorization", "Bearer "+token)

		// Call the existing Firebase auth middleware
		return m.RequireAuth(c)
	}
}
