package service

import (
	"context"
	"github.com/google/uuid"

	"transaction-service/internal/domain/repository"
)

type WalletService interface {
	SendMoney(ctx context.Context, amount int) error
	GetBalance(ctx context.Context, id uuid.UUID) (amount int, err error)
}

type walletService struct {
	repository.WalletRepository
}

func NewWalletService(repository repository.WalletRepository) WalletService {
	return &walletService{repository}
}

func (w walletService) SendMoney(ctx context.Context, amount int) error {
	//TODO implement me
	panic("implement me")
}

func (w walletService) GetBalance(ctx context.Context, id uuid.UUID) (amount int, err error) {
	//TODO implement me
	panic("implement me")
}
