package service

import (
	"context"
	"transaction-service/internal/domain/model"
	"transaction-service/internal/domain/repository"
)

type TransactionService interface {
	GetNTransactions(ctx context.Context, n int) ([]model.Transaction, error)
}

type transactionService struct {
	repository.TransactionRepository
}

func (t transactionService) GetNTransactions(ctx context.Context, n int) ([]model.Transaction, error) {
	transactions, err := t.GetTransactions(ctx)
	if err != nil {
		return nil, err
	}
	return transactions[:n], nil
}

func NewTransactionService(repository repository.TransactionRepository) TransactionService {
	return &transactionService{repository}
}
