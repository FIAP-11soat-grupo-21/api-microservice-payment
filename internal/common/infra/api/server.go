package api

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"payment_microservice/internal/adapters/driver/api/routes"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/infra/api/middlewares"
	_ "payment_microservice/internal/common/infra/api/swagger"
	"payment_microservice/internal/common/infra/database"
)

func Init() {
	config := env.GetConfig()

	if config.IsProduction() {
		log.Printf("Running in production mode on [%s]", config.API.URL)
		gin.SetMode(gin.ReleaseMode)
	}

	database.Connect()

	if config.Database.RunMigrations {
		database.RunMigrations()
	}

	ginRouter := gin.Default()

	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ginRouter.Use(gin.Logger())
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(middlewares.ErrorHandlerMiddleware())

	v1Routes := ginRouter.Group("/v1")

	routes.RegisterPaymentRoutes(v1Routes.Group("/payments"))

	ginRouter.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "healthy"})
	})

	ginRouter.Run(config.API.URL)
}
