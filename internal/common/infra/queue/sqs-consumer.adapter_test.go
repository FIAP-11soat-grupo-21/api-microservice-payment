package queue

import (
	"context"
	"errors"
	"testing"
	"time"

	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/ports"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
)

type mockSQSClient struct {
	receiveOutput *sqs.ReceiveMessageOutput
	receiveErr    error
	deleteErr     error

	receiveCalls int
	deleteCalls  int
}

func (m *mockSQSClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	m.receiveCalls++
	return m.receiveOutput, m.receiveErr
}

func (m *mockSQSClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	m.deleteCalls++
	if m.deleteErr != nil {
		return nil, m.deleteErr
	}
	return &sqs.DeleteMessageOutput{}, nil
}

func stubAWSConfigLoaderSuccess(t *testing.T) {
	original := loadAWSConfig
	loadAWSConfig = func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
		return aws.Config{
			Region:      "us-east-1",
			Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("test-access", "test-secret", "")),
		}, nil
	}
	t.Cleanup(func() {
		loadAWSConfig = original
	})
}

func stubAWSConfigLoaderError(t *testing.T, err error) {
	original := loadAWSConfig
	loadAWSConfig = func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
		return aws.Config{}, err
	}
	t.Cleanup(func() {
		loadAWSConfig = original
	})
}

func stubQueueFatal(t *testing.T) *bool {
	called := new(bool)
	original := queueLogFatalf
	queueLogFatalf = func(format string, v ...interface{}) {
		*called = true
	}
	t.Cleanup(func() {
		queueLogFatalf = original
	})
	return called
}

func setupEnv(t *testing.T) func() {
	cleanup := env.SetupTestEnv(t)
	stubAWSConfigLoaderSuccess(t)
	return cleanup
}

func TestNewSQSConsumer_DefaultConfig(t *testing.T) {
	cleanup := setupEnv(t)
	defer cleanup()

	consumer := NewSQSConsumer()

	assert.NotNil(t, consumer)
	assert.NotNil(t, consumer.client)
	assert.NotNil(t, consumer.cfg)
	assert.NotNil(t, consumer.ctx)
	assert.NotNil(t, consumer.cancelFn)
	assert.False(t, consumer.isRunning)
}

func TestNewSQSConsumer_WithDifferentAWSConfigs(t *testing.T) {
	cleanup := setupEnv(t)
	defer cleanup()

	cfg := env.GetConfig()

	t.Run("with endpoint and explicit credentials", func(t *testing.T) {
		cfg.AWS.Endpoint = "http://localhost:9324"
		cfg.AWS.AccessKeyID = "access-key"
		cfg.AWS.SecretAccessKey = "secret-key"

		consumer := NewSQSConsumer()
		assert.NotNil(t, consumer)
	})

	t.Run("with endpoint and without explicit credentials", func(t *testing.T) {
		cfg.AWS.Endpoint = "http://localhost:9324"
		cfg.AWS.AccessKeyID = ""
		cfg.AWS.SecretAccessKey = ""

		consumer := NewSQSConsumer()
		assert.NotNil(t, consumer)
	})

	t.Run("without endpoint and without explicit credentials", func(t *testing.T) {
		cfg.AWS.Endpoint = ""
		cfg.AWS.AccessKeyID = ""
		cfg.AWS.SecretAccessKey = ""

		consumer := NewSQSConsumer()
		assert.NotNil(t, consumer)
	})
}

func TestNewSQSConsumer_LoadConfigError(t *testing.T) {
	cleanup := setupEnv(t)
	defer cleanup()

	stubAWSConfigLoaderError(t, errors.New("loader error"))
	fatalCalled := stubQueueFatal(t)

	consumer := NewSQSConsumer()
	assert.Nil(t, consumer)
	assert.True(t, *fatalCalled)
}

func TestConsumeQueue_SetsRunningAndStarts(t *testing.T) {
	cleanup := setupEnv(t)
	defer cleanup()

	consumer := NewSQSConsumer()
	mockClient := &mockSQSClient{
		receiveOutput: &sqs.ReceiveMessageOutput{Messages: []types.Message{}},
	}
	consumer.client = mockClient
	consumer.ctx, consumer.cancelFn = context.WithCancel(context.Background())

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789/test-queue"
	handler := func(body []byte) error { return nil }

	// iniciar consumo
	err := consumer.ConsumeQueue(queueURL, handler)
	assert.NoError(t, err)
	assert.True(t, consumer.isRunning)

	// encerrar a goroutine de polling
	consumer.cancelFn()
	time.Sleep(20 * time.Millisecond)
}

func TestConsumeQueue_AlreadyRunning(t *testing.T) {
	consumer := &SQSConsumer{
		client:    &mockSQSClient{},
		cfg:       &env.Config{},
		ctx:       context.Background(),
		cancelFn:  func() {},
		isRunning: true,
	}

	err := consumer.ConsumeQueue("queue-url", func(body []byte) error { return nil })
	assert.NoError(t, err)
	assert.True(t, consumer.isRunning)
}

func TestPollMessages_ProcessAndStopOnContextCancel(t *testing.T) {
	cleanup := setupEnv(t)
	defer cleanup()

	consumer := NewSQSConsumer()
	mockClient := &mockSQSClient{
		receiveOutput: &sqs.ReceiveMessageOutput{
			Messages: []types.Message{
				{
					Body:          aws.String("test-body"),
					MessageId:     aws.String("msg-1"),
					ReceiptHandle: aws.String("rh-1"),
				},
			},
		},
	}
	consumer.client = mockClient
	consumer.ctx, consumer.cancelFn = context.WithCancel(context.Background())
	consumer.isRunning = true

	processed := 0
	handler := func(body []byte) error {
		processed++
		consumer.cancelFn()
		return nil
	}

	consumer.pollMessages("queue-url", handler)

	assert.Equal(t, 1, processed)
	assert.GreaterOrEqual(t, mockClient.receiveCalls, 1)
	assert.GreaterOrEqual(t, mockClient.deleteCalls, 1)
	assert.False(t, consumer.isRunning)
}

func TestPollMessages_ReceiveErrorPath(t *testing.T) {
	cleanup := setupEnv(t)
	defer cleanup()

	originalSleep := sqsConsumerSleep
	sqsConsumerSleep = func(d time.Duration) {}
	defer func() { sqsConsumerSleep = originalSleep }()

	consumer := NewSQSConsumer()
	mockClient := &mockSQSClient{
		receiveErr: errors.New("receive error"),
	}
	consumer.client = mockClient
	consumer.ctx, consumer.cancelFn = context.WithCancel(context.Background())
	consumer.isRunning = true

	handler := func(body []byte) error { return nil }

	go func() {
		time.Sleep(20 * time.Millisecond)
		consumer.cancelFn()
	}()

	consumer.pollMessages("queue-url", handler)

	assert.True(t, mockClient.receiveCalls >= 1)
	assert.False(t, consumer.isRunning)
}

func TestPollMessages_ProcessMessageError(t *testing.T) {
	cleanup := setupEnv(t)
	defer cleanup()

	consumer := NewSQSConsumer()
	mockClient := &mockSQSClient{
		receiveOutput: &sqs.ReceiveMessageOutput{
			Messages: []types.Message{
				{
					Body:          aws.String("body"),
					MessageId:     aws.String("msg"),
					ReceiptHandle: aws.String("rh"),
				},
			},
		},
	}
	consumer.client = mockClient
	consumer.ctx, consumer.cancelFn = context.WithCancel(context.Background())
	consumer.isRunning = true

	handler := func(body []byte) error {
		consumer.cancelFn()
		return errors.New("handler failure")
	}

	consumer.pollMessages("queue-url", handler)

	assert.True(t, mockClient.receiveCalls >= 1)
	assert.Equal(t, 0, mockClient.deleteCalls)
	assert.False(t, consumer.isRunning)
}

func TestProcessMessage_Success(t *testing.T) {
	consumer := &SQSConsumer{
		client:   &mockSQSClient{},
		ctx:      context.Background(),
		cancelFn: func() {},
	}

	message := types.Message{
		Body:          aws.String("ok-body"),
		MessageId:     aws.String("msg-1"),
		ReceiptHandle: aws.String("rh-1"),
	}

	handlerCalled := false
	handler := func(body []byte) error {
		handlerCalled = true
		assert.Equal(t, "ok-body", string(body))
		return nil
	}

	err := consumer.processMessage("queue-url", message, handler)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

func TestProcessMessage_NilBody(t *testing.T) {
	consumer := &SQSConsumer{
		client:   &mockSQSClient{},
		ctx:      context.Background(),
		cancelFn: func() {},
	}

	message := types.Message{
		Body:          nil,
		MessageId:     aws.String("msg-1"),
		ReceiptHandle: aws.String("rh-1"),
	}

	err := consumer.processMessage("queue-url", message, func(body []byte) error {
		return errors.New("should not be called")
	})
	assert.NoError(t, err)
}

func TestProcessMessage_HandlerError(t *testing.T) {
	consumer := &SQSConsumer{
		client:   &mockSQSClient{},
		ctx:      context.Background(),
		cancelFn: func() {},
	}

	message := types.Message{
		Body:          aws.String("body"),
		MessageId:     aws.String("msg-1"),
		ReceiptHandle: aws.String("rh-1"),
	}

	expectedErr := errors.New("handler error")
	err := consumer.processMessage("queue-url", message, func(body []byte) error {
		return expectedErr
	})

	assert.Equal(t, expectedErr, err)
}

func TestProcessMessage_DeleteError(t *testing.T) {
	mockClient := &mockSQSClient{
		deleteErr: errors.New("delete error"),
	}
	consumer := &SQSConsumer{
		client:   mockClient,
		ctx:      context.Background(),
		cancelFn: func() {},
	}

	message := types.Message{
		Body:          aws.String("body"),
		MessageId:     aws.String("msg-1"),
		ReceiptHandle: aws.String("rh-1"),
	}

	err := consumer.processMessage("queue-url", message, func(body []byte) error { return nil })
	assert.Equal(t, mockClient.deleteErr, err)
}

func TestConsumeQueue_AcceptsMessageHandlerType(t *testing.T) {
	cleanup := setupEnv(t)
	defer cleanup()

	consumer := NewSQSConsumer()
	consumer.client = &mockSQSClient{
		receiveOutput: &sqs.ReceiveMessageOutput{Messages: []types.Message{}},
	}
	consumer.ctx, consumer.cancelFn = context.WithCancel(context.Background())

	var handler ports.MessageHandler = func(body []byte) error { return nil }

	err := consumer.ConsumeQueue("queue-url", handler)
	assert.NoError(t, err)
	consumer.cancelFn()
}
