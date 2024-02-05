variable "api_name" {
    description = "arn of the lambda to attach to the gateway"
}

variable "openapi" {
    description = "The OpenAPI definition"
}

variable "lambda_arns" {
  description = "Array of Lambda ARNs"
  type        = list(string)
  default     = []
}

