package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	t.Run("should return singleton config instance", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		config1 := GetConfig()
		config2 := GetConfig()

		assert.NotNil(t, config1)
		assert.NotNil(t, config2)
		assert.Same(t, config1, config2, "GetConfig should return the same instance (singleton)")
	})
}

func TestConfig_Load(t *testing.T) {
	t.Run("should load all environment variables correctly", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		config := &Config{}
		config.Load()

		assert.Equal(t, "test", config.GoEnv)
		assert.Equal(t, "8080", config.API.Port)
		assert.Equal(t, "0.0.0.0", config.API.Host)
		assert.Equal(t, "0.0.0.0:8080", config.API.URL)
		assert.Equal(t, "http://localhost:8080/v1/payments/webhook", config.API.WebhookURL)
		assert.False(t, config.Database.RunMigrations)
		assert.Equal(t, "localhost", config.Database.Host)
		assert.Equal(t, "test_db", config.Database.Name)
		assert.Equal(t, "5432", config.Database.Port)
		assert.Equal(t, "test_user", config.Database.Username)
		assert.Equal(t, "test_password", config.Database.Password)
		assert.Equal(t, "TEST_ACCESS_TOKEN", config.MercadoPago.AccessToken)
		assert.Equal(t, "1234567890", config.MercadoPago.CollectorID)
		assert.Equal(t, "test_pos_id", config.MercadoPago.ExternalPosID)
		assert.Equal(t, "https://api.test.com", config.MercadoPago.ApiBaseURL)
		assert.Equal(t, "amqp://guest:guest@localhost:5672/", config.RabbitMQ.URL)
		assert.Equal(t, "test_exchange", config.RabbitMQ.Exchange)
		assert.Equal(t, "create.kitchen-order", config.RabbitMQ.Topics.CreateKitchenOrder)
	})

	t.Run("should set RunMigrations to true when DB_RUN_MIGRATIONS is true", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("DB_RUN_MIGRATIONS", "true")

		config := &Config{}
		config.Load()

		assert.True(t, config.Database.RunMigrations)
	})

	t.Run("should set RunMigrations to false when DB_RUN_MIGRATIONS is false", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("DB_RUN_MIGRATIONS", "false")

		config := &Config{}
		config.Load()

		assert.False(t, config.Database.RunMigrations)
	})

	t.Run("should set RunMigrations to false when DB_RUN_MIGRATIONS is any other value", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("DB_RUN_MIGRATIONS", "invalid")

		config := &Config{}
		config.Load()

		assert.False(t, config.Database.RunMigrations)
	})

	t.Run("should construct API.URL from Host and Port", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("API_HOST", "127.0.0.1")
		os.Setenv("API_PORT", "9090")

		config := &Config{}
		config.Load()

		assert.Equal(t, "127.0.0.1", config.API.Host)
		assert.Equal(t, "9090", config.API.Port)
		assert.Equal(t, "127.0.0.1:9090", config.API.URL)
	})
}

func TestConfig_IsProduction(t *testing.T) {
	t.Run("should return true when GO_ENV is production", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("GO_ENV", "production")
		config := &Config{}
		config.Load()

		result := config.IsProduction()

		assert.True(t, result)
	})

	t.Run("should return false when GO_ENV is not production", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("GO_ENV", "development")
		config := &Config{}
		config.Load()

		result := config.IsProduction()

		assert.False(t, result)
	})

	t.Run("should return false when GO_ENV is test", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("GO_ENV", "test")
		config := &Config{}
		config.Load()

		result := config.IsProduction()

		assert.False(t, result)
	})
}

func TestConfig_IsDevelopment(t *testing.T) {
	t.Run("should return true when GO_ENV is development", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("GO_ENV", "development")
		config := &Config{}
		config.Load()

		result := config.IsDevelopment()

		assert.True(t, result)
	})

	t.Run("should return false when GO_ENV is not development", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("GO_ENV", "production")
		config := &Config{}
		config.Load()

		result := config.IsDevelopment()

		assert.False(t, result)
	})

	t.Run("should return false when GO_ENV is test", func(t *testing.T) {
		cleanup := SetupTestEnv(t)
		defer cleanup()

		os.Setenv("GO_ENV", "test")
		config := &Config{}
		config.Load()

		result := config.IsDevelopment()

		assert.False(t, result)
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("should return environment variable value when set", func(t *testing.T) {
		key := "TEST_VAR_EXISTS"
		expectedValue := "test_value"
		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key)

		value := getEnv(key)

		assert.Equal(t, expectedValue, value)
	})
}
