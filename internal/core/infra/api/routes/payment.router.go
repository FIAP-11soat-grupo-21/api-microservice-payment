package routes

import (
	"github.com/gin-gonic/gin"

	"payment_microservice/internal/core/infra/api/handlers"
)

func RegisterPaymentRoutes(router *gin.RouterGroup) {
	paymentHandler := handlers.NewPaymentHandler()

	router.POST("/pix", paymentHandler.CreatePixBilling)
	router.POST("/webhook", paymentHandler.ConfirmPayment)
	router.GET("/order/:orderID", paymentHandler.FindByOrderID)
}
