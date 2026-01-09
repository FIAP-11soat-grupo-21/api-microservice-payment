module "payment_api" {
  source     = "git::ssh://git@github.com/FIAP-11soat-grupo-21/infra-core.git//modules/ECS-Service?ref=main"
  depends_on = [aws_lb_listener.listener]

  cluster_id            = data.terraform_remote_state.infra.outputs.ecs_cluster_id
  ecs_security_group_id = data.terraform_remote_state.infra.outputs.ecs_security_group_id

  cloudwatch_log_group     = data.terraform_remote_state.infra.outputs.ecs_cloudwatch_log_group
  ecs_container_image      = var.image_name
  ecs_container_name       = var.application_name
  ecs_container_port       = var.image_port
  ecs_service_name         = var.application_name
  ecs_desired_count        = var.desired_count
  registry_credentials_arn = data.terraform_remote_state.infra.outputs.ecr_registry_credentials_arn

  ecs_container_environment_variables = merge(
    var.container_environment_variables
    , {
      DB_HOST : data.terraform_remote_state.infra.outputs.rds_address,
      DB_USERNAME : data.terraform_remote_state.infra.outputs.rds_postgres_db_username
    }
  )

  ecs_container_secrets = merge(
    var.container_secrets
    , {
      DB_PASSWORD : data.terraform_remote_state.infra.outputs.rds_secret_arn
    }
  )

  private_subnet_ids      = data.terraform_remote_state.infra.outputs.private_subnet_id
  task_execution_role_arn = data.terraform_remote_state.infra.outputs.ecs_task_execution_role_arn
  task_role_policy_arns   = var.task_role_policy_arns
  alb_target_group_arn    = aws_alb_target_group.target_group.arn
  alb_security_group_id   = data.terraform_remote_state.infra.outputs.alb_security_group_id

  project_common_tags = data.terraform_remote_state.infra.outputs.project_common_tags
}

module "GetPaymentAPIRoute" {
  source     = "git::ssh://git@github.com/FIAP-11soat-grupo-21/infra-core.git//modules/API-Gateway-Routes?ref=main"
  depends_on = [module.payment_api]

  api_id       = data.terraform_remote_state.infra.outputs.api_gateway_id
  alb_proxy_id = aws_apigatewayv2_integration.alb_proxy.id

  endpoints = {
    get_payment_by_order = {
      route_key           = "GET /payments/order/{orderId}"
      restricted          = false
      auth_integration_id = data.terraform_remote_state.auth.outputs.auth_id
    },
    create_pix_payment = {
      route_key           = "POST /payments/pix"
      restricted          = false
      auth_integration_id = data.terraform_remote_state.auth.outputs.auth_id
    },
    webhook_notification = {
      route_key           = "POST /payments/webhook"
      restricted          = false
      auth_integration_id = data.terraform_remote_state.auth.outputs.auth_id
    },
  }
}