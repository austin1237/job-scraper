package sited

import (
	"errors"
	"log"
	"net/http"
	"scraper/interest"
	"scraper/job"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func scanSiteD(siteDBaseUrl string) []job.Job {
	var newJobs = []job.Job{}
	url := siteDBaseUrl + "/remote-jobs/developer/"
	response, err := http.Get(url)
	if err != nil {
		log.Println("SiteD: Failed to get site", err)
		return newJobs
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("SiteD: HTTP request failed with status: %s", response.Status)
		return newJobs
	}

	// Parse the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("SiteD: Failed to parse site", err)
		return newJobs
	}

	doc.Find("a.card.m-0.border-left-0.border-right-0.border-top-0.border-bottom").Each(func(i int, s *goquery.Selection) {
		jobURL, exists := s.Attr("href")
		if exists {
			jobTitle := s.Find(".font-weight-bold.larger").Text()
			postTime := strings.TrimSpace(s.Find(".float-right.d-none.d-md-inline.text-secondary small").Text())
			companyInfo := strings.TrimSpace(s.Find("p.m-0.text-secondary").First().Text())
			company := strings.TrimSpace(strings.Split(companyInfo, "|")[0])

			newJob := job.Job{
				Title:   jobTitle,
				Link:    siteDBaseUrl + jobURL,
				Company: company,
			}

			if strings.Contains(postTime, "hours") {
				newJobs = append(newJobs, newJob)
			}
		}
	})

	return newJobs
}

func getSiteDJobInfo(jobUrl string, proxyUrl string) (string, error) {
	response, err := http.Get(proxyUrl + "/proxy?url=" + jobUrl)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("HTTP request failed with status: " + response.Status)
	}

	// Parse the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	jobInfo := doc.Find("div.job_description").Text()
	return jobInfo, nil
}

func ScanNewJobs(siteDBaseUrl string, proxyUrl string) []job.Job {
	jobs := scanSiteD(siteDBaseUrl)
	log.Println("siteD total jobs found", len(jobs))
	interestingJobs := interest.FilterInterest(proxyUrl, jobs, getSiteDJobInfo)
	log.Println("siteD interesting jobs", len(interestingJobs))
	return interestingJobs
}
