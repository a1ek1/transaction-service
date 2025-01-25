// Package handler implements HTTP handlers for transaction-related operations.
package handler

import (
	"net/http"
	"strconv"
	"transaction-service/internal/usecase"

	"github.com/labstack/echo"
)

// TransactionHandler defines HTTP endpoints for transaction operations.
type TransactionHandler interface {
	// GetLastTransactions handles the request to retrieve recent transactions.
	GetLastTransactions(c echo.Context) error
}

type transactionHandlerImpl struct {
	TransactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) TransactionHandler {
	return &transactionHandlerImpl{TransactionUsecase: transactionUsecase}
}

func (h *transactionHandlerImpl) GetLastTransactions(c echo.Context) error {
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
