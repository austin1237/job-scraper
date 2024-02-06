package sitec

import (
	"errors"
	"log"
	"net/http"

	"scraper/interest"
	"scraper/job"

	"github.com/PuerkitoBio/goquery"
)

func scanSiteC(siteCBaseUrl string, suburl string) []job.Job {
	url := siteCBaseUrl + suburl
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

	doc.Find("section.jobs ul li").Each(func(i int, s *goquery.Selection) {
		if s.Find("span.new").Length() > 0 {
			jobTitle := s.Find("span.title").Text()
			companyName := s.Find("span.company").First().Text()
			jobURL, exists := s.Find("a:has(span.title)").Attr("href")
			if exists {
				newJob := job.Job{
					Title:   jobTitle,
					Company: companyName,
					Link:    siteCBaseUrl + jobURL,
				}
				newJobs = append(newJobs, newJob)
			}
		}
	})

	return newJobs
}

func getSiteCJobInfo(jobUrl string, proxyUrl string) (string, error) {
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

	description := ""
	// Find and print all <p> elements
	doc.Find("section.job").Each(func(i int, s *goquery.Selection) {
		s.Find("*").Each(func(i int, s *goquery.Selection) {
			description += s.Text() + "\n"
		})
	})
	return description, nil
}

func ScanNewJobs(sitecBaseUrl string, proxyUrl string) []job.Job {
	var jobs = []job.Job{}
	jobChannel := make(chan []job.Job, 2)

	go func() {
		fullStack := scanSiteC(sitecBaseUrl, "/categories/remote-full-stack-programming-jobs#job-listings")
		jobChannel <- fullStack
	}()

	go func() {
		backEnd := scanSiteC(sitecBaseUrl, "/categories/remote-back-end-programming-jobs#job-listings")
		jobChannel <- backEnd
	}()

	for i := 0; i < 2; i++ {
		jobs = append(jobs, <-jobChannel...)
	}

	jobs = job.DeduplicatedLinks(jobs)
	log.Println("siteC total jobs found", len(jobs))
	interestingJobs := interest.FilterInterest(proxyUrl, jobs, getSiteCJobInfo)
	log.Println("siteC interesting jobs", len(interestingJobs))
	return interestingJobs
}
