package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type TokenData struct {
	UserID    *string
	TokenID   *string
	IssuedAt  *time.Time
	ExpiresAt *time.Time
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func:
		return v.IsNil()
	}
	return false
}

// BuildMapIfNotNil creates map[string]interface{} from fields that are not nil
func BuildMapIfNotNil(pairs ...FieldPair) map[string]interface{} {
	result := make(map[string]interface{})
	for _, pair := range pairs {
		if !IsNil(pair.Value) {
			result[pair.Key] = pair.Value
		}
	}
	return result
}

type FieldPair struct {
	Key   string
	Value interface{}
}

func GetClientIP(c *fiber.Ctx) string {
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return c.IP()
}

func GenerateRefreshTokenHash(refreshToken string) string {
	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	return hex.EncodeToString(hasher.Sum(nil))
}

func DownloadFile(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func ConvertToUint(v int) (uint, error) {
	if v < 0 {
		return 0, fmt.Errorf("value %d cannot be negative", v)
	}
	return uint(v), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // DefaultCost = 10
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
