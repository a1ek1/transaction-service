// Package usecase implements application-specific logic for wallets.
package usecase

import (
	"context"
	"fmt"
	"transaction-service/internal/domain/model"
	"transaction-service/internal/domain/service"

	"github.com/google/uuid"
)

// WalletUsecase defines application-level logic for wallets.
type WalletUsecase interface {
	// SendMoney transfers funds between wallets.
	SendMoney(ctx context.Context, fromID, toID string, amount float64) error

	// GetBalance retrieves the balance of a wallet by its string ID.
	GetBalance(ctx context.Context, walletID string) (float64, error)

	GetAllWallets(ctx context.Context) ([]*model.Wallet, error)
}

type walletUsecase struct {
	walletService service.WalletService
}

func NewWalletUsecase(walletService service.WalletService) WalletUsecase {
	return &walletUsecase{
		walletService: walletService,
	}
}

func (u *walletUsecase) SendMoney(ctx context.Context, fromID, toID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	fromUUID, err := uuid.Parse(fromID)
	if err != nil {
		return fmt.Errorf("invalid 'from' wallet ID: %w", err)
	}

	toUUID, err := uuid.Parse(toID)
	if err != nil {
		return fmt.Errorf("invalid 'to' wallet ID: %w", err)
	}

	amountInCents := int(amount * 100)
	if err := u.walletService.SendMoney(ctx, fromUUID, toUUID, amountInCents); err != nil {
		return fmt.Errorf("failed to send money: %w", err)
	}
	return nil
}

func (u *walletUsecase) GetBalance(ctx context.Context, walletID string) (float64, error) {
	walletUUID, err := uuid.Parse(walletID)
	if err != nil {
		return 0, fmt.Errorf("invalid wallet ID: %w", err)
	}

	amount, err := u.walletService.GetBalance(ctx, walletUUID)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return float64(amount) / 100, nil
}

func (u *walletUsecase) GetAllWallets(ctx context.Context) ([]*model.Wallet, error) {
	wallets, err := u.walletService.FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wallets: %w", err)
	}
	return wallets, nil
}
