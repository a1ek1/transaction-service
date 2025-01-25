// Package repository defines interfaces for interacting with persistent storage.
package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"transaction-service/internal/domain/model"
)

// WalletRepository defines methods for managing wallets in the database.
type WalletRepository interface {
	// FetchByID retrieves a wallet by its unique ID.
	FetchByID(ctx context.Context, id uuid.UUID) (*model.Wallet, error)

	// Create adds a new wallet to the database and returns its ID.
	Create(ctx context.Context) (uuid.UUID, error)

	// Update modifies the details of an existing wallet.
	Update(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error)

	// Delete removes a wallet from the database by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// BeginTransaction starts a new database transaction.
	BeginTransaction() (*sqlx.Tx, error)

	// IsServiceInitialized shows if there are 10 records in the database.
	IsServiceInitialized(ctx context.Context) (bool, error)

	// SetServiceInitialized adds 10 records to the database
	SetServiceInitialized(ctx context.Context) error

	// FetchAll returns all records from the database
	FetchAll(ctx context.Context) ([]*model.Wallet, error)
}
