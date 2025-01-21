package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"transaction-service/internal/domain/model"
)

type WalletRepository interface {
	FetchByID(ctx context.Context, id uuid.UUID) (*model.Wallet, error)
	Create(ctx context.Context) (uuid.UUID, error)
	Update(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error)
	Delete(ctx context.Context, id uuid.UUID) error
	BeginTransaction() (*sqlx.Tx, error)
}
