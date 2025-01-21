package repository

import (
	"context"
	"github.com/google/uuid"
	"transaction-service/internal/domain/model"
)

type TransactionRepository interface {
	Create(ctx context.Context, wallet *model.Transaction) (uuid.UUID, error)
	GetTransactions(ctx context.Context) ([]model.Transaction, error)
}
