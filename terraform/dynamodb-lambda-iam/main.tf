resource "aws_iam_policy" "lambda_dynamodb_policy" {
  name = "LambdaDynamoDBPolicy-${terraform.workspace}-${var.dynamodb_name}-${terraform.workspace}"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow",
        Action   = [
          "dynamodb:GetItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:PutItem"
        ],
        Resource = [
          "${var.dynamodb_arn}",
        ],
      },
    ],
  })
}

# Loop through the lambda functions and attach the DynamoDB policy to each
resource "aws_iam_role_policy_attachment" "lambda_dynamodb_policy_attachment" {
  for_each   = toset(var.lambda_roles)
  policy_arn = aws_iam_policy.lambda_dynamodb_policy.arn
  role       = each.value
}