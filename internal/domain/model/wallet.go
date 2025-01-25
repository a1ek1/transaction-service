// Package model defines the core data models used in the transaction service.
package model

import "github.com/google/uuid"

// Wallet represents a digital wallet with a unique ID and a balance.
type Wallet struct {
	ID     uuid.UUID // Unique identifier for the wallet
	Amount int       // Current balance in the wallet, in cents
}
