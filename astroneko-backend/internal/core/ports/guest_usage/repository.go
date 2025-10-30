package guest_usage

import (
	"context"

	"astroneko-backend/internal/core/domain/guest_usage"
)

// Repository defines the interface for guest API usage operations
type Repository interface {
	// GetByCompositeKey retrieves guest usage by composite fingerprint key
	GetByCompositeKey(ctx context.Context, compositeKey, endpoint string, windowResetAt string) (*guest_usage.GuestAPIUsage, error)

	// Create creates a new guest usage record
	Create(ctx context.Context, usage *guest_usage.GuestAPIUsage) error

	// IncrementUsage increments the usage count for existing record
	IncrementUsage(ctx context.Context, id string) error

	// ResetExpiredWindows resets usage counts for expired time windows
	ResetExpiredWindows(ctx context.Context) error

	// GetByIPAddress gets all usage records for an IP (for abuse detection)
	GetByIPAddress(ctx context.Context, ipAddress string, since string) ([]*guest_usage.GuestAPIUsage, error)

	// BlockGuest blocks a guest from making requests
	BlockGuest(ctx context.Context, compositeKey string, reason string) error

	// DeleteOldRecords cleanup records older than retention period
	DeleteOldRecords(ctx context.Context, olderThan string) error
}
