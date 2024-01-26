package main

import (
	"context"
	"log"
	"os"
	"scraper/discord"
	"scraper/sitea"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	proxyURL            string
	scraperWebhook      string
	scraperSiteABaseURL string
)

func init() {
	proxyURL = os.Getenv("PROXY_URL")
	if proxyURL == "" {
		log.Fatal("Environment variable PROXY_URL must be set")
	}

	scraperWebhook = os.Getenv("SCRAPER_WEBHOOK")
	if scraperWebhook == "" {
		log.Fatal("Environment variable SCRAPER_WEBHOOK must be set")
	}

	scraperSiteABaseURL = os.Getenv("SCRAPER_SITEA_BASEURL")
	if scraperSiteABaseURL == "" {
		log.Fatal("Environment variable SCRAPER_SITEA_BASEURL must be set")
	}

}

func lookForNewJobs() {
	siteAjobs := sitea.ScanNewJobs(scraperSiteABaseURL, proxyURL)
	discord.SendJobsToDiscord(siteAjobs, scraperWebhook)
}

func handler(ctx context.Context) error {
	lookForNewJobs()
	return nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" && os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		// This will run locally
		lookForNewJobs()
	} else {
		lambda.Start(handler)
	}
}
