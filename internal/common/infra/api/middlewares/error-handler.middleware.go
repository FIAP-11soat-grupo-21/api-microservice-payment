package middlewares

import (
	"net/http"
	"payment_microservice/internal/adapters/driver/api/http_errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			errorHasBinHandled := http_errors.HandleDomainErrors(err, ctx)

			if !errorHasBinHandled {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}

			ctx.Abort()
		}
	}
}
