package http_errors

import (
	"net/http"
	"payment_microservice/internal/core/domain/exceptions"

	"github.com/gin-gonic/gin"
)

func HandleDomainErrors(err error, ctx *gin.Context) bool {
	switch e := err.(type) {
	case *exceptions.PaymentNotFoundException:
		ctx.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		return true

	case *exceptions.InvalidPaymentDataException:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return true
	}

	return false
}
