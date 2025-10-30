package firebase

import (
	"context"
	"encoding/json"
	"fmt"

	"astroneko-backend/configs"
	"astroneko-backend/internal/core/domain/shared"
	"astroneko-backend/pkg/apprequest"
	"astroneko-backend/pkg/logger"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"google.golang.org/api/option"
)

var (
	FirebaseClient *auth.Client
	firebaseLogger logger.Logger
	httpClient     apprequest.HTTPRequest
)

func InitFirebaseClient(log logger.Logger) error {
	firebaseLogger = log
	httpClient = apprequest.NewRequester()

	ctx := context.Background()
	cfg := configs.GetViper().Firebase

	if cfg.Credential == "" {
		firebaseLogger.Error("Firebase credential is empty", logger.Field{Key: "module", Value: "firebase"})
		return shared.ErrFirebaseCredentialInvalid
	}

	credential := []byte(cfg.Credential)
	clientOption := option.WithCredentialsJSON(credential)

	app, err := firebase.NewApp(ctx, nil, clientOption)
	if err != nil {
		firebaseLogger.Error("Failed to initialize Firebase app",
			logger.Field{Key: "module", Value: "firebase"},
			logger.Field{Key: "error", Value: err.Error()})
		return fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	FirebaseClient, err = app.Auth(ctx)
	if err != nil {
		firebaseLogger.Error("Failed to initialize Firebase Auth client",
			logger.Field{Key: "module", Value: "firebase"},
			logger.Field{Key: "error", Value: err.Error()})
		return fmt.Errorf("failed to initialize Firebase Auth client: %w", err)
	}

	firebaseLogger.Info("Firebase client initialized successfully",
		logger.Field{Key: "module", Value: "firebase"})
	return nil
}

type RefreshTokenRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
	ProjectID    string `json:"project_id"`
}

// RefreshFirebaseToken refreshes a Firebase refresh token and returns new tokens
func RefreshFirebaseToken(refreshToken string) (*RefreshTokenResponse, error) {
	if FirebaseClient == nil {
		firebaseLogger.Error("Firebase client not initialized", logger.Field{Key: "module", Value: "firebase"})
		return nil, shared.ErrFirebaseNotInitialized
	}

	cfg := configs.GetViper().Firebase
	if cfg.WebAPIKey == "" {
		firebaseLogger.Error("Firebase web API key not configured", logger.Field{Key: "module", Value: "firebase"})
		return nil, shared.ErrFirebaseWebAPIKeyMissing
	}

	// Refreshing token - no logging needed

	url := fmt.Sprintf("https://securetoken.googleapis.com/v1/token?key=%s", cfg.WebAPIKey)
	reqBody := RefreshTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		firebaseLogger.Error("Failed to marshal refresh token request",
			logger.Field{Key: "module", Value: "firebase"},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, resp := httpClient.NewRequest(jsonBody, apprequest.POST, url)
	req.Header.SetContentTypeBytes(apprequest.ApplicationJSON)

	{
		err = fasthttp.Do(req, resp)
		if err != nil {
			logrus.Error("Failed to make refresh token request: ", err)
			firebaseLogger.Error("Failed to make refresh token request",
				logger.Field{Key: "module", Value: "firebase"},
				logger.Field{Key: "error", Value: err.Error()})
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
	}

	{
		bodyBytes := resp.Body()
		if resp.StatusCode() != fasthttp.StatusOK {
			firebaseLogger.Error("Firebase token refresh failed",
				logger.Field{Key: "module", Value: "firebase"},
				logger.Field{Key: "status_code", Value: resp.StatusCode()},
				logger.Field{Key: "response", Value: string(bodyBytes)})
			return nil, shared.ErrFirebaseTokenRefreshFailed
		}

		var tokenResp RefreshTokenResponse
		if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
			logrus.Error("error on unmarshal: ", err)
			firebaseLogger.Error("Failed to decode refresh token response",
				logger.Field{Key: "module", Value: "firebase"},
				logger.Field{Key: "error", Value: err.Error()})
			return nil, shared.ErrFirebaseTokenDecodeFailed
		}

		// Token refreshed - no logging needed
		return &tokenResp, nil
	}
}

type LoginRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type LoginResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

func SignInWithEmailPassword(email, password string) (*LoginResponse, error) {
	if FirebaseClient == nil {
		firebaseLogger.Error("Firebase client not initialized", logger.Field{Key: "module", Value: "firebase"})
		return nil, shared.ErrFirebaseNotInitialized
	}

	cfg := configs.GetViper().Firebase
	if cfg.WebAPIKey == "" {
		firebaseLogger.Error("Firebase web API key not configured", logger.Field{Key: "module", Value: "firebase"})
		return nil, shared.ErrFirebaseWebAPIKeyMissing
	}

	// Signing in - no logging needed

	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", cfg.WebAPIKey)
	reqBody := map[string]any{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		firebaseLogger.Error("Failed to marshal sign in request",
			logger.Field{Key: "module", Value: "firebase"},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, resp := httpClient.NewRequest(jsonBody, apprequest.POST, url)
	req.Header.SetContentTypeBytes(apprequest.ApplicationJSON)

	{
		err = fasthttp.Do(req, resp)
		if err != nil {
			logrus.Error("Failed to make sign in request: ", err)
			firebaseLogger.Error("Failed to make sign in request",
				logger.Field{Key: "module", Value: "firebase"},
				logger.Field{Key: "error", Value: err.Error()})
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
	}

	{
		bodyBytes := resp.Body()
		if resp.StatusCode() != fasthttp.StatusOK {
			firebaseLogger.Error("Firebase authentication failed",
				logger.Field{Key: "module", Value: "firebase"},
				logger.Field{Key: "status_code", Value: resp.StatusCode()},
				logger.Field{Key: "response", Value: string(bodyBytes)})
			return nil, shared.ErrFirebaseAuthFailed
		}

		var authResp LoginResponse
		if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
			logrus.Error("error on unmarshal: ", err)
			firebaseLogger.Error("Failed to decode sign in response",
				logger.Field{Key: "module", Value: "firebase"},
				logger.Field{Key: "error", Value: err.Error()})
			return nil, shared.ErrFirebaseTokenDecodeFailed
		}

		// Auth successful - no logging needed
		return &authResp, nil
	}
}
