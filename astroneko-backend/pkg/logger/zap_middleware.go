package logger

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

const (
	LogMessageHTTPRequest = "HTTP Request"
)

var sensitiveFields = map[string]struct{}{
	"password":     {},
	"token":        {},
	"accessToken":  {},
	"refreshToken": {},
	"secret":       {},
	"ssn":          {},
	"email":        {},
}

func ZapLoggerMiddleware(zapLogger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		reqBody := string(c.Body())
		maskedBody := MaskSensitiveData(reqBody)

		err := c.Next()

		statusCode := c.Response().StatusCode()
		latency := time.Since(start)
		respBody := string(c.Response().Body())

		// Fields for Zap logger (stdout)
		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
			zap.String("query", c.Context().QueryArgs().String()),
			zap.String("requestBody", truncate(maskedBody, 300)),
			zap.Int("status", statusCode),
			zap.String("latency", latency.String()),
			zap.String("ip", c.IP()),
		}

		if statusCode >= 400 {
			fields = append(fields, zap.String("responseBody", truncate(respBody, 300)))
		}

		// Write to stdout (zap)
		switch {
		case statusCode >= 500:
			zapLogger.Error(LogMessageHTTPRequest, fields...)
		case statusCode >= 400:
			zapLogger.Warn(LogMessageHTTPRequest, fields...)
		default:
			zapLogger.Info(LogMessageHTTPRequest, fields...)
		}

		return err
	}
}

func ZapRecoveryMiddleware(logger *zap.Logger) fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {

			reqBody := string(c.Body())
			maskedBody := MaskSensitiveData(reqBody)

			// Convert panic to error
			var err error
			switch x := e.(type) {
			case string:
				err = fmt.Errorf("%s", x)
			case error:
				err = x
			default:
				err = fmt.Errorf("panic: %v", e)
			}

			logger.Error("Recovered from panic",
				zap.Any("error", err),
				zap.ByteString("stack", debug.Stack()),
				zap.String("url", c.OriginalURL()),
				zap.String("method", c.Method()),
				zap.String("ip", c.IP()),
				zap.String("requestBody", truncate(maskedBody, 300)),
			)
		},
	})
}

func truncate(s string, limit int) string {
	if len(s) > limit {
		return s[:limit] + "..."
	}
	return s
}

func MaskSensitiveData(jsonStr string) string {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return jsonStr
	}
	maskMap(raw)
	masked, err := json.Marshal(raw)
	if err != nil {
		return jsonStr
	}
	return string(masked)
}

func maskMap(data map[string]interface{}) {
	for k, v := range data {
		if _, found := sensitiveFields[k]; found {
			data[k] = "***"
		} else if subMap, ok := v.(map[string]interface{}); ok {
			maskMap(subMap)
		} else if subSlice, ok := v.([]interface{}); ok {
			maskSlice(subSlice)
		}
	}
}

func maskSlice(arr []interface{}) {
	for i, v := range arr {
		if subMap, ok := v.(map[string]interface{}); ok {
			maskMap(subMap)
			arr[i] = subMap
		}
	}
}
