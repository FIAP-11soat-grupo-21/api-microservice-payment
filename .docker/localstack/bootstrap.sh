#!/usr/bin/env bash
set -e

ORDER_ERROR_TOPIC=order-error-topic
CREATE_PAYMENT_QUEUE=create-payment-queue
CREATE_KITCHEN_ORDER_QUEUE=create-kitchen-order-queue
ORDER_ERROR_QUEUE=payment-order-error-queue

echo "▶ Creating SNS topic..."
TOPIC_ARN=$(awslocal sns create-topic \
  --name "$ORDER_ERROR_TOPIC" \
  --query 'TopicArn' \
  --output text)

echo "▶ Creating SQS queues..."
PAYMENT_QUEUE_URL=$(awslocal sqs create-queue \
  --queue-name "$CREATE_PAYMENT_QUEUE" \
  --query 'QueueUrl' \
  --output text)

KITCHEN_QUEUE_URL=$(awslocal sqs create-queue \
  --queue-name "$CREATE_KITCHEN_ORDER_QUEUE" \
  --query 'QueueUrl' \
  --output text)

ERROR_QUEUE_URL=$(awslocal sqs create-queue \
  --queue-name "$ORDER_ERROR_QUEUE" \
  --query 'QueueUrl' \
  --output text)

echo "▶ Getting Queue ARN..."
PAYMENT_QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url "$PAYMENT_QUEUE_URL" \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' \
  --output text)

echo "▶ Subscribing SQS to SNS..."
awslocal sns subscribe \
  --topic-arn "$TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$PAYMENT_QUEUE_ARN"

echo "✅ LocalStack bootstrap finished successfully"
