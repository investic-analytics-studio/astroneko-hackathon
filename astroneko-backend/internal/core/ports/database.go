package ports

import (
	"context"
)

// DatabaseInterface defines the contract for database operations
type DatabaseInterface interface {
	// Transaction management
	WithContext(ctx context.Context) DatabaseInterface
	Begin() DatabaseInterface
	Commit() error
	Rollback() error

	// CRUD operations
	Create(value any) error
	Save(value any) error
	First(dest any, conds ...any) error
	Find(dest any, conds ...any) error
	Where(query any, args ...any) DatabaseInterface
	Model(value any) DatabaseInterface
	Update(column string, value any) error
	Updates(values any) error
	Delete(value any, conds ...any) error

	// Query building
	Select(query any, args ...any) DatabaseInterface
	Omit(columns ...string) DatabaseInterface
	Order(value any) DatabaseInterface
	Limit(limit int) DatabaseInterface
	Offset(offset int) DatabaseInterface
	Count(count *int64) error

	// Raw queries
	Exec(sql string, values ...any) error
	Raw(sql string, values ...any) DatabaseInterface
	Scan(dest any) error
}
