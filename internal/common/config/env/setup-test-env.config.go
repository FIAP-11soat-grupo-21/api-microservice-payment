package env

import (
	"os"
	"testing"
)

func SetupTestEnv(t *testing.T) func() {
	// Salva os valores originais
	originalValues := map[string]string{
		"GO_ENV":                              os.Getenv("GO_ENV"),
		"API_PORT":                            os.Getenv("API_PORT"),
		"API_HOST":                            os.Getenv("API_HOST"),
		"WEBHOOK_URL":                         os.Getenv("WEBHOOK_URL"),
		"DB_HOST":                             os.Getenv("DB_HOST"),
		"DB_NAME":                             os.Getenv("DB_NAME"),
		"DB_PORT":                             os.Getenv("DB_PORT"),
		"DB_USERNAME":                         os.Getenv("DB_USERNAME"),
		"DB_PASSWORD":                         os.Getenv("DB_PASSWORD"),
		"DB_RUN_MIGRATIONS":                   os.Getenv("DB_RUN_MIGRATIONS"),
		"MERCADOPAGO_ACCESS_TOKEN":            os.Getenv("MERCADOPAGO_ACCESS_TOKEN"),
		"MERCADOPAGO_COLLECTOR_ID":            os.Getenv("MERCADOPAGO_COLLECTOR_ID"),
		"MERCADOPAGO_EXTERNAL_POS_ID":         os.Getenv("MERCADOPAGO_EXTERNAL_POS_ID"),
		"MERCADOPAGO_API_URL":                 os.Getenv("MERCADOPAGO_API_URL"),
		"RABBITMQ_URL":                        os.Getenv("RABBITMQ_URL"),
		"RABBITMQ_EXCHANGE":                   os.Getenv("RABBITMQ_EXCHANGE"),
		"RABBITMQ_CREATE_PAYMENT_TOPIC":       os.Getenv("RABBITMQ_CREATE_PAYMENT_TOPIC"),
		"RABBITMQ_CREATE_KITCHEN_ORDER_TOPIC": os.Getenv("RABBITMQ_CREATE_KITCHEN_ORDER_TOPIC"),
		"RABBITMQ_ORDER_ERROR_TOPIC":          os.Getenv("RABBITMQ_ORDER_ERROR_TOPIC"),
	}

	// Define valores de teste
	os.Setenv("GO_ENV", "test")
	os.Setenv("API_PORT", "8080")
	os.Setenv("API_HOST", "0.0.0.0")
	os.Setenv("WEBHOOK_URL", "http://localhost:8080/v1/payments/webhook")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	os.Setenv("MERCADOPAGO_ACCESS_TOKEN", "TEST_ACCESS_TOKEN")
	os.Setenv("MERCADOPAGO_COLLECTOR_ID", "1234567890")
	os.Setenv("MERCADOPAGO_EXTERNAL_POS_ID", "test_pos_id")
	os.Setenv("MERCADOPAGO_API_URL", "https://api.test.com")
	os.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	os.Setenv("RABBITMQ_EXCHANGE", "test_exchange")
	os.Setenv("RABBITMQ_CREATE_PAYMENT_TOPIC", "create.payment")
	os.Setenv("RABBITMQ_CREATE_KITCHEN_ORDER_TOPIC", "create.kitchen-order")
	os.Setenv("RABBITMQ_ORDER_ERROR_TOPIC", "order-error")

	// Retorna função de cleanup
	return func() {
		for key, value := range originalValues {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}
}
