resource "aws_iam_policy" "lambda_logs_policy" {
  name        = "lambda-cloudwatch-logs-policy-${var.name}"
  description = "Allows Lambda function to write logs to CloudWatch Logs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
EOF
}

resource "aws_iam_role" "lambda_role" {
  name = "my-lambda-role-${terraform.workspace}-${var.name}"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs_attachment" {
  policy_arn = aws_iam_policy.lambda_logs_policy.arn
  role       = aws_iam_role.lambda_role.name
}

resource "aws_lambda_function" "lambda_zip" {
  count = var.package_type == "Zip" ? 1 : 0

  function_name = var.name
  role          = aws_iam_role.lambda_role.arn
  handler       = var.handler
  timeout       = var.timeout
  memory_size   = var.memory_size
  package_type  = var.package_type

  filename         = var.zip_location
  source_code_hash = filesha256(var.zip_location)
  runtime          = var.run_time

  environment {
    variables = var.env_vars
  }
}

resource "aws_lambda_function" "lambda_image" {
  count = var.package_type == "Image" ? 1 : 0

  function_name = var.name
  role          = aws_iam_role.lambda_role.arn
  timeout       = var.timeout
  memory_size   = var.memory_size
  package_type  = var.package_type

  image_uri = var.image_uri

  environment {
    variables = var.env_vars
  }
}
