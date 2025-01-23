package service

import (
	"context"
	"sync"
	"transaction-service/internal/domain/model"
	"transaction-service/internal/domain/repository"
)

type TransactionService interface {
	GetNTransactions(ctx context.Context, n int) ([]model.Transaction, error)
}

type transactionService struct {
	repository     repository.TransactionRepository
	workerPoolSize int
}

func (t *transactionService) GetNTransactions(ctx context.Context, n int) ([]model.Transaction, error) {
	transactions, err := t.repository.GetTransactions(ctx)
	if err != nil {
		return nil, err
	}

	// Ограничиваем количество рабочих потоков
	workerCount := min(n, t.workerPoolSize)
	results := make([]model.Transaction, 0, n)
	resultsMu := sync.Mutex{} // Мьютекс для защиты результатов
	resultCh := make(chan model.Transaction, n)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := start; j < len(transactions); j += workerCount {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					resultsMu.Lock()
					if len(results) < n {
						results = append(results, transactions[j])
						resultsMu.Unlock()
					} else {
						resultsMu.Unlock()
						return
					}
				}
			}
		}(i)
	}

	// Ждем завершения всех горутин
	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	// Проверяем наличие ошибок
	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return results[:min(len(results), n)], nil
}

func NewTransactionService(repository repository.TransactionRepository) TransactionService {
	return &transactionService{repository: repository, workerPoolSize: 5}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
