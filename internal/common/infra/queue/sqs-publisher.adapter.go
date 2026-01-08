package queue

import (
	"context"
	"log"
	"payment_microservice/internal/common/config/env"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSPublisher struct {
	sqsClient *sqs.Client
	snsClient *sns.Client
	cfg       *env.Config
}

func NewSQSPublisher() *SQSPublisher {
	ctx := context.Background()

	appCfg := env.GetConfig()

	var awsCfg aws.Config
	var err error

	// Configura um resolver de endpoint customizado, se fornecido (ex.: ElasticMQ)
	var endpointResolver aws.EndpointResolverWithOptions
	if appCfg.AWS.SQS.Endpoint != "" {
		endpointResolver = aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if service == sqs.ServiceID || service == sns.ServiceID {
					return aws.Endpoint{
						URL: appCfg.AWS.SQS.Endpoint,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			},
		)
	}

	if appCfg.AWS.AccessKeyID != "" && appCfg.AWS.SecretAccessKey != "" {
		// Usa credenciais explícitas se fornecidas
		if endpointResolver != nil {
			awsCfg, err = config.LoadDefaultConfig(ctx,
				config.WithRegion(appCfg.AWS.Region),
				config.WithEndpointResolverWithOptions(endpointResolver),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
					appCfg.AWS.AccessKeyID,
					appCfg.AWS.SecretAccessKey,
					"",
				)),
			)
		} else {
			awsCfg, err = config.LoadDefaultConfig(ctx,
				config.WithRegion(appCfg.AWS.Region),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
					appCfg.AWS.AccessKeyID,
					appCfg.AWS.SecretAccessKey,
					"",
				)),
			)
		}
	} else {
		// Usa credenciais padrão (IAM role, environment variables, etc)
		if endpointResolver != nil {
			awsCfg, err = config.LoadDefaultConfig(ctx,
				config.WithRegion(appCfg.AWS.Region),
				config.WithEndpointResolverWithOptions(endpointResolver),
			)
		} else {
			awsCfg, err = config.LoadDefaultConfig(ctx,
				config.WithRegion(appCfg.AWS.Region),
			)
		}
	}

	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
		return nil
	}

	sqsClient := sqs.NewFromConfig(awsCfg)

	snsClient := sns.NewFromConfig(awsCfg)

	return &SQSPublisher{
		sqsClient: sqsClient,
		snsClient: snsClient,
		cfg:       appCfg,
	}
}

func (p *SQSPublisher) PublishOnQueue(ctx context.Context, queueName string, message []byte) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueName),
		MessageBody: aws.String(string(message)),
	}

	_, err := p.sqsClient.SendMessage(ctx, input)
	if err != nil {
		log.Printf("Failed to publish message to SQS: %v", err)
		return err
	}

	return nil
}

func (p *SQSPublisher) PublishOnTopic(ctx context.Context, topic string, message []byte) error {
	input := &sns.PublishInput{
		TopicArn: aws.String(topic),
		Message:  aws.String(string(message)),
	}

	_, err := p.snsClient.Publish(ctx, input)
	if err != nil {
		log.Printf("Failed to publish message to SNS topic: %v", err)
		return err
	}

	return nil
}
