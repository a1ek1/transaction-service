package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
	"transaction-service/internal/domain/model"
	"transaction-service/internal/domain/repository"
)

// WalletService defines methods for wallet-related operations.
type WalletService interface {
	// SendMoney transfers funds between two wallets.
	SendMoney(ctx context.Context, fromID, toID uuid.UUID, amount int) error

	// GetBalance retrieves the balance of a wallet by its ID.
	GetBalance(ctx context.Context, id uuid.UUID) (amount int, err error)

	// InitializeWallets create 10 wallets for first launch
	InitializeWallets(ctx context.Context) error
	FetchAll(ctx context.Context) ([]*model.Wallet, error)
}

type walletService struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository

	lockMap sync.Map
}

// FetchAll returns all records from the database
func (w *walletService) FetchAll(ctx context.Context) ([]*model.Wallet, error) {
	wallets, err := w.walletRepo.FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wallets: %w", err)
	}
	return wallets, nil
}

// NewWalletService creates a new instance of WalletService.
func NewWalletService(walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository) WalletService {
	return &walletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

func (w *walletService) InitializeWallets(ctx context.Context) error {
	initialized, err := w.walletRepo.IsServiceInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check initialization state: %w", err)
	}

	if initialized {
		return nil // Кошельки уже были созданы
	}

	tx, err := w.walletRepo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i := 0; i < 10; i++ {
		if _, err := w.walletRepo.Create(ctx); err != nil {
			return fmt.Errorf("failed to create wallet #%d: %w", i+1, err)
		}
	}

	if err := w.walletRepo.SetServiceInitialized(ctx); err != nil {
		return fmt.Errorf("failed to set service initialized: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit initialization: %w", err)
	}

	return nil
}

func (w *walletService) SendMoney(ctx context.Context, fromID, toID uuid.UUID, amount int) error {
	if fromID == toID {
		return fmt.Errorf("cannot send money to the same wallet")
	}
	if amount <= 0 || amount > 10000000 {
		return fmt.Errorf("amount must be between 0 and 10,000")
	}

	fromLock := w.getLock(fromID)
	toLock := w.getLock(toID)

	fromLock.Lock()
	defer fromLock.Unlock()

	toLock.Lock()
	defer toLock.Unlock()

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

func (w *walletService) GetBalance(ctx context.Context, id uuid.UUID) (amount int, err error) {
	wallet, err := w.walletRepo.FetchByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch wallet: %w", err)
	}
	return wallet.Amount, nil
}

// getLock возвращает блокировку для кошелька
func (w *walletService) getLock(walletID uuid.UUID) *sync.Mutex {
	lock, _ := w.lockMap.LoadOrStore(walletID, &sync.Mutex{})
	return lock.(*sync.Mutex)
}
