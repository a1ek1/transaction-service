package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"transaction-service/internal/domain/repository"
)

type WalletService interface {
	SendMoney(ctx context.Context, fromID, toID uuid.UUID, amount int) error
	GetBalance(ctx context.Context, id uuid.UUID) (amount int, err error)
}

type walletService struct {
	repository.WalletRepository
}

func NewWalletService(repository repository.WalletRepository) WalletService {
	return &walletService{repository}
}

func (w walletService) SendMoney(ctx context.Context, fromID, toID uuid.UUID, amount int) error {
	tx, err := w.WalletRepository.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	senderWallet, err := w.WalletRepository.FetchByID(ctx, fromID)
	if err != nil {
		return err
	}

	if senderWallet.Amount < amount {
		return fmt.Errorf("insufficient funds")
	}

	senderWallet.Amount -= amount
	_, err = w.WalletRepository.Update(ctx, senderWallet)
	if err != nil {
		return err
	}

	receiverWallet, err := w.WalletRepository.FetchByID(ctx, toID)
	if err != nil {
		return err
	}

	receiverWallet.Amount += amount
	_, err = w.WalletRepository.Update(ctx, receiverWallet)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (w walletService) GetBalance(ctx context.Context, id uuid.UUID) (amount int, err error) {
	wallet, err := w.WalletRepository.FetchByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch wallet: %w", err)
	}
	return wallet.Amount, nil
}
