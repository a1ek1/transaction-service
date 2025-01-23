package router

import (
	"github.com/labstack/echo"
	"transaction-service/internal/presenter/http/handler"
)

func NewRouter(e *echo.Echo, h handler.AppHandler) {
	api := e.Group("/api")
	{
		api.POST("/send", h.SendMoney)
		api.GET("/transactions", h.GetLastTransactions)
		api.GET("/wallet/:address/balance", h.GetBalance)
	}
}
