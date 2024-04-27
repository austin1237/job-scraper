package main

import (
	"context"
	"log"
	"os"
	"scraper/job"
	"scraper/remotive"
	"scraper/sitea"
	"scraper/siteb"
	"scraper/sitec"
	"scraper/sited"
	"scraper/sitee"
	"scraper/sitef"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Site struct {
	ScanNewJobs func(string, string) []job.Job
	BaseURL     string
}

type Result struct {
	Elapsed time.Duration
	URL     string
}

var (
	proxyURL            string
	scraperSiteABaseURL string
	scraperSiteBBaseURL string
	scraperSiteCBaseURL string
	scraperSiteDBaseURL string
	scraperSiteEBaseURL string
	scraperSiteFBaseURL string
	jobURL              string
)

func init() {
	proxyURL = os.Getenv("PROXY_URL")
	if proxyURL == "" {
		log.Fatal("Environment variable PROXY_URL must be set")
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

	jobURL = os.Getenv("JOB_URL")
	if jobURL == "" {
		log.Fatal("Environment variable JOB_URL must be set")
	}

}

func lookForNewJobs() {
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

	results := make([]Result, 0, len(sites))
	doneChannel := make(chan Result, len(sites))
	for _, site := range sites {
		go func(site Site) {
			start := time.Now()
			interestingJobs := site.ScanNewJobs(site.BaseURL, proxyURL)
			results, err := job.SendJobs(jobURL, interestingJobs)
			elapsed := time.Since(start)
			if err != nil {
				log.Println("Error sending to job api", err)
				doneChannel <- Result{Elapsed: elapsed, URL: site.BaseURL}
				return
			}
			log.Println(site.BaseURL+", Total Jobs: ", results.Total, "Uncached Jobs: ", results.Uncached, "Duplicates: ", results.Duplicates)
			doneChannel <- Result{Elapsed: elapsed, URL: site.BaseURL}
		}(site)
	}

	// Wait for all goroutines to finish
	for range sites {
		result := <-doneChannel
		results = append(results, result)
	}

	for _, result := range results {
		log.Printf("Execution took %s for %s \n", result.Elapsed, result.URL)
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
