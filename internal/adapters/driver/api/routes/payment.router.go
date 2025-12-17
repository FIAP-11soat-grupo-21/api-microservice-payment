package routes

import (
	"payment_microservice/internal/adapters/driver/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(router *gin.RouterGroup) {
	paymentHandler := handlers.NewPaymentHandler()

	router.POST("/pix", paymentHandler.CreatePixBilling)
	router.POST("/webhook", paymentHandler.ConfirmPayment)
	router.GET("/order/:orderID", paymentHandler.FindByOrderID)
}
