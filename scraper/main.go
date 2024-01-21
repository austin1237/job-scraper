package main

import (
	"log"
	"os"
	"scraper/discord"
	"scraper/interest"
	"scraper/scanner"
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

func main() {
	links := scanner.ScanNewJobs(scraperSiteABaseURL)
	interestingJobs := interest.FilterInterest(scraperSiteABaseURL, proxyURL, links)
	discord.SendJobsToDiscord(interestingJobs, scraperWebhook)

}
