package repository

import (
	"context"
	"transaction-service/internal/domain/model"
)

type TransactionRepository interface {
	Create(ctx context.Context, wallet *model.Transaction) (int, error)
}
