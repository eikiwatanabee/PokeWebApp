package uow

import "context"

// UnitOfWork manages transactional boundaries.
// Implemented in infrastructure layer with GORM.
type UnitOfWork interface {
	// Do executes the given function within a transaction.
	// If fn returns an error, the transaction is rolled back.
	// If fn returns nil, the transaction is committed.
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
