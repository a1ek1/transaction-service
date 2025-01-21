package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"transaction-service/internal/domain/service"
	"transaction-service/internal/infrastructure/datastore"
)

func main() {
	db, err := sqlx.Open("postgres", "###")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	repository := datastore.NewWalletRepositoryImpl(db)
	walletService := service.NewWalletService(repository)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wallet1ID, err := repository.Create(ctx)
	if err != nil {
		log.Fatalf("Failed to create wallet 1: %v", err)
	}
	wallet2ID, err := repository.Create(ctx)
	if err != nil {
		log.Fatalf("Failed to create wallet 2: %v", err)
	}
	fmt.Printf("Created wallets: %v (Wallet 1), %v (Wallet 2)\n", wallet1ID, wallet2ID)

	balance, err := walletService.GetBalance(ctx, wallet1ID)
	if err != nil {
		log.Fatalf("Failed to get balance of wallet 1: %v", err)
	}
	fmt.Printf("Balance of wallet 1: %d\n", balance)

	user1, err := repository.FetchByID(ctx, wallet1ID)
	if err != nil {
		log.Fatalf("Failed to fetch amount to transfer: %v", err)
	}

	amountToTransfer := user1.Amount

	fmt.Printf("Transferring %d units from wallet 1 to wallet 2...\n", amountToTransfer)
	err = walletService.SendMoney(ctx, wallet1ID, wallet2ID, amountToTransfer)
	if err != nil {
		log.Fatalf("Failed to transfer money: %v", err)
	}

	balance1, err := walletService.GetBalance(ctx, wallet1ID)
	if err != nil {
		log.Fatalf("Failed to get balance of wallet 1 after transfer: %v", err)
	}
	balance2, err := walletService.GetBalance(ctx, wallet2ID)
	if err != nil {
		log.Fatalf("Failed to get balance of wallet 2 after transfer: %v", err)
	}

	fmt.Printf("Balances after transfer: Wallet 1 = %d, Wallet 2 = %d\n", balance1, balance2)
}
