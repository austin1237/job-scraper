# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

variable "SCRAPER_WEBHOOK" {
    sensitive = true
}

variable "SCRAPER_SITEA_BASEURL" {
    sensitive = true
}

variable "SCRAPER_SITEB_BASEURL" {
    sensitive = true
}

variable "SCRAPER_SITEC_BASEURL" {
    sensitive = true
}

variable "SCRAPER_SITED_BASEURL" {
    sensitive = true
}

variable "SCRAPER_SITEE_BASEURL" {
    sensitive = true
}

variable "AWS_ACCOUNT_ID" {
    sensitive = true
}

variable "DOCKER_IMAGE_SHA" {
    sensitive = true
}

