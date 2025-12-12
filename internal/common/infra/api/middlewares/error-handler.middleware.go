package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	payment_http_errors "payment_microservice/internal/core/infra/api/http_errors"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			errorHasBinHandled := payment_http_errors.HandleDomainErrors(err, ctx)

			if !errorHasBinHandled {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}

			ctx.Abort()
		}
	}
}
