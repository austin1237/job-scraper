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
    "SCRAPER_SITEB_BASEURL" = "${var.SCRAPER_SITEB_BASEURL}"
    "SCRAPER_SITEC_BASEURL" = "${var.SCRAPER_SITEC_BASEURL}"
    "SCRAPER_SITED_BASEURL" = "${var.SCRAPER_SITED_BASEURL}"
    "SCRAPER_SITEE_BASEURL" = "${var.SCRAPER_SITEE_BASEURL}"
    "SCRAPER_SITEF_BASEURL" = "${var.SCRAPER_SITEF_BASEURL}"
    "DYNAMO_TABLE" = "${aws_dynamodb_table.job_scraper_company_cache.name}"
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

module "headless_lambda" {
  source         = "./lambda"
  name = "headless-${terraform.workspace}"
  memory_size        = 2048
  timeout = 30
  image_uri        = "${var.AWS_ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/headless@${var.DOCKER_IMAGE_SHA}"
  package_type     = "Image"
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
  api_name = "proxy-${terraform.workspace}"
  openapi = "../openapi.json"
  lambda_arns = [ module.headless_lambda.arn, module.proxy_lambda.arn]
}

# ---------------------------------------------------------------------------------------------------------------------
# DynamoDb
# ---------------------------------------------------------------------------------------------------------------------
resource "aws_dynamodb_table" "job_scraper_company_cache" {
  name           = "job-scraper-company-cache-${terraform.workspace}"
  billing_mode   = "PAY_PER_REQUEST"  # On-demand capacity mode
  hash_key       = "company"

  ttl {
    attribute_name = "ExpirationTime"
    enabled        = true
  }

  attribute {
    name = "company"
    type = "S"  # String data type for company attribute
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# Lambda -> DynamoDb IAM
# ---------------------------------------------------------------------------------------------------------------------
module "dynamodb_lambda_iam" {
  source = "./dynamodb-lambda-iam"
  dynamodb_name = aws_dynamodb_table.job_scraper_company_cache.name
  dynamodb_arn = aws_dynamodb_table.job_scraper_company_cache.arn
  lambda_roles = [module.scraper_lambda.role_name]
}