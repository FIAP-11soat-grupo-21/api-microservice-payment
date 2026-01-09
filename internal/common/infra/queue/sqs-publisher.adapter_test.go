package queue

import (
	"context"
	"errors"
	"testing"

	"payment_microservice/internal/common/config/env"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
)

type mockSQSSender struct {
	lastInput *sqs.SendMessageInput
	err       error
	calls     int
}

func (m *mockSQSSender) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	m.calls++
	m.lastInput = params
	if m.err != nil {
		return nil, m.err
	}
	return &sqs.SendMessageOutput{}, nil
}

type mockSNSPublisher struct {
	lastInput *sns.PublishInput
	err       error
	calls     int
}

func (m *mockSNSPublisher) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.calls++
	m.lastInput = params
	if m.err != nil {
		return nil, m.err
	}
	return &sns.PublishOutput{}, nil
}

func setupPublisherEnv(t *testing.T) func() {
	cleanup := env.SetupTestEnv(t)
	stubAWSConfigLoaderSuccess(t)
	return cleanup
}

func TestNewSQSPublisher(t *testing.T) {
	cleanup := setupPublisherEnv(t)
	defer cleanup()

	publisher := NewSQSPublisher()

	assert.NotNil(t, publisher)
	assert.NotNil(t, publisher.sqsClient)
	assert.NotNil(t, publisher.snsClient)
	assert.NotNil(t, publisher.cfg)
}

func TestNewSQSPublisher_WithDifferentAWSConfigs(t *testing.T) {
	cleanup := setupPublisherEnv(t)
	defer cleanup()

	cfg := env.GetConfig()

	t.Run("with endpoint and explicit credentials", func(t *testing.T) {
		cfg.AWS.Endpoint = "http://localhost:9324"
		cfg.AWS.AccessKeyID = "access-key"
		cfg.AWS.SecretAccessKey = "secret-key"

		publisher := NewSQSPublisher()
		assert.NotNil(t, publisher)
	})

	t.Run("with endpoint and without explicit credentials", func(t *testing.T) {
		cfg.AWS.Endpoint = "http://localhost:9324"
		cfg.AWS.AccessKeyID = ""
		cfg.AWS.SecretAccessKey = ""

		publisher := NewSQSPublisher()
		assert.NotNil(t, publisher)
	})

	t.Run("without endpoint and without explicit credentials", func(t *testing.T) {
		cfg.AWS.Endpoint = ""
		cfg.AWS.AccessKeyID = ""
		cfg.AWS.SecretAccessKey = ""

		publisher := NewSQSPublisher()
		assert.NotNil(t, publisher)
	})
}

func TestNewSQSPublisher_LoadConfigError(t *testing.T) {
	cleanup := setupPublisherEnv(t)
	defer cleanup()

	stubAWSConfigLoaderError(t, errors.New("loader error"))
	fatalCalled := stubQueueFatal(t)

	publisher := NewSQSPublisher()
	assert.Nil(t, publisher)
	assert.True(t, *fatalCalled)
}

func TestPublishOnQueue_Success(t *testing.T) {
	cleanup := setupPublisherEnv(t)
	defer cleanup()

	mockSQS := &mockSQSSender{}
	publisher := &SQSPublisher{
		sqsClient: mockSQS,
		snsClient: &mockSNSPublisher{},
		cfg:       env.GetConfig(),
	}

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789/test"
	body := []byte("hello world")

	err := publisher.PublishOnQueue(context.Background(), queueURL, body)

	assert.NoError(t, err)
	assert.Equal(t, 1, mockSQS.calls)
	if assert.NotNil(t, mockSQS.lastInput) {
		assert.Equal(t, queueURL, aws.ToString(mockSQS.lastInput.QueueUrl))
		assert.Equal(t, string(body), aws.ToString(mockSQS.lastInput.MessageBody))
	}
}

func TestPublishOnQueue_Error(t *testing.T) {
	cleanup := setupPublisherEnv(t)
	defer cleanup()

	expectedErr := errors.New("send error")
	mockSQS := &mockSQSSender{err: expectedErr}
	publisher := &SQSPublisher{
		sqsClient: mockSQS,
		snsClient: &mockSNSPublisher{},
		cfg:       env.GetConfig(),
	}

	err := publisher.PublishOnQueue(context.Background(), "queue-url", []byte("payload"))

	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, mockSQS.calls)
}

func TestPublishOnTopic_Success(t *testing.T) {
	cleanup := setupPublisherEnv(t)
	defer cleanup()

	mockSNS := &mockSNSPublisher{}
	publisher := &SQSPublisher{
		sqsClient: &mockSQSSender{},
		snsClient: mockSNS,
		cfg:       env.GetConfig(),
	}

	topicARN := "arn:aws:sns:us-east-1:123456789012:topic"
	body := []byte("topic message")

	err := publisher.PublishOnTopic(context.Background(), topicARN, body)

	assert.NoError(t, err)
	assert.Equal(t, 1, mockSNS.calls)
	if assert.NotNil(t, mockSNS.lastInput) {
		assert.Equal(t, topicARN, aws.ToString(mockSNS.lastInput.TopicArn))
		assert.Equal(t, string(body), aws.ToString(mockSNS.lastInput.Message))
	}
}

func TestPublishOnTopic_Error(t *testing.T) {
	cleanup := setupPublisherEnv(t)
	defer cleanup()

	expectedErr := errors.New("publish error")
	mockSNS := &mockSNSPublisher{err: expectedErr}
	publisher := &SQSPublisher{
		sqsClient: &mockSQSSender{},
		snsClient: mockSNS,
		cfg:       env.GetConfig(),
	}

	err := publisher.PublishOnTopic(context.Background(), "arn:test", []byte("payload"))

	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, mockSNS.calls)
}
