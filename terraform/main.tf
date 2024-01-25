terraform {
  backend "s3" {
    bucket         = "job-scraper-state"
    key            = "global/s3/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "job-scraper-state-lock"
    encrypt        = true
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.33"
    }
  }
  required_version = "~> 1.7"
}

# ---------------------------------------------------------------------------------------------------------------------
# Lambdas
# ---------------------------------------------------------------------------------------------------------------------

module "scraper_lambda" {
  source         = "./lambda"
  zip_location   = "../scraper/bootstrap.zip"
  name           = "job-scraper-${terraform.workspace}"
  handler        = "bootstrap"
  run_time       = "provided.al2"
  timeout        = 300
  env_vars = {
    "PROXY_URL" = "${module.proxy_gateway.api_url}"
    "SCRAPER_WEBHOOK" = "${var.SCRAPER_WEBHOOK}"
    "SCRAPER_SITEA_BASEURL" = "${var.SCRAPER_SITEA_BASEURL}"
  } 
}

module "proxy_lambda" {
  source         = "./lambda"
  zip_location   = "../proxy/bootstrap.zip"
  name           = "proxy-${terraform.workspace}"
  handler        = "bootstrap"
  run_time       = "provided.al2"
  timeout        = 300
  env_vars = {} 
}

# ---------------------------------------------------------------------------------------------------------------------
# Cloudwatch that will trigger the scraper lambda
# ---------------------------------------------------------------------------------------------------------------------
module "scraper_lambda_trigger" {
  source               = "./cloudwatch-lambda-trigger"
  # Every Weekday at 5pm MDT
  start_time           = "cron(0 0 * * ? *)"
  name                 = "scraper-lambda-trigger-${terraform.workspace}"
  lambda_function_name = "${module.scraper_lambda.name}"
  description          = "The timed trigger for ${module.scraper_lambda.name}"
  lambda_arn           = "${module.scraper_lambda.arn}"
}

# ---------------------------------------------------------------------------------------------------------------------
# API gateway
# ---------------------------------------------------------------------------------------------------------------------

module "proxy_gateway" {
  source         = "./api-gateway"
  lambda_name =  "${module.proxy_lambda.name}"
  lambda_invoke_arn =   "${module.proxy_lambda.invoke_arn}"
  api_name = "proxy-${terraform.workspace}"
  route = "proxy"
}