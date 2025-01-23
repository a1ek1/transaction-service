package usecase

import (
	"context"
	"fmt"
	"transaction-service/internal/domain/service"
)

type TransactionUsecase interface {
	GetLastTransactions(ctx context.Context, count int) ([]TransactionDTO, error)
}

type transactionUsecase struct {
	transactionService service.TransactionService
}

func (u *transactionUsecase) GetLastTransactions(ctx context.Context, count int) ([]TransactionDTO, error) {
	transactions, err := u.transactionService.GetNTransactions(ctx, count)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	transactionDTOs := make([]TransactionDTO, len(transactions))
	for i, t := range transactions {
		transactionDTOs[i] = TransactionDTO{
			ID:        t.ID.String(),
			From:      t.From,
			To:        t.To,
			Amount:    float64(t.Amount) / 100,
			CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return transactionDTOs, nil
}

func NewTransactionUsecase(transactionService service.TransactionService) TransactionUsecase {
	return &transactionUsecase{
		transactionService: transactionService,
	}
}

type TransactionDTO struct {
	ID        string  `json:"id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
}
