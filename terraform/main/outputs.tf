# DNS público del ALB (frontend y API)
output "alb_dns" {
  value       = aws_lb.main.dns_name
  description = "DNS del Application Load Balancer"
}

# ECS Cluster ARN
output "ecs_cluster_arn" {
  value       = aws_ecs_cluster.main.arn
  description = "ARN del ECS Cluster"
}

# ECS API Service ARN
output "api_service_arn" {
  value       = aws_ecs_service.api_service.arn
  description = "ARN del ECS Service de la API"
}

# ECS Cliente Service ARN
output "client_service_arn" {
  value       = aws_ecs_service.client_service.arn
  description = "ARN del ECS Service del cliente Vue"
}

# Task Definition ARN de API
output "api_task_definition_arn" {
  value       = aws_ecs_task_definition.api_task.arn
  description = "ARN de la task definition de la API"
}

# Task Definition ARN del Cliente
output "client_task_definition_arn" {
  value       = aws_ecs_task_definition.client_task.arn
  description = "ARN de la task definition del cliente"
}

# Target Groups
output "api_target_group_arn" {
  value       = aws_lb_target_group.api_tg.arn
  description = "ARN del target group de la API"
}

output "client_target_group_arn" {
  value       = aws_lb_target_group.app_tg.arn
  description = "ARN del target group del cliente"
}
