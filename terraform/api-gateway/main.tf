resource "aws_iam_role" "api_gateway_role" {
  name = "${var.api_name}-api-gateway-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "apigateway.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy" "api_gateway_policy" {
  name = "${var.api_name}-api-gateway-policy"
  role = aws_iam_role.api_gateway_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "lambda:InvokeFunction"
        ]
        Effect   = "Allow"
        Resource = var.lambda_arns
      },
    ]
  })
}

data "aws_caller_identity" "current" {}


# This will be needed to replace the account-id/iam-role in the OpenAPI definition
data "local_file" "openapi" {
  filename = var.openapi
}


resource "aws_apigatewayv2_api" "api" {
  name          = var.api_name
  protocol_type = "HTTP"
  # replaces the placeholders in the OpenAPI definition with the actual values
  body = replace(
  replace(
    data.local_file.openapi.content, 
    "{account-id}", 
    data.aws_caller_identity.current.account_id
  ),
  "{iam-role}",
  aws_iam_role.api_gateway_role.name
  )

}

resource "aws_apigatewayv2_deployment" "api_deployment" {
  api_id      = aws_apigatewayv2_api.api.id
  description = "${var.api_name} deployment"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_cloudwatch_log_group" "api_log_group" {
  name = "${var.api_name}-logs"
}

resource "aws_apigatewayv2_stage" "api_stage" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "prod"
  description = "Production stage"
  auto_deploy   = true
  default_route_settings {
    logging_level = "INFO"
    throttling_burst_limit = 50
    throttling_rate_limit  = 100
  }
  
    access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_log_group.arn
    format          = "$context.identity.sourceIp - - [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId $context.error.message $context.integration.error"
  }
}

