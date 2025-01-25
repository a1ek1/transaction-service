// Package model defines the core data models used in the transaction service.
package model

import (
	"github.com/google/uuid"
	"time"
)

// Transaction represents a financial transaction between two wallets.
type Transaction struct {
	ID        uuid.UUID // Unique identifier for the transaction
	From      string    // Wallet ID of the sender
	To        string    // Wallet ID of the receiver
	Amount    int       // Transaction amount in cents
	CreatedAt time.Time // Timestamp of when the transaction was created
}
