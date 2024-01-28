package siteb

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"scraper/interest"
	"strings"

	"scraper/job"

	"github.com/PuerkitoBio/goquery"
)

func scanSiteB(siteBBaseUrl string) []job.Job {
	url := siteBBaseUrl + "/jobs"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("HTTP request failed with status: %s", response.Status)
	}

	// Parse the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var newJobs = []job.Job{}

	// Find the div with class "row search-result"
	doc.Find("div.row.search-result").Each(func(i int, s *goquery.Selection) {
		// Extract the href attribute from the <a> element with rel="canonical"
		newJob := job.Job{}
		jobLink, exists := s.Find("a[rel='canonical']").Attr("href")
		if exists {
			newJob.Link = siteBBaseUrl + jobLink
		}

		// Extract the company name from the <li> element with the <i> child having class "fa fa-building"
		companyName := s.Find("li").Has("i.fa.fa-building").Text()
		newJob.Company = companyName

		// Extract the job title from the <h2> element with class "jobl-title"
		jobTitle := s.Find("h2.jobl-title").Text()
		newJob.Title = jobTitle

		// Extract the posted time from the <li> element with the <i> child having class "fa fa-calendar"
		postedTime := s.Find("li").Has("i.fa.fa-calendar").Text()
		if strings.Contains(postedTime, "hours") {
			newJobs = append(newJobs, newJob)
		}
	})

	return newJobs
}

func getSiteBJobInfo(jobUrl string, proxyUrl string) (string, error) {
	response, err := http.Get(proxyUrl + "/proxy?url=" + jobUrl)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := errors.New("source HTTP request failed with status: " + response.Status)
		return "", err
	}

	// Parse the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	div := doc.Find(".job-details")
	description := ""
	// Find and print all <p> elements
	div.Find("div div p").Each(func(i int, s *goquery.Selection) {
		description += " " + s.Text()
	})

	// Find and print all <ul> elements
	div.Find("div div ul").Each(func(i int, s *goquery.Selection) {
		description += " " + s.Text()
	})
	return description, nil
}

func ScanNewJobs(sitebBaseUrl string, proxyUrl string) []job.Job {
	jobs := scanSiteB(sitebBaseUrl)
	fmt.Println("siteB total jobs found", len(jobs))
	interestingJobs := interest.FilterInterest(proxyUrl, jobs, getSiteBJobInfo)
	fmt.Println("siteB interesting jobs", len(interestingJobs))
	return interestingJobs
}
