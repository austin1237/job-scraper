output "invoke_arn" {
  value = var.package_type == "Zip" ? aws_lambda_function.lambda_zip[0].invoke_arn : aws_lambda_function.lambda_image[0].invoke_arn
}

output "qualified_arn" {
  value = var.package_type == "Zip" ? aws_lambda_function.lambda_zip[0].qualified_arn : aws_lambda_function.lambda_image[0].qualified_arn
}

output "version" {
  value = var.package_type == "Zip" ? aws_lambda_function.lambda_zip[0].version : aws_lambda_function.lambda_image[0].version
}

