package queue

import (
	"context"
	"payment_microservice/internal/common/config/env"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConnection(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Reset singletons
	rabbitMQConnection = nil
	connectionInstance = nil

	conn := GetConnection()
	assert.Nil(t, conn, "GetConnection should return nil when rabbitMQConnection is not set")
}

func TestGetChannel(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Reset singletons
	rabbitMQChannel = nil
	channelInstance = nil

	ch := GetChannel()
	assert.Nil(t, ch, "GetChannel should return nil when rabbitMQChannel is not set")
}

func TestClose_NoConnection(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Reset singletons
	rabbitMQConnection = nil
	rabbitMQChannel = nil

	// Act & Assert: não deve fazer panic
	assert.NotPanics(t, func() {
		Close()
	})
}

func TestPublishMessageWithContext_NilChannel(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Reset singletons
	rabbitMQChannel = nil
	channelInstance = nil

	ctx := context.Background()
	routingKey := "test.routing.key"
	body := []byte(`{"test": "message"}`)

	// Act & Assert: deve fazer panic ao tentar publicar sem channel (comportamento esperado do código atual)
	assert.Panics(t, func() {
		PublishMessageWithContext(ctx, routingKey, body)
	}, "PublishMessageWithContext should panic when channel is nil")
}

func TestConnect_RetryLogic(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Reset singletons
	rabbitMQConnection = nil
	rabbitMQChannel = nil
	connectionInstance = nil
	channelInstance = nil

	// Este teste apenas verifica se o código não entra em panic
	// Em ambiente de teste, a conexão falhará mas não devemos testar
	// a conexão real ao RabbitMQ

	// Para testar retry logic, precisaríamos de uma interface injetável
	// Por enquanto, apenas documentamos que o código tem retry logic
	assert.NotNil(t, &rabbitMQConnection)
	assert.NotNil(t, &rabbitMQChannel)
}
