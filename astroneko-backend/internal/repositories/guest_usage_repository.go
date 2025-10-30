package repositories

import (
	"context"
	"database/sql"
	"time"

	"astroneko-backend/internal/core/domain/guest_usage"
	"astroneko-backend/internal/core/ports"
	"astroneko-backend/pkg/logger"
)

type GuestUsageRepository struct {
	db     ports.DatabaseInterface
	logger logger.Logger
}

func NewGuestUsageRepository(db ports.DatabaseInterface, log logger.Logger) *GuestUsageRepository {
	return &GuestUsageRepository{
		db:     db,
		logger: log,
	}
}

type guestUsageModel struct {
	ID            string         `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	IPAddress     string         `gorm:"column:ip_address;type:inet"`
	UserAgentHash string         `gorm:"column:user_agent_hash"`
	CompositeKey  string         `gorm:"column:composite_key"`
	Endpoint      string         `gorm:"column:endpoint"`
	UsageCount    int            `gorm:"column:usage_count;default:1"`
	DailyLimit    int            `gorm:"column:daily_limit;default:3"`
	WindowResetAt time.Time      `gorm:"column:window_reset_at;type:timestamptz"`
	LastRequestAt time.Time      `gorm:"column:last_request_at;type:timestamptz"`
	IsBlocked     bool           `gorm:"column:is_blocked;default:false"`
	BlockedReason sql.NullString `gorm:"column:blocked_reason"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`
}

func (guestUsageModel) TableName() string {
	return "astroneko_guest_api_usage"
}

func (m *guestUsageModel) toDomain() *guest_usage.GuestAPIUsage {
	var blockedReason *string
	if m.BlockedReason.Valid {
		blockedReason = &m.BlockedReason.String
	}

	return &guest_usage.GuestAPIUsage{
		ID:            m.ID,
		IPAddress:     m.IPAddress,
		UserAgentHash: m.UserAgentHash,
		CompositeKey:  m.CompositeKey,
		Endpoint:      m.Endpoint,
		UsageCount:    m.UsageCount,
		DailyLimit:    m.DailyLimit,
		WindowResetAt: m.WindowResetAt,
		LastRequestAt: m.LastRequestAt,
		IsBlocked:     m.IsBlocked,
		BlockedReason: blockedReason,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// GetByCompositeKey retrieves guest usage by composite key and endpoint
func (r *GuestUsageRepository) GetByCompositeKey(ctx context.Context, compositeKey, endpoint string, windowResetAt string) (*guest_usage.GuestAPIUsage, error) {
	var model guestUsageModel

	err := r.db.WithContext(ctx).
		Where("composite_key = ? AND endpoint = ? AND window_reset_at = ?", compositeKey, endpoint, windowResetAt).
		First(&model)

	if err != nil {
		// GORM returns ErrRecordNotFound when no record is found
		if err.Error() == "record not found" {
			return nil, nil // Not found, not an error for our use case
		}
		r.logger.Error("Failed to get guest usage by composite key",
			logger.Field{Key: "composite_key", Value: compositeKey},
			logger.Field{Key: "endpoint", Value: endpoint},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	return model.toDomain(), nil
}

// Create creates a new guest usage record
func (r *GuestUsageRepository) Create(ctx context.Context, usage *guest_usage.GuestAPIUsage) error {
	model := &guestUsageModel{
		// Don't set ID - let DB generate UUID via gen_random_uuid()
		IPAddress:     usage.IPAddress,
		UserAgentHash: usage.UserAgentHash,
		CompositeKey:  usage.CompositeKey,
		Endpoint:      usage.Endpoint,
		UsageCount:    usage.UsageCount,
		DailyLimit:    usage.DailyLimit,
		WindowResetAt: usage.WindowResetAt,
		LastRequestAt: usage.LastRequestAt,
		IsBlocked:     usage.IsBlocked,
	}

	if usage.BlockedReason != nil {
		model.BlockedReason = sql.NullString{String: *usage.BlockedReason, Valid: true}
	}

	// Use Omit to skip ID field, letting DB generate it
	err := r.db.WithContext(ctx).Omit("id").Create(model)
	if err != nil {
		r.logger.Error("Failed to create guest usage record",
			logger.Field{Key: "composite_key", Value: usage.CompositeKey},
			logger.Field{Key: "endpoint", Value: usage.Endpoint},
			logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	// Update the domain object with generated values
	usage.ID = model.ID
	usage.CreatedAt = model.CreatedAt
	usage.UpdatedAt = model.UpdatedAt

	return nil
}

// IncrementUsage increments usage count and updates last request time
func (r *GuestUsageRepository) IncrementUsage(ctx context.Context, id string) error {
	// Use raw SQL to increment usage_count atomically
	err := r.db.WithContext(ctx).Exec(
		"UPDATE astroneko_guest_api_usage SET usage_count = usage_count + 1, last_request_at = NOW() WHERE id = ?",
		id,
	)

	if err != nil {
		r.logger.Error("Failed to increment guest usage",
			logger.Field{Key: "id", Value: id},
			logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	return nil
}

// ResetExpiredWindows resets usage for expired time windows
func (r *GuestUsageRepository) ResetExpiredWindows(ctx context.Context) error {
	// Use raw SQL to update with interval arithmetic
	err := r.db.WithContext(ctx).Exec(
		"UPDATE astroneko_guest_api_usage SET usage_count = 0, window_reset_at = window_reset_at + INTERVAL '1 day' WHERE window_reset_at < NOW() AND is_blocked = false",
	)

	if err != nil {
		r.logger.Error("Failed to reset expired windows",
			logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	r.logger.Info("Reset expired guest usage windows completed")

	return nil
}

// GetByIPAddress retrieves all usage records for an IP (abuse detection)
func (r *GuestUsageRepository) GetByIPAddress(ctx context.Context, ipAddress string, since string) ([]*guest_usage.GuestAPIUsage, error) {
	var models []guestUsageModel

	err := r.db.WithContext(ctx).
		Where("ip_address = ? AND created_at >= ?", ipAddress, since).
		Order("created_at DESC").
		Find(&models)

	if err != nil {
		r.logger.Error("Failed to get usage by IP",
			logger.Field{Key: "ip_address", Value: ipAddress},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	usages := make([]*guest_usage.GuestAPIUsage, len(models))
	for i, model := range models {
		usages[i] = model.toDomain()
	}

	return usages, nil
}

// BlockGuest blocks a guest from making requests
func (r *GuestUsageRepository) BlockGuest(ctx context.Context, compositeKey string, reason string) error {
	err := r.db.WithContext(ctx).
		Model(&guestUsageModel{}).
		Where("composite_key = ?", compositeKey).
		Updates(map[string]interface{}{
			"is_blocked":     true,
			"blocked_reason": reason,
		})

	if err != nil {
		r.logger.Error("Failed to block guest",
			logger.Field{Key: "composite_key", Value: compositeKey},
			logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	r.logger.Warn("Blocked guest user",
		logger.Field{Key: "composite_key", Value: compositeKey},
		logger.Field{Key: "reason", Value: reason})

	return nil
}

// DeleteOldRecords deletes records older than retention period (cleanup)
func (r *GuestUsageRepository) DeleteOldRecords(ctx context.Context, olderThan string) error {
	err := r.db.WithContext(ctx).
		Where("created_at < ? AND is_blocked = ?", olderThan, false).
		Delete(&guestUsageModel{})

	if err != nil {
		r.logger.Error("Failed to delete old guest usage records",
			logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	r.logger.Info("Deleted old guest usage records successfully")

	return nil
}
