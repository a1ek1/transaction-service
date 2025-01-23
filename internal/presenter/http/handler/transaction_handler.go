package handler

import (
	"net/http"
	"strconv"
	"transaction-service/internal/usecase"

	"github.com/labstack/echo"
)

// TransactionHandler интерфейс для обработки запросов транзакций
type TransactionHandler interface {
	GetLastTransactions(c echo.Context) error
}

// Реализация TransactionHandler
type transactionHandlerImpl struct {
	TransactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) TransactionHandler {
	return &transactionHandlerImpl{TransactionUsecase: transactionUsecase}
}

func (h *transactionHandlerImpl) GetLastTransactions(c echo.Context) error {
	// Чтение параметра count из запроса
	countParam := c.QueryParam("count")
	if countParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "count parameter is required"})
	}

	count, err := strconv.Atoi(countParam)
	if err != nil || count <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid count parameter"})
	}

	transactions, err := h.TransactionUsecase.GetLastTransactions(c.Request().Context(), count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, transactions)
}
