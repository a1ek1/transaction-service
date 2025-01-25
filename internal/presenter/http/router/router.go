// Package router configures HTTP routes for the transaction service.
package router

import (
	"github.com/labstack/echo"
	"transaction-service/internal/presenter/http/handler"
)

// NewRouter sets up routes for the transaction service.
func NewRouter(e *echo.Echo, h handler.AppHandler) {
	api := e.Group("/api")
	{
		api.POST("/send", h.SendMoney)
		api.GET("/transactions", h.GetLastTransactions)
		api.GET("/wallets", h.GetAllWallets)
		api.GET("/wallet/:address/balance", h.GetBalance)
	}
}
