output "invoke_arn" {
  value = var.package_type == "Zip" ? aws_lambda_function.lambda_zip[0].invoke_arn : aws_lambda_function.lambda_image[0].invoke_arn
}

output "name" {
  value = var.package_type == "Zip" ? aws_lambda_function.lambda_zip[0].function_name : aws_lambda_function.lambda_image[0].function_name
}

output "arn" {
  value = var.package_type == "Zip" ? aws_lambda_function.lambda_zip[0].arn : aws_lambda_function.lambda_image[0].arn
}

