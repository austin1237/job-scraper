variable "dynamodb_name" {
  description = "The name of the DynamoDB table"
  type        = string
}

variable "dynamodb_arn" {
  description = "The ARN of the DynamoDB table"
  type        = string
}

variable "lambda_roles" {
  description = "List of IAM role names for Lambda functions"
  type        = list(string)
}