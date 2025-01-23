package datastore

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
	"transaction-service/internal/domain/model"
	"transaction-service/internal/domain/repository"
)

type transactionRepositoryImpl struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) repository.TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (tr *transactionRepositoryImpl) BeginTransaction() (*sqlx.Tx, error) {
	tx, err := tr.db.Beginx()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (tr *transactionRepositoryImpl) Create(ctx context.Context, transaction *model.Transaction) (uuid.UUID, error) {
	if transaction == nil {
		return uuid.Nil, fmt.Errorf("transaction cannot be nil")
	}

	tx, err := tr.BeginTransaction()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
        INSERT INTO transactions (id, "from", "to", amount, created_at) 
        VALUES (:id, :from, :to, :amount, NOW()) 
        RETURNING id
    `
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	var id uuid.UUID
	err = stmt.GetContext(ctx, &id, map[string]interface{}{
		"id":     transaction.ID,
		"from":   transaction.From,
		"to":     transaction.To,
		"amount": transaction.Amount,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (tr *transactionRepositoryImpl) GetTransactions(ctx context.Context) ([]model.Transaction, error) {
	conn, err := tr.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var transactions []dbTransaction
	query := `SELECT id, "from", "to", amount, created_at FROM transactions ORDER BY created_at DESC LIMIT 100`
	if err := conn.SelectContext(ctx, &transactions, query); err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	return lo.Map(transactions, func(transaction dbTransaction, _ int) model.Transaction {
		return model.Transaction(transaction)
	}), nil
}

type dbTransaction struct {
	ID        uuid.UUID `db:"id"`
	From      string    `db:"from"`
	To        string    `db:"to"`
	Amount    int       `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}
