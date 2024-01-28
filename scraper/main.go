package main

import (
	"context"
	"log"
	"os"
	"scraper/discord"
	"scraper/sitea"
	"scraper/siteb"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	proxyURL            string
	scraperWebhook      string
	scraperSiteABaseURL string
	scraperSiteBBaseURL string
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

	scraperSiteBBaseURL = os.Getenv("SCRAPER_SITEB_BASEURL")
	if scraperSiteBBaseURL == "" {
		log.Fatal("Environment variable SCRAPER_SITEB_BASEURL must be set")
	}

}

func lookForNewJobs() {
	doneChannel := make(chan bool)

	go func() {
		siteAjobs := sitea.ScanNewJobs(scraperSiteABaseURL, proxyURL)
		discord.SendJobsToDiscord(siteAjobs, scraperWebhook)
		doneChannel <- true
	}()

	go func() {
		siteBJobs := siteb.ScanNewJobs(scraperSiteBBaseURL, proxyURL)
		discord.SendJobsToDiscord(siteBJobs, scraperWebhook)
		doneChannel <- true
	}()

	// Wait for both goroutines to finish
	<-doneChannel
	<-doneChannel
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
