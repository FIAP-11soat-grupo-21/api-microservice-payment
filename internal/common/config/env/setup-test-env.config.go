package env

import (
	"os"
	"testing"
)

func SetupTestEnv(t *testing.T) func() {
	// Salva os valores originais
	originalValues := map[string]string{
		"GO_ENV":                             os.Getenv("GO_ENV"),
		"API_PORT":                           os.Getenv("API_PORT"),
		"API_HOST":                           os.Getenv("API_HOST"),
		"WEBHOOK_URL":                        os.Getenv("WEBHOOK_URL"),
		"DB_HOST":                            os.Getenv("DB_HOST"),
		"DB_NAME":                            os.Getenv("DB_NAME"),
		"DB_PORT":                            os.Getenv("DB_PORT"),
		"DB_USERNAME":                        os.Getenv("DB_USERNAME"),
		"DB_PASSWORD":                        os.Getenv("DB_PASSWORD"),
		"DB_RUN_MIGRATIONS":                  os.Getenv("DB_RUN_MIGRATIONS"),
		"MERCADOPAGO_ACCESS_TOKEN":           os.Getenv("MERCADOPAGO_ACCESS_TOKEN"),
		"MERCADOPAGO_COLLECTOR_ID":           os.Getenv("MERCADOPAGO_COLLECTOR_ID"),
		"MERCADOPAGO_EXTERNAL_POS_ID":        os.Getenv("MERCADOPAGO_EXTERNAL_POS_ID"),
		"MERCADOPAGO_API_URL":                os.Getenv("MERCADOPAGO_API_URL"),
		"AWS_REGION":                         os.Getenv("AWS_REGION"),
		"AWS_ACCESS_KEY_ID":                  os.Getenv("AWS_ACCESS_KEY_ID"),
		"AWS_SECRET_ACCESS_KEY":              os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"AWS_SQS_CREATE_PAYMENT_QUEUE":       os.Getenv("AWS_SQS_CREATE_PAYMENT_QUEUE"),
		"AWS_SQS_CREATE_KITCHEN_ORDER_QUEUE": os.Getenv("AWS_SQS_CREATE_KITCHEN_ORDER_QUEUE"),
		"AWS_SQS_ORDER_ERROR_QUEUE":          os.Getenv("AWS_SQS_ORDER_ERROR_QUEUE"),
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
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_ACCESS_KEY_ID", "test_access_key_id")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test_secret_access_key")
	os.Setenv("AWS_SQS_CREATE_PAYMENT_QUEUE", "test_create_payment_queue")
	os.Setenv("AWS_SQS_CREATE_KITCHEN_ORDER_QUEUE", "test_create_kitchen_order_queue")
	os.Setenv("AWS_SQS_ORDER_ERROR_QUEUE", "test_order_error_queue")

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
