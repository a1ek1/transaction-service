package model

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID        uuid.UUID
	From      string
	To        string
	Amount    int
	CreatedAt time.Time
}
