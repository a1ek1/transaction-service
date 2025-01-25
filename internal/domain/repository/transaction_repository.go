// Package repository defines interfaces for interacting with persistent storage.
package repository

import (
	"context"
	"github.com/google/uuid"
	"transaction-service/internal/domain/model"
)

// TransactionRepository defines methods for managing transactions in the database.
type TransactionRepository interface {
	// Create adds a new transaction to the database.
	Create(ctx context.Context, transaction *model.Transaction) (uuid.UUID, error)

	// GetTransactions retrieves a list of transactions from the database.
	GetTransactions(ctx context.Context) ([]model.Transaction, error)
}
