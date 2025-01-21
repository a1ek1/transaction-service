package datastore

import (
	"context"
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
	tx, err := tr.BeginTransaction()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	var id uuid.UUID

	err = tx.QueryRowxContext(ctx,
		`INSERT INTO transactions (id, "from", "to", amount, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id`,
		transaction.ID, transaction.From, transaction.To, transaction.Amount,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
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
	if err := conn.SelectContext(ctx, &transactions, "SELECT * FROM transactions"); err != nil {
		return nil, err
	}

	return lo.Map(transactions, func(transaction dbTransaction, _ int) model.Transaction { return model.Transaction(transaction) }), nil
}

type dbTransaction struct {
	ID        uuid.UUID `db:"id"`
	From      string    `db:"from"`
	To        string    `db:"to"`
	Amount    int       `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}
