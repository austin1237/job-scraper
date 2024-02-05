output "api_url" {
  description = "The URL of the API Gateway"
  value       = aws_apigatewayv2_stage.api_stage.invoke_url
  sensitive   = true
}