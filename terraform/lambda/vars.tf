variable "zip_location" {
  description = "path to the ziped lambda"
  default = ""
}

variable "name" {
  description = "The name of the lambda function"
}

variable "handler" {
  description = "name of the lambdas handler"
  default = ""
}

variable "run_time" {
  description = "run time of the lambda"
  default = ""
}

variable "env_vars" {
  type        = map(string)
  description = "environment variables for the lambda"
}

variable "memory_size" {
  description = "Amount of memory in MB your Lambda Function can use at runtime. CPU is implicitly tied to this."
  default     = 128
}

variable "timeout" {
  description = "The max number of seconds the lambda can run"
  default     = 3
}

variable "package_type" {
  description = "The package type of the lambda only valid for docker based lambdas"
  default = "Zip"
}

variable "image_uri" {
  description = "The docker image uri of the lambda"
  default = ""
}
