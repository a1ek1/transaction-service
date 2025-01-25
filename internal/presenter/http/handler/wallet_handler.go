// Package handler implements HTTP handlers for wallet-related operations.
package handler

import (
	"github.com/google/uuid"
	"net/http"
	"transaction-service/internal/usecase"

	"github.com/labstack/echo"
)

// WalletHandler defines HTTP endpoints for wallet operations.
type WalletHandler interface {
	// GetBalance handles the request to retrieve a wallet's balance.
	GetBalance(c echo.Context) error

	// SendMoney handles the request to transfer money between wallets.
	SendMoney(c echo.Context) error

	// GetAllWallets handles the request to return all of the wallets
	GetAllWallets(c echo.Context) error
}

type walletHandlerImpl struct {
	WalletUsecase usecase.WalletUsecase
}

func NewWalletHandler(walletUsecase usecase.WalletUsecase) WalletHandler {
	return &walletHandlerImpl{WalletUsecase: walletUsecase}
}

func (h *walletHandlerImpl) GetBalance(c echo.Context) error {
	address := c.Param("address")
	balance, err := h.WalletUsecase.GetBalance(c.Request().Context(), address)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"balance": balance})
}

func (h *walletHandlerImpl) SendMoney(c echo.Context) error {
	var request struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}
	if err := c.Bind(&request); err != nil || request.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	fromUUID, err := uuid.Parse(request.From)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid 'from' UUID"})
	}

	toUUID, err := uuid.Parse(request.To)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid 'to' UUID"})
	}

	err = h.WalletUsecase.SendMoney(c.Request().Context(), fromUUID.String(), toUUID.String(), request.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

func (h *walletHandlerImpl) GetAllWallets(c echo.Context) error {
	wallets, err := h.WalletUsecase.GetAllWallets(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, wallets)
}
