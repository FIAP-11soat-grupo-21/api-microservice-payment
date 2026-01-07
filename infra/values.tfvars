application_name = "payment-api"
image_name       = "GHCR_IMAGE_TAG"
image_port       = 8088
app_path_pattern = ["/payments*"]

# =======================================================
# Configurações do ECS Service
# =======================================================
container_environment_variables = {
  GO_ENV : "production"
  API_PORT : "8088"
  API_HOST : "0.0.0.0"
  WEBHOOK_URL : "https://tech-challenge.com.br/v1/payments/webhook"

  DB_RUN_MIGRATIONS : "true"
  DB_NAME : "payment_db"
  DB_PORT : "5432"

  MERCADOPAGO_ACCESS_TOKEN : "APP_USR-1759552174955251-121514-0add78ff150666b46a342875686de7a7-3068451447"
  MERCADOPAGO_COLLECTOR_ID : "3068451447"
  MERCADOPAGO_EXTERNAL_POS_ID : "tccaixafiapf2"
  MERCADOPAGO_API_URL : "https://api.mercadopago.com"

  RABBITMQ_URL : "amqp://guest:guest@rabbitmq:5672/"
  RABBITMQ_EXCHANGE : "amq.topic"
  RABBITMQ_CREATE_PAYMENT_TOPIC : "create.payment"
  RABBITMQ_CREATE_KITCHEN_ORDER_TOPIC : "create.kitchen-order"
  RABBITMQ_ORDER_ERROR_TOPIC : "order-error"
}

container_secrets = {}
health_check_path = "/health"
task_role_policy_arns = [
  "arn:aws:iam::aws:policy/AmazonRDSFullAccess",
  "arn:aws:iam::aws:policy/AmazonSQSFullAccess",
  "arn:aws:iam::aws:policy/SecretsManagerReadWrite"
]

# =======================================================
# Configurações do API Gateway
# =======================================================

apigw_integration_type       = "HTTP_PROXY"
apigw_integration_method     = "ANY"
apigw_payload_format_version = "1.0"
apigw_connection_type        = "VPC_LINK"

authorization_name = "CognitoAuthorizer"

# =======================================================
# Configurações do RDS
# =======================================================
# TODO

# =======================================================
# Configurações do SQS
# =======================================================
# TODO