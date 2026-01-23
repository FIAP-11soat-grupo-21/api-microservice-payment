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

  MERCADOPAGO_API_URL : "https://api.mercadopago.com"
  MERCADOPAGO_COLLECTOR_ID : "MERCADOPAGO_COLLECTOR_ID_PLACEHOLDER"
  MERCADOPAGO_EXTERNAL_POS_ID : "MERCADOPAGO_EXTERNAL_POS_ID_PLACEHOLDER"
  MERCADOPAGO_ACCESS_TOKEN : "MERCADOPAGO_ACCESS_TOKEN_PLACEHOLDER"
}

container_secrets = {}
health_check_path = "/health"
task_role_policy_arns = [
  "arn:aws:iam::aws:policy/AmazonRDSFullAccess",
  "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
]

# =======================================================
# Configurações do API Gateway
# =======================================================

apigw_integration_type       = "HTTP_PROXY"
apigw_integration_method     = "ANY"
apigw_payload_format_version = "1.0"
apigw_connection_type        = "VPC_LINK"

authorization_name = "CognitoAuthorizer"
