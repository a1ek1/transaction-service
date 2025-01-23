package interactor

import (
	"github.com/jmoiron/sqlx"
	"transaction-service/internal/domain/repository"
	"transaction-service/internal/domain/service"
	"transaction-service/internal/infrastructure/datastore"
	"transaction-service/internal/presenter/http/handler"
	"transaction-service/internal/usecase"
)

type Interactor interface {
	NewWalletRepository() repository.WalletRepository
	NewTransactionRepository() repository.TransactionRepository
	NewWalletService() service.WalletService
	NewTransactionService() service.TransactionService
	NewWalletUsecase() usecase.WalletUsecase
	NewTransactionUsecase() usecase.TransactionUsecase
	NewWalletHandler() handler.WalletHandler
	NewTransactionHandler() handler.TransactionHandler
	NewAppHandler() handler.AppHandler
}

type interactor struct {
	DB *sqlx.DB
}

func NewInteractor(db *sqlx.DB) Interactor {
	return &interactor{DB: db}
}

type appHandler struct {
	handler.WalletHandler
	handler.TransactionHandler
}

func (i *interactor) NewAppHandler() handler.AppHandler {
	return &appHandler{
		WalletHandler:      i.NewWalletHandler(),
		TransactionHandler: i.NewTransactionHandler(),
	}
}

func (i *interactor) NewWalletRepository() repository.WalletRepository {
	return datastore.NewWalletRepositoryImpl(i.DB)
}

func (i *interactor) NewTransactionRepository() repository.TransactionRepository {
	return datastore.NewTransactionRepository(i.DB)
}

func (i *interactor) NewWalletService() service.WalletService {
	return service.NewWalletService(i.NewWalletRepository(), i.NewTransactionRepository())
}

func (i *interactor) NewTransactionService() service.TransactionService {
	return service.NewTransactionService(i.NewTransactionRepository())
}

func (i *interactor) NewWalletUsecase() usecase.WalletUsecase {
	return usecase.NewWalletUsecase(i.NewWalletService())
}

func (i *interactor) NewTransactionUsecase() usecase.TransactionUsecase {
	return usecase.NewTransactionUsecase(i.NewTransactionService())
}

func (i *interactor) NewWalletHandler() handler.WalletHandler {
	return handler.NewWalletHandler(i.NewWalletUsecase())
}

func (i *interactor) NewTransactionHandler() handler.TransactionHandler {
	return handler.NewTransactionHandler(i.NewTransactionUsecase())
}
