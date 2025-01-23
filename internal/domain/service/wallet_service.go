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
	if fromID == toID {
		return fmt.Errorf("cannot send money to the same wallet")
	}
	if amount <= 0 || amount > 10000000 { // Лимит: 10,000.00
		return fmt.Errorf("amount must be between 0 and 10,000")
	}

	tx, err := w.walletRepo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	senderWallet, err := w.walletRepo.FetchByID(ctx, fromID)
	if err != nil {
		return fmt.Errorf("failed to fetch sender wallet: %w", err)
	}
	if senderWallet.Amount < amount {
		return fmt.Errorf("insufficient funds")
	}

	receiverWallet, err := w.walletRepo.FetchByID(ctx, toID)
	if err != nil {
		return fmt.Errorf("failed to fetch receiver wallet: %w", err)
	}

	senderWallet.Amount -= amount
	receiverWallet.Amount += amount

	if _, err := w.walletRepo.Update(ctx, senderWallet); err != nil {
		return fmt.Errorf("failed to update sender wallet: %w", err)
	}
	if _, err := w.walletRepo.Update(ctx, receiverWallet); err != nil {
		return fmt.Errorf("failed to update receiver wallet: %w", err)
	}

	transaction := &model.Transaction{
		ID:        uuid.New(),
		From:      fromID.String(),
		To:        toID.String(),
		Amount:    amount,
		CreatedAt: time.Now(),
	}
	if _, err := w.transactionRepo.Create(ctx, transaction); err != nil {
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
