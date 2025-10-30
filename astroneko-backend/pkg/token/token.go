package token

import (
	"astroneko-backend/configs"
	"astroneko-backend/internal/core/domain/shared"
	"fmt"
	"time"

	"astroneko-backend/internal/core/domain/auth"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type TokenType string

const (
	AccessTokenConst  TokenType = "access_token"
	RefreshTokenConst TokenType = "refresh_token"
)

type (
	Token struct {
		AccessToken  *string `json:"access_token,omitempty" mapstructure:"access_token"`
		RefreshToken *string `json:"refresh_token,omitempty" mapstructure:"refresh_token"`
	}

	TokenResponse struct {
		Iss     *string       `json:"iss" mapstructure:"iss"`
		Sub     *string       `json:"sub" mapstructure:"sub"`
		Exp     *int64        `json:"exp" mapstructure:"exp"`
		Iat     *int64        `json:"iat" mapstructure:"iat"`
		Payload *TokenPayload `json:"payload,omitempty" mapstructure:"payload"`
		TokenID *string       `json:"tokenID,omitempty" mapstructure:"tokenID"`
	}

	TokenPayload struct {
		UserID *string `json:"user_id,omitempty" mapstructure:"user_id"`
		Email  *string `json:"email,omitempty" mapstructure:"email"`
		RoleID *int64  `json:"role_id,omitempty" mapstructure:"role_id"`
		Role   *string `json:"role,omitempty" mapstructure:"role"`
	}
)

const (
	AccessTokenDuration  = time.Minute * 15
	RefreshTokenDuration = time.Hour * 24 * 7
)

func NewToken(payload *auth.TokenPayloadRequest) (*Token, error) {
	var token Token

	// Create Access Token
	if err := token.CreateToken(payload, AccessTokenConst); err != nil {
		return nil, fmt.Errorf("%w: %w", shared.ErrInvalidAccessToken, err)
	}

	// Create Refresh Token
	if err := token.CreateToken(payload, RefreshTokenConst); err != nil {
		return nil, fmt.Errorf("%w: %w", shared.ErrInvalidRefreshToken, err)
	}

	return &token, nil
}

func (t *Token) CreateToken(payload *auth.TokenPayloadRequest, op TokenType) error {
	tokenDuration := t.getTokenDuration(op)
	claims := t.createClaims(payload, tokenDuration)

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	signedToken, err := token.SignedString([]byte(configs.GetViper().App.JWT))
	if err != nil {
		return err
	}

	t.setToken(op, &signedToken)
	return nil
}

func (t *Token) DecodeToken(tokenType TokenType) (*TokenResponse, error) {
	var tokenString *string

	switch tokenType {
	case AccessTokenConst:
		if t.AccessToken == nil {
			return nil, shared.ErrInvalidAccessToken
		}
		tokenString = t.AccessToken
	case RefreshTokenConst:
		if t.RefreshToken == nil {
			return nil, shared.ErrInvalidRefreshToken
		}
		tokenString = t.RefreshToken
	default:
		return nil, fmt.Errorf("%w: %s", shared.ErrInvalidTokenType, tokenType)
	}

	tokenPayload, err := ParseAndValidateToken(*tokenString)
	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", shared.ErrInvalidTokenType, tokenType, err)
	}

	return tokenPayload, nil
}

func (t *Token) DecodeAccessToken() (*TokenResponse, error) {
	return t.DecodeToken(AccessTokenConst)
}

func (t *Token) DecodeRefreshToken() (*TokenResponse, error) {
	return t.DecodeToken(RefreshTokenConst)
}

func ParseAndValidateToken(tokenString string) (*TokenResponse, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", shared.ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return []byte(configs.GetViper().App.JWT), nil
	})

	if err != nil {
		return nil, err
	}

	var tokenPayload TokenResponse
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &tokenPayload,
		WeaklyTypedInput: true,
		TagName:          "mapstructure",
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", shared.ErrInvalidTokenType, err)
	}

	if err := decoder.Decode(claims); err != nil {
		return nil, fmt.Errorf("%w: %w", shared.ErrInvalidTokenType, err)
	}

	return &tokenPayload, nil
}

func (t *Token) getTokenDuration(op TokenType) time.Duration {
	if op == AccessTokenConst {
		return AccessTokenDuration
	}
	return RefreshTokenDuration
}

func (t *Token) createClaims(payload *auth.TokenPayloadRequest, duration time.Duration) jwt.MapClaims {
	now := time.Now()
	issuer := fmt.Sprintf("%s-%s", configs.GetViper().App.Project, configs.GetViper().App.Env)

	return jwt.MapClaims{
		"iss":     issuer,
		"sub":     payload.Email,
		"exp":     now.Add(duration).Unix(),
		"iat":     now.Unix(),
		"payload": payload,
		"tokenID": uuid.New().String(),
	}
}

func (t *Token) setToken(op TokenType, token *string) {
	if op == AccessTokenConst {
		t.AccessToken = token
	} else {
		t.RefreshToken = token
	}
}
