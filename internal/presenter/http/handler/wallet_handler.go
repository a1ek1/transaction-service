package handler

import (
	"net/http"
	"transaction-service/internal/usecase"

	"github.com/labstack/echo"
)

// WalletHandler интерфейс для обработки запросов кошельков
type WalletHandler interface {
	GetBalance(c echo.Context) error
	SendMoney(c echo.Context) error
}

// Реализация WalletHandler
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
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := h.WalletUsecase.SendMoney(c.Request().Context(), request.From, request.To, request.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}
