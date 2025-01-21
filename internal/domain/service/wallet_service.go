package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
	"transaction-service/internal/domain/model"

	"transaction-service/internal/domain/repository"
)

type WalletService interface {
	SendMoney(ctx context.Context, fromID, toID uuid.UUID, amount int) error
	GetBalance(ctx context.Context, id uuid.UUID) (amount int, err error)
}

type walletService struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
}

func NewWalletService(walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository) WalletService {
	return &walletService{walletRepo: walletRepo, transactionRepo: transactionRepo}
}

func (w walletService) SendMoney(ctx context.Context, fromID, toID uuid.UUID, amount int) error {
	tx, err := w.walletRepo.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	senderWallet, err := w.walletRepo.FetchByID(ctx, fromID)
	if err != nil {
		return err
	}

	if senderWallet.Amount < amount {
		return fmt.Errorf("insufficient funds")
	}

	senderWallet.Amount -= amount
	_, err = w.walletRepo.Update(ctx, senderWallet)
	if err != nil {
		return err
	}

	receiverWallet, err := w.walletRepo.FetchByID(ctx, toID)
	if err != nil {
		return err
	}

	receiverWallet.Amount += amount
	_, err = w.walletRepo.Update(ctx, receiverWallet)
	if err != nil {
		return err
	}

	transaction := &model.Transaction{
		ID:        uuid.New(),
		From:      fromID.String(),
		To:        toID.String(),
		Amount:    amount,
		CreatedAt: time.Now(),
	}

	_, err = w.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (w walletService) GetBalance(ctx context.Context, id uuid.UUID) (amount int, err error) {
	wallet, err := w.walletRepo.FetchByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch wallet: %w", err)
	}
	return wallet.Amount, nil
}
