package repositories

import (
	"context"

	"astroneko-backend/internal/core/ports"
	"astroneko-backend/pkg/utils"

	"github.com/google/uuid"
)

// BaseModel defines common fields for all models
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt string    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt string    `gorm:"autoUpdateTime" json:"updated_at"`
}

// GenericRepository provides common CRUD operations
type GenericRepository[T any] struct {
	db ports.DatabaseInterface
}

// NewGenericRepository creates a new generic repository
func NewGenericRepository[T any](db ports.DatabaseInterface) *GenericRepository[T] {
	return &GenericRepository[T]{
		db: db,
	}
}

// Create creates a new record with UUID
func (r *GenericRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	if err := r.db.WithContext(ctx).Create(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// GetByID retrieves a record by ID with UUID validation
func (r *GenericRepository[T]) GetByID(ctx context.Context, id string) (*T, error) {
	parsedID, err := utils.ValidateAndParseUUID(id)
	if err != nil {
		return nil, err
	}

	var entity T
	if err := r.db.WithContext(ctx).Where("id = ?", parsedID).First(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update updates a record
func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity)
}

// Delete deletes a record by ID
func (r *GenericRepository[T]) Delete(ctx context.Context, id string) error {
	parsedID, err := utils.ValidateAndParseUUID(id)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Delete(new(T), "id = ?", parsedID)
}

// List retrieves multiple records with pagination
func (r *GenericRepository[T]) List(ctx context.Context, limit, offset int) ([]*T, error) {
	var entities []*T
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&entities); err != nil {
		return nil, err
	}
	return entities, nil
}

// Count returns the total count of records
func (r *GenericRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&count); err != nil {
		return 0, err
	}
	return count, nil
}
