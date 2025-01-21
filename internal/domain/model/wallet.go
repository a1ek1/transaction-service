package model

import "github.com/google/uuid"

type Wallet struct {
	ID     uuid.UUID
	Amount int
}

func NewWallet() *Wallet {
	return &Wallet{
		ID:     uuid.New(),
		Amount: 0,
	}
}
