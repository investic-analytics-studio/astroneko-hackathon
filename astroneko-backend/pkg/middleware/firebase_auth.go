package middleware

import (
	"context"
	"fmt"
	"strings"

	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/internal/core/domain/user"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/logger"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

const (
	BearerPrefix = "Bearer "
	ModuleName   = "auth_middleware"
)

type AuthMiddleware struct {
	firebaseApp *auth.Client
	userService *services.UserService
	logger      logger.Logger
}

func NewFirebaseAuthMiddleware(firebaseApp *auth.Client, userService *services.UserService, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		firebaseApp: firebaseApp,
		userService: userService,
		logger:      log,
	}
}

// extractTokenFromHeader extracts and validates the Bearer token from the Authorization header
func (m *AuthMiddleware) extractTokenFromHeader(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header required")
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	idToken := strings.TrimPrefix(authHeader, BearerPrefix)
	if idToken == "" {
		return "", fmt.Errorf("token not provided")
	}

	return idToken, nil
}

// verifyToken verifies the Firebase ID token and returns the token claims
func (m *AuthMiddleware) verifyToken(idToken string) (*auth.Token, error) {
	ctx := context.Background()
	token, err := m.firebaseApp.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	return token, nil
}

// getUserFromToken retrieves the user from database using Firebase UID
func (m *AuthMiddleware) getUserFromToken(c *fiber.Ctx, token *auth.Token) (*user.User, error) {
	user, err := m.userService.GetUserByFirebaseUID(c.Context(), token.UID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// setUserContext adds user data to the Fiber context
func (m *AuthMiddleware) setUserContext(c *fiber.Ctx, user *user.User, token *auth.Token) {
	c.Locals("user", user)
	c.Locals("firebase_token", token)
	c.Locals("firebase_uid", token.UID)
	if email, ok := token.Claims["email"]; ok {
		c.Locals("user_email", email)
	}
}

// RequireAuth middleware validates Firebase ID token and adds user data to context
func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	if m.firebaseApp == nil {
		m.logger.Error("Firebase not configured in middleware", logger.Field{Key: "module", Value: ModuleName})
		status, response := shared.NewErrorResponse("ERR_1001", "Firebase not configured")
		return c.Status(status).JSON(response)
	}

	idToken, err := m.extractTokenFromHeader(c)
	if err != nil {
		status, response := shared.NewErrorResponse("ERR_1014", err.Error())
		return c.Status(status).JSON(response)
	}

	token, err := m.verifyToken(idToken)
	if err != nil {
		m.logger.Error("Failed to verify ID token",
			logger.Field{Key: "module", Value: ModuleName},
			logger.Field{Key: "error", Value: err.Error()})
		status, response := shared.NewErrorResponse("ERR_1003", "Invalid or expired token")
		return c.Status(status).JSON(response)
	}

	user, err := m.getUserFromToken(c, token)
	if err != nil {
		m.logger.Error("User not found for Firebase UID",
			logger.Field{Key: "module", Value: ModuleName},
			logger.Field{Key: "firebase_uid", Value: token.UID},
			logger.Field{Key: "error", Value: err.Error()})
		status, response := shared.NewErrorResponse("ERR_1023", "User not found")
		return c.Status(status).JSON(response)
	}

	m.setUserContext(c, user, token)
	return c.Next()
}

// OptionalAuth middleware validates Firebase ID token if provided but doesn't require it
func (m *AuthMiddleware) OptionalAuth(c *fiber.Ctx) error {
	if m.firebaseApp == nil {
		return c.Next()
	}

	idToken, err := m.extractTokenFromHeader(c)
	if err != nil {
		return c.Next() // Continue without authentication
	}

	token, err := m.verifyToken(idToken)
	if err != nil {
		m.logger.Warn("Failed to verify optional ID token",
			logger.Field{Key: "module", Value: ModuleName},
			logger.Field{Key: "error", Value: err.Error()})
		return c.Next() // Continue without authentication
	}

	user, err := m.getUserFromToken(c, token)
	if err != nil {
		m.logger.Warn("User not found for optional authentication",
			logger.Field{Key: "module", Value: ModuleName},
			logger.Field{Key: "firebase_uid", Value: token.UID},
			logger.Field{Key: "error", Value: err.Error()})
		return c.Next() // Continue without user data
	}

	m.setUserContext(c, user, token)

	return c.Next()
}

func (m *AuthMiddleware) OptionalAuthWithReferralCheck(c *fiber.Ctx) error {
	if m.firebaseApp == nil {
		c.Locals("user_type", "guest")
		return c.Next()
	}

	idToken, err := m.extractTokenFromHeader(c)
	if err != nil {
		c.Locals("user_type", "guest")
		return c.Next() // Continue without authentication
	}

	token, err := m.verifyToken(idToken)
	if err != nil {
		m.logger.Warn("Failed to verify ID token for agent API",
			logger.Field{Key: "module", Value: ModuleName},
			logger.Field{Key: "error", Value: err.Error()})
		c.Locals("user_type", "guest")
		return c.Next() // Continue without authentication
	}

	user, err := m.getUserFromToken(c, token)
	if err != nil {
		m.logger.Warn("User not found for agent API authentication",
			logger.Field{Key: "module", Value: ModuleName},
			logger.Field{Key: "firebase_uid", Value: token.UID},
			logger.Field{Key: "error", Value: err.Error()})
		c.Locals("user_type", "guest")
		return c.Next() // Continue without user data
	}

	// CRITICAL: Check if referral is activated for agent API access
	if !user.IsActivatedReferral {
		m.logger.Info("Logged-in user without activated referral accessing agent API",
			logger.Field{Key: "module", Value: ModuleName},
			logger.Field{Key: "user_id", Value: user.ID.String()},
			logger.Field{Key: "email", Value: user.Email})
		// Set user context for logged-in users without referral (3 requests per day)
		m.setUserContext(c, user, token)
		c.Locals("user_type", "logged_in_no_referral")
		return c.Next()
	}

	// User has activated referral - grant unlimited access
	m.setUserContext(c, user, token)
	c.Locals("user_type", "logged_in_with_referral")

	return c.Next()
}
