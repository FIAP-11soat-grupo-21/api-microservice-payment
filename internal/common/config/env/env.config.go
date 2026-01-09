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
	AWS struct {
		Endpoint        string
		Region          string
		AccessKeyID     string
		SecretAccessKey string
		SQS             struct {
			Queues struct {
				OrderError         string
				CreateKitchenOrder string
				CreatePayment      string
			}
		}
		SNS struct {
			Topics struct {
				OrderError string
			}
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

	// API
	c.API.Port = getEnv("API_PORT")
	c.API.Host = getEnv("API_HOST")
	c.API.URL = c.API.Host + ":" + c.API.Port
	c.API.WebhookURL = getEnv("WEBHOOK_URL")

	// Database
	c.Database.RunMigrations = getEnv("DB_RUN_MIGRATIONS") == "true"
	c.Database.Host = getEnv("DB_HOST")
	c.Database.Name = getEnv("DB_NAME")
	c.Database.Port = getEnv("DB_PORT")
	c.Database.Username = getEnv("DB_USERNAME")
	c.Database.Password = getEnv("DB_PASSWORD")

	// Mercado Pago API
	c.MercadoPago.AccessToken = getEnv("MERCADOPAGO_ACCESS_TOKEN")
	c.MercadoPago.CollectorID = getEnv("MERCADOPAGO_COLLECTOR_ID")
	c.MercadoPago.ExternalPosID = getEnv("MERCADOPAGO_EXTERNAL_POS_ID")
	c.MercadoPago.ApiBaseURL = getEnv("MERCADOPAGO_API_URL")

	// AWS
	c.AWS.Region = getEnv("AWS_REGION")
	c.AWS.Endpoint = os.Getenv("AWS_ENDPOINT")                 // Optional
	c.AWS.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")         // Optional
	c.AWS.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY") // Optional

	// SQS
	c.AWS.SQS.Queues.CreatePayment = getEnv("AWS_SQS_CREATE_PAYMENT_QUEUE")
	c.AWS.SQS.Queues.CreateKitchenOrder = getEnv("AWS_SQS_CREATE_KITCHEN_ORDER_QUEUE")
	c.AWS.SQS.Queues.OrderError = getEnv("AWS_SQS_ORDER_ERROR_QUEUE")

	// SNS
	c.AWS.SNS.Topics.OrderError = getEnv("AWS_SNS_ORDER_ERROR_TOPIC")
}

func (c *Config) IsProduction() bool {
	return c.GoEnv == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.GoEnv == "development"
}
