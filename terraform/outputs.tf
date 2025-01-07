output "api_gateway_url" {
  description = "Base URL for API Gateway stage"
  value       = "${aws_api_gateway_stage.rvc.invoke_url}/"
}

output "custom_domain_url" {
  description = "Custom domain URL"
  value       = "https://${var.domain_name}/"
}
