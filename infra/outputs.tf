output "payment_queue_url" {
  description = "URL do Payment API no ALB"
  value       = data.terraform_remote_state.infra.outputs.sqs_payments_queue_url 
}

output "kitchen_order_queue_url" {
  description = "URL da fila de pedidos da cozinha"
  value       = data.terraform_remote_state.infra.outputs.sqs_kitchen_orders_queue_url
}

output "payment_order_error_queue_url" {
  description = "URL da fila de erros de pedidos"
  value       = data.terraform_remote_state.infra.outputs.sqs_payments_order_error_queue_url
}

output "ecs_service_id" {
  description = "ID do serviço ECS do Payment API"
  value       = module.payment_api.service_id
}