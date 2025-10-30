package guest_usage

import "time"

// GuestAPIUsage represents a guest user's API usage record
type GuestAPIUsage struct {
	ID            string
	IPAddress     string
	UserAgentHash string
	CompositeKey  string
	Endpoint      string
	UsageCount    int
	DailyLimit    int
	WindowResetAt time.Time
	LastRequestAt time.Time
	IsBlocked     bool
	BlockedReason *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CanMakeRequest checks if guest can make another request
func (g *GuestAPIUsage) CanMakeRequest() bool {
	if g.IsBlocked {
		return false
	}

	// Check if window has expired (reset daily limit)
	if time.Now().After(g.WindowResetAt) {
		return true
	}

	return g.UsageCount < g.DailyLimit
}

// IncrementUsage increments the usage count
func (g *GuestAPIUsage) IncrementUsage() {
	g.UsageCount++
	g.LastRequestAt = time.Now()
}

// RemainingRequests returns how many requests are left
func (g *GuestAPIUsage) RemainingRequests() int {
	if time.Now().After(g.WindowResetAt) {
		return g.DailyLimit
	}

	remaining := g.DailyLimit - g.UsageCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ResetWindow resets the usage window to next day
func (g *GuestAPIUsage) ResetWindow() {
	g.UsageCount = 0
	g.WindowResetAt = GetNextResetTime()
}

// GetNextResetTime returns the next midnight UTC
func GetNextResetTime() time.Time {
	now := time.Now().UTC()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return nextMidnight
}

// CreateGuestUsageRequest represents request to create new guest usage
type CreateGuestUsageRequest struct {
	IPAddress  string
	UserAgent  string
	Endpoint   string
	DailyLimit int
}
