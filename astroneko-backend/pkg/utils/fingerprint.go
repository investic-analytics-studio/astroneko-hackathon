package utils

import (
	"crypto/sha256"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GuestFingerprint represents a guest's unique fingerprint
type GuestFingerprint struct {
	IPAddress     string
	UserAgent     string
	UserAgentHash string
	CompositeKey  string
	WindowResetAt time.Time
}

// GenerateGuestFingerprint creates a multi-factor fingerprint to identify guests
// Combines IP + UserAgent + Daily Window to prevent incognito bypass
func GenerateGuestFingerprint(c *fiber.Ctx) *GuestFingerprint {
	ip := GetRealIP(c)
	userAgent := c.Get("User-Agent")

	// Normalize user agent (remove version numbers that change frequently)
	normalizedUA := NormalizeUserAgent(userAgent)

	// Hash the user agent for privacy
	uaHash := HashString(normalizedUA)

	// Get next reset time (midnight UTC)
	resetTime := GetNextMidnightUTC()
	resetDay := resetTime.Format("2006-01-02")

	// Create composite key: SHA256(IP + UserAgent + Day)
	// This prevents incognito bypass while allowing legitimate daily resets
	compositeKey := GenerateCompositeKey(ip, normalizedUA, resetDay)

	return &GuestFingerprint{
		IPAddress:     ip,
		UserAgent:     userAgent,
		UserAgentHash: uaHash,
		CompositeKey:  compositeKey,
		WindowResetAt: resetTime,
	}
}

// GetRealIP extracts the real client IP, handling proxies and load balancers
func GetRealIP(c *fiber.Ctx) string {
	// Priority order for IP detection:
	// 1. CF-Connecting-IP (Cloudflare)
	// 2. X-Real-IP (Nginx)
	// 3. X-Forwarded-For (Standard proxy header)
	// 4. c.IP() (Direct connection)

	if cfIP := c.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}

	if realIP := c.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	if forwardedFor := c.Get("X-Forwarded-For"); forwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs: "client, proxy1, proxy2"
		// Take the first (leftmost) IP which is the original client
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if isValidIP(clientIP) {
				return clientIP
			}
		}
	}

	return c.IP()
}

// NormalizeUserAgent removes version-specific info that changes frequently
// This helps identify the same browser even after updates
func NormalizeUserAgent(ua string) string {
	ua = strings.TrimSpace(ua)
	if ua == "" {
		return "unknown"
	}

	// Convert to lowercase for consistent hashing
	ua = strings.ToLower(ua)

	// Remove specific version numbers but keep major identifiers
	// Example: "Chrome/120.0.0.0" -> "Chrome"
	// This is a simplified version; production might use a UA parser library

	// Keep only major browser identifiers
	var normalized strings.Builder
	normalized.WriteString(detectBrowser(ua))
	normalized.WriteString(detectOS(ua))

	return normalized.String()
}

// detectBrowser identifies the browser type from user agent string
func detectBrowser(ua string) string {
	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edge") {
		return "chrome/"
	}
	if strings.Contains(ua, "firefox") {
		return "firefox/"
	}
	if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		return "safari/"
	}
	if strings.Contains(ua, "edge") {
		return "edge/"
	}
	if strings.Contains(ua, "opera") {
		return "opera/"
	}
	return "other/"
}

// detectOS identifies the operating system from user agent string
func detectOS(ua string) string {
	if strings.Contains(ua, "windows") {
		return "windows"
	}
	if strings.Contains(ua, "mac") {
		return "mac"
	}
	if strings.Contains(ua, "linux") {
		return "linux"
	}
	if strings.Contains(ua, "android") {
		return "android"
	}
	if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		return "ios"
	}
	return "unknown"
}

// HashString creates a SHA-256 hash of a string
func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

// GenerateCompositeKey creates a unique key from IP, UA, and day
func GenerateCompositeKey(ip, userAgent, day string) string {
	// Combine all factors with a separator
	combined := fmt.Sprintf("%s|%s|%s", ip, userAgent, day)
	return HashString(combined)
}

// GetNextMidnightUTC returns the next midnight in UTC timezone
func GetNextMidnightUTC() time.Time {
	now := time.Now().UTC()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return nextMidnight
}

// isValidIP checks if a string is a valid IP address
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// GetWindowResetString returns the window reset time as a string for DB queries
func GetWindowResetString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05-07")
}

// EnhancedFingerprint can be extended with additional factors
type EnhancedFingerprint struct {
	*GuestFingerprint
	AcceptLanguage string
	AcceptEncoding string
	ScreenInfo     string // Optional: from client-side JS
	CanvasHash     string // Optional: from client-side canvas fingerprinting
}

// GenerateEnhancedFingerprint creates a more robust fingerprint
// Use this if you need stronger protection against sophisticated bypasses
func GenerateEnhancedFingerprint(c *fiber.Ctx) *EnhancedFingerprint {
	base := GenerateGuestFingerprint(c)

	return &EnhancedFingerprint{
		GuestFingerprint: base,
		AcceptLanguage:   c.Get("Accept-Language"),
		AcceptEncoding:   c.Get("Accept-Encoding"),
		ScreenInfo:       c.Get("X-Screen-Info"),        // Custom header from client
		CanvasHash:       c.Get("X-Canvas-Fingerprint"), // Custom header from client
	}
}
