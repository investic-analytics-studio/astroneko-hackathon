package adapters

import (
	"context"

	"astroneko-backend/internal/core/ports"
	"gorm.io/gorm"
)

// GormAdapter wraps gorm.DB to implement DatabaseInterface
type GormAdapter struct {
	db *gorm.DB
}

// NewGormAdapter creates a new GORM adapter
func NewGormAdapter(db *gorm.DB) ports.DatabaseInterface {
	return &GormAdapter{db: db}
}

// WithContext returns a new adapter with the given context
func (g *GormAdapter) WithContext(ctx context.Context) ports.DatabaseInterface {
	return &GormAdapter{db: g.db.WithContext(ctx)}
}

// Begin starts a transaction
func (g *GormAdapter) Begin() ports.DatabaseInterface {
	return &GormAdapter{db: g.db.Begin()}
}

// Commit commits the transaction
func (g *GormAdapter) Commit() error {
	return g.db.Commit().Error
}

// Rollback rolls back the transaction
func (g *GormAdapter) Rollback() error {
	return g.db.Rollback().Error
}

// Create inserts a new record
func (g *GormAdapter) Create(value any) error {
	return g.db.Create(value).Error
}

// Save updates all fields of a record
func (g *GormAdapter) Save(value any) error {
	return g.db.Save(value).Error
}

// First finds the first record matching given conditions
func (g *GormAdapter) First(dest any, conds ...any) error {
	return g.db.First(dest, conds...).Error
}

// Find finds all records matching given conditions
func (g *GormAdapter) Find(dest any, conds ...any) error {
	return g.db.Find(dest, conds...).Error
}

// Where adds WHERE condition
func (g *GormAdapter) Where(query any, args ...any) ports.DatabaseInterface {
	return &GormAdapter{db: g.db.Where(query, args...)}
}

// Model specifies the model to operate on
func (g *GormAdapter) Model(value any) ports.DatabaseInterface {
	return &GormAdapter{db: g.db.Model(value)}
}

// Update updates a single column
func (g *GormAdapter) Update(column string, value any) error {
	return g.db.Update(column, value).Error
}

// Updates updates multiple columns
func (g *GormAdapter) Updates(values any) error {
	return g.db.Updates(values).Error
}

// Delete deletes records matching given conditions
func (g *GormAdapter) Delete(value any, conds ...any) error {
	return g.db.Delete(value, conds...).Error
}

// Select specifies fields to retrieve
func (g *GormAdapter) Select(query any, args ...any) ports.DatabaseInterface {
	return &GormAdapter{db: g.db.Select(query, args...)}
}

// Omit specifies fields to ignore when creating/updating
func (g *GormAdapter) Omit(columns ...string) ports.DatabaseInterface {
	return &GormAdapter{db: g.db.Omit(columns...)}
}

// Order specifies order when retrieving records
func (g *GormAdapter) Order(value any) ports.DatabaseInterface {
	return &GormAdapter{db: g.db.Order(value)}
}

// Limit specifies number of records to retrieve
func (g *GormAdapter) Limit(limit int) ports.DatabaseInterface {
	return &GormAdapter{db: g.db.Limit(limit)}
}

// Offset specifies number of records to skip before starting to return records
func (g *GormAdapter) Offset(offset int) ports.DatabaseInterface {
	return &GormAdapter{db: g.db.Offset(offset)}
}

// Count gets the count of records
func (g *GormAdapter) Count(count *int64) error {
	return g.db.Count(count).Error
}

// Exec executes raw SQL
func (g *GormAdapter) Exec(sql string, values ...any) error {
	return g.db.Exec(sql, values...).Error
}

// Raw creates a raw query
func (g *GormAdapter) Raw(sql string, values ...any) ports.DatabaseInterface {
	return &GormAdapter{db: g.db.Raw(sql, values...)}
}

// Scan scans result into destination
func (g *GormAdapter) Scan(dest any) error {
	return g.db.Scan(dest).Error
}
