package routes

import (
	"payment_microservice/internal/common/config/env"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterPaymentRoutes(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	paymentGroup := router.Group("/v1/payments")

	RegisterPaymentRoutes(paymentGroup)

	// Verificar se as rotas foram registradas
	routes := router.Routes()

	// Deve ter 3 rotas registradas
	assert.Len(t, routes, 3, "Should have 3 payment routes registered")

	// Verificar se cada rota esperada está presente
	routePaths := make(map[string]bool)
	for _, route := range routes {
		routePaths[route.Method+":"+route.Path] = true
	}

	assert.True(t, routePaths["POST:/v1/payments/pix"], "Should have POST /v1/payments/pix route")
	assert.True(t, routePaths["POST:/v1/payments/webhook"], "Should have POST /v1/payments/webhook route")
	assert.True(t, routePaths["GET:/v1/payments/order/:orderID"], "Should have GET /v1/payments/order/:orderID route")
}

func TestRegisterPaymentRoutes_EndpointsAccessible(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	paymentGroup := router.Group("/v1/payments")

	RegisterPaymentRoutes(paymentGroup)

	// Verificar que as rotas foram registradas e podem ser acessadas
	routes := router.Routes()

	// Deve ter handlers para cada rota
	for _, route := range routes {
		assert.NotNil(t, route.Handler, "Route %s %s should have a handler", route.Method, route.Path)
		assert.NotEmpty(t, route.Handler, "Route %s %s handler should not be empty", route.Method, route.Path)
	}
}

func TestRegisterPaymentRoutes_GroupPrefix(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Test com prefixo diferente
	customGroup := router.Group("/api/v2/custom-payments")
	RegisterPaymentRoutes(customGroup)

	routes := router.Routes()

	// Verificar que as rotas usam o prefixo customizado
	routePaths := make([]string, 0)
	for _, route := range routes {
		routePaths = append(routePaths, route.Path)
	}

	assert.Contains(t, routePaths, "/api/v2/custom-payments/pix")
	assert.Contains(t, routePaths, "/api/v2/custom-payments/webhook")
	assert.Contains(t, routePaths, "/api/v2/custom-payments/order/:orderID")
}
