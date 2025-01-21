package datastore

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"transaction-service/internal/domain/model"
	"transaction-service/internal/domain/repository"
)

type walletRepositoryImpl struct {
	db *sqlx.DB
}

func NewWalletRepositoryImpl(db *sqlx.DB) repository.WalletRepository {
	return &walletRepositoryImpl{db: db}
}

func (w *walletRepositoryImpl) BeginTransaction() (*sqlx.Tx, error) {
	tx, err := w.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

func (w *walletRepositoryImpl) Create(ctx context.Context) (uuid.UUID, error) {
	tx, err := w.BeginTransaction()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	var id uuid.UUID
	err = tx.QueryRowxContext(
		ctx,
		"INSERT INTO wallets (id, amount) VALUES ($1, $2) RETURNING id",
		uuid.New(),
		rand.Intn(2147483647),
	).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (w *walletRepositoryImpl) FetchByID(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	conn, err := w.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var wallet dbWallet
	if err := conn.GetContext(ctx, &wallet, `SELECT * FROM wallets WHERE id = $1`, id); err != nil {
		return nil, err
	}

	return (*model.Wallet)(&wallet), nil
}

func (w *walletRepositoryImpl) Update(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	tx, err := w.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`UPDATE wallets SET amount = $1 WHERE id = $2`,
		wallet.Amount,
		wallet.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (w *walletRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := w.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM wallets WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

type dbWallet struct {
	ID     uuid.UUID `db:"id"`
	Amount int       `db:"amount"`
}
