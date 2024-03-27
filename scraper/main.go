package main

import (
	"context"
	"log"
	"os"
	"scraper/cache"
	"scraper/discord"
	"scraper/dynamo"
	"scraper/job"
	"scraper/remotive"
	"scraper/sitea"
	"scraper/siteb"
	"scraper/sitec"
	"scraper/sited"
	"scraper/sitee"
	"scraper/sitef"

	"github.com/aws/aws-lambda-go/lambda"
)

type Site struct {
	ScanNewJobs func(string, string, *cache.Cache) ([]job.Job, []job.Job)
	BaseURL     string
}

var (
	proxyURL            string
	scraperWebhook      string
	scraperSiteABaseURL string
	scraperSiteBBaseURL string
	scraperSiteCBaseURL string
	scraperSiteDBaseURL string
	scraperSiteEBaseURL string
	scraperSiteFBaseURL string
	dynamoTable         string
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

	scraperSiteCBaseURL = os.Getenv("SCRAPER_SITEC_BASEURL")
	if scraperSiteCBaseURL == "" {
		log.Fatal("Environment variable SCRAPER_SITEC_BASEURL must be set")
	}

	scraperSiteDBaseURL = os.Getenv("SCRAPER_SITED_BASEURL")
	if scraperSiteDBaseURL == "" {
		log.Fatal("Environment variable SCRAPER_SITED_BASEURL must be set")
	}

	scraperSiteEBaseURL = os.Getenv("SCRAPER_SITEE_BASEURL")
	if scraperSiteEBaseURL == "" {
		log.Fatal("Environment variable SCRAPER_SITEE_BASEURL must be set")
	}

	scraperSiteFBaseURL = os.Getenv("SCRAPER_SITEF_BASEURL")
	if scraperSiteFBaseURL == "" {
		log.Fatal("Environment variable SCRAPER_SITEF_BASEURL must be set")
	}

	dynamoTable = os.Getenv("DYNAMO_TABLE")
	if dynamoTable == "" {
		log.Fatal("Environment variable DYNAMO_TABLE must be set")
	}

}

func lookForNewJobs() {
	table, err := dynamo.NewTable(dynamoTable, "us-east-1") // replace with your table name
	if err != nil {
		log.Fatal(err)
	}

	cache := cache.NewCache(table)
	var sites = []Site{
		{ScanNewJobs: sitea.ScanNewJobs, BaseURL: scraperSiteABaseURL},
		{ScanNewJobs: siteb.ScanNewJobs, BaseURL: scraperSiteBBaseURL},
		{ScanNewJobs: sitec.ScanNewJobs, BaseURL: scraperSiteCBaseURL},
		{ScanNewJobs: sited.ScanNewJobs, BaseURL: scraperSiteDBaseURL},
		{ScanNewJobs: sitee.ScanNewJobs, BaseURL: scraperSiteEBaseURL},
		{ScanNewJobs: sitef.ScanNewJobs, BaseURL: scraperSiteFBaseURL},
		{ScanNewJobs: remotive.ScanNewJobs, BaseURL: "https://remotive.com"},
		// Add more sites here
	}

	doneChannel := make(chan bool, len(sites))
	for _, site := range sites {
		go func(site Site) {
			uncachedJobs, interestingJobs := site.ScanNewJobs(site.BaseURL, proxyURL, cache)
			errs := discord.SendJobsToDiscord(interestingJobs, scraperWebhook)
			if len(errs) == 0 {
				cache.WriteCompaniesToCache(uncachedJobs)
			} else {
				log.Println("Error sending to discord", errs)
			}
			doneChannel <- true
		}(site)
	}

	// Wait for all goroutines to finish
	for range sites {
		<-doneChannel
	}

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
