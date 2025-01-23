package usecase

import (
	"context"
	"fmt"
	"transaction-service/internal/domain/service"

	"github.com/google/uuid"
)

type WalletUsecase interface {
	SendMoney(ctx context.Context, fromID, toID string, amount float64) error
	GetBalance(ctx context.Context, walletID string) (float64, error)
}

type walletUsecase struct {
	walletService service.WalletService
}

func NewWalletUsecase(walletService service.WalletService) WalletUsecase {
	return &walletUsecase{
		walletService: walletService,
	}
}

// SendMoney обрабатывает перевод средств между кошельками
func (u *walletUsecase) SendMoney(ctx context.Context, fromID, toID string, amount float64) error {
	fromUUID, err := uuid.Parse(fromID)
	if err != nil {
		return fmt.Errorf("invalid from wallet ID: %w", err)
	}

	toUUID, err := uuid.Parse(toID)
	if err != nil {
		return fmt.Errorf("invalid to wallet ID: %w", err)
	}

	amountInCents := int(amount * 100)
	if err := u.walletService.SendMoney(ctx, fromUUID, toUUID, amountInCents); err != nil {
		return fmt.Errorf("failed to send money: %w", err)
	}
	return nil
}

// GetBalance возвращает баланс кошелька
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
