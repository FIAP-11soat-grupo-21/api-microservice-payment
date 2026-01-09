package queue

import (
	"context"
	"log"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/ports"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type sqsAPI interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

var sqsConsumerSleep = time.Sleep

type SQSConsumer struct {
	client    sqsAPI
	cfg       *env.Config
	ctx       context.Context
	cancelFn  context.CancelFunc
	isRunning bool
}

func NewSQSConsumer() *SQSConsumer {
	ctx := context.Background()
	appCfg := env.GetConfig()

	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(appCfg.AWS.Region),
		config.WithBaseEndpoint(appCfg.AWS.Endpoint),
	}

	if appCfg.AWS.AccessKeyID != "" && appCfg.AWS.SecretAccessKey != "" {
		optFns = append(optFns, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				appCfg.AWS.AccessKeyID,
				appCfg.AWS.SecretAccessKey,
				"",
			),
		))
	}

	awsCfg, err := loadAWSConfig(ctx, optFns...)

	if err != nil {
		queueLogFatalf("unable to load AWS SDK config, %v", err)
		return nil
	}

	client := sqs.NewFromConfig(awsCfg)

	consumerCtx, cancel := context.WithCancel(ctx)

	return &SQSConsumer{
		client:    client,
		cfg:       appCfg,
		ctx:       consumerCtx,
		cancelFn:  cancel,
		isRunning: false,
	}
}

func (c *SQSConsumer) ConsumeQueue(queueURL string, handler ports.MessageHandler) error {
	if c.isRunning {
		log.Printf("Consumer already running for queue: %s", queueURL)
		return nil
	}

	c.isRunning = true
	log.Printf(" [*] Starting to consume messages from SQS queue: %s", queueURL)

	go c.pollMessages(queueURL, handler)

	return nil
}

func (c *SQSConsumer) pollMessages(queueURL string, handler ports.MessageHandler) {
	for {
		select {
		case <-c.ctx.Done():
			log.Printf("Stopping consumer for queue: %s", queueURL)
			c.isRunning = false
			return
		default:
			result, err := c.client.ReceiveMessage(c.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueURL),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20, // Long polling
				VisibilityTimeout:   30,
			})

			if err != nil {
				log.Println("QueueURL ->", queueURL)
				log.Printf("Error receiving messages from SQS: %v", err)
				sqsConsumerSleep(5 * time.Second)
				continue
			}

			for _, message := range result.Messages {
				err = c.processMessage(queueURL, message, handler)
				if err != nil {
					log.Printf("Error processing message ID %s: %v", *message.MessageId, err)
				}
			}
		}
	}
}

func (c *SQSConsumer) processMessage(queueURL string, message types.Message, handler ports.MessageHandler) error {
	if message.Body == nil {
		log.Printf("Received message with nil body")
		return nil
	}

	err := handler([]byte(*message.Body))

	if err != nil {
		return err
	}

	_, err = c.client.DeleteMessage(c.ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: message.ReceiptHandle,
	})

	if err != nil {
		return err
	}

	return nil
}
