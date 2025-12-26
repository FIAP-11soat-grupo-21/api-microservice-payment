package queue

import (
	"context"
	"log"
	"payment_microservice/internal/common/config/env"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitMQChannel *amqp.Channel
	channelInstance *amqp.Channel
	channelOnce     sync.Once
)

var (
	rabbitMQConnection *amqp.Connection
	connectionInstance *amqp.Connection
	connectionOnce     sync.Once
)

func GetConnection() *amqp.Connection {
	connectionOnce.Do(func() {
		connectionInstance = rabbitMQConnection
	})
	return connectionInstance
}

func GetChannel() *amqp.Channel {
	channelOnce.Do(func() {
		channelInstance = rabbitMQChannel
	})
	return channelInstance
}

func Connect() {
	cfg := env.GetConfig()

	rabbitMQConnection = new(amqp.Connection)
	var err error
	maxRetries := 5
	retryInterval := 2 * time.Second

	for i := range maxRetries {
		rabbitMQConnection, err = amqp.Dial(cfg.RabbitMQ.URL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := rabbitMQConnection.Channel()

	if err != nil {
		log.Panicf("Failed to open a channel: %v", err)
	}

	rabbitMQChannel = ch
}

func Close() {
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}

	if rabbitMQConnection != nil {
		rabbitMQConnection.Close()
	}
}

func PublishMessageWithContext(ctx context.Context, routingKey string, body []byte) error {
	ch := GetChannel()

	cfg := env.GetConfig()

	err := ch.PublishWithContext(
		ctx,                   // context
		cfg.RabbitMQ.Exchange, // exchange
		routingKey,            // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	return err
}

func RegisterConsumer(queueName, binding string, handler func(amqp.Delivery)) {
	ch := GetChannel()
	cfg := env.GetConfig()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}

	err = ch.QueueBind(
		q.Name,                // queue name
		binding,               // routing key
		cfg.RabbitMQ.Exchange, // exchange
		false,                 // no-wait
		nil,                   // args
	)

	if err != nil {
		log.Panicf("Failed to bind a queue: %v", err)
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		log.Panicf("Failed to register a consumer: %v", err)
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	log.Printf(" [*] Waiting for messages in queue: %s. To exit press CTRL+C", queueName)
	<-forever
}
