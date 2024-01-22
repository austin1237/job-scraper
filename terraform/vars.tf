# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

variable "PROXY_URL" {
    sensitive = true
}

variable "SCRAPER_WEBHOOK" {
    sensitive = true
}

variable "SCRAPER_SITEA_BASEURL" {
    sensitive = true
}

