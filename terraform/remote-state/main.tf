terraform {
  backend "local" {

  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.33"
    }
  }
  required_version = "~> 1.7"
}

provider "aws" {
  region  = "us-east-1"
}

resource "aws_s3_bucket" "terraform_state" {
  bucket = "job-scraper-state"

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_s3_bucket_versioning" "example" {
  bucket = aws_s3_bucket.terraform_state.id

  versioning_configuration {
    status = "Enabled"
  }
}


resource "aws_dynamodb_table" "terraform_state_lock" {
  name           = "job-scraper-state-lock"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}