package env

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	GoEnv string
	API   struct {
		Port       string
		Host       string
		URL        string
		WebhookURL string
	}
	Database struct {
		RunMigrations bool
		Host          string
		Name          string
		Port          string
		Username      string
		Password      string
	}
	MercadoPago struct {
		AccessToken   string
		CollectorID   string
		ExternalPosID string
		ApiBaseURL    string
	}
	RabbitMQ struct {
		URL      string
		Exchange string
		Topics   struct {
			OrderError         string
			CreateKitchenOrder string
			CreatePayment      string
		}
	}
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		instance.Load()
	})
	return instance
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func (c *Config) Load() {
	dotEnvPath := ".env"
	_, err := os.Stat(dotEnvPath)

	if err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	c.GoEnv = getEnv("GO_ENV")

	c.API.Port = getEnv("API_PORT")
	c.API.Host = getEnv("API_HOST")
	c.API.URL = c.API.Host + ":" + c.API.Port
	c.API.WebhookURL = getEnv("WEBHOOK_URL")

	c.Database.RunMigrations = getEnv("DB_RUN_MIGRATIONS") == "true"
	c.Database.Host = getEnv("DB_HOST")
	c.Database.Name = getEnv("DB_NAME")
	c.Database.Port = getEnv("DB_PORT")
	c.Database.Username = getEnv("DB_USERNAME")
	c.Database.Password = getEnv("DB_PASSWORD")

	c.MercadoPago.AccessToken = getEnv("MERCADOPAGO_ACCESS_TOKEN")
	c.MercadoPago.CollectorID = getEnv("MERCADOPAGO_COLLECTOR_ID")
	c.MercadoPago.ExternalPosID = getEnv("MERCADOPAGO_EXTERNAL_POS_ID")
	c.MercadoPago.ApiBaseURL = getEnv("MERCADOPAGO_API_URL")

	c.RabbitMQ.URL = getEnv("RABBITMQ_URL")
	c.RabbitMQ.Exchange = getEnv("RABBITMQ_EXCHANGE")
	c.RabbitMQ.Topics.CreatePayment = getEnv("RABBITMQ_CREATE_PAYMENT_TOPIC")
	c.RabbitMQ.Topics.CreateKitchenOrder = getEnv("RABBITMQ_CREATE_KITCHEN_ORDER_TOPIC")
	c.RabbitMQ.Topics.OrderError = getEnv("RABBITMQ_ORDER_ERROR_TOPIC")
}

func (c *Config) IsProduction() bool {
	return c.GoEnv == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.GoEnv == "development"
}
