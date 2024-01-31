package sitea

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"scraper/interest"
	"scraper/job"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func scanSiteA(siteABaseUrl string) []job.Job {
	possibleJobs := []job.Job{}
	finished := false
	page := 1

	for !finished || page > 15 {
		pageStr := strconv.Itoa(page)
		url := siteABaseUrl + "/jobs/remote/nationwide/dev-engineering?page=" + pageStr
		// Make an HTTP GET request
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
		currentJobCount := len(possibleJobs)

		doc.Find("div[id^='job-card-']").Each(func(i int, s *goquery.Selection) {
			recent := false
			titleCheck := false
			companyLink, _ := s.Find("a[href^='/company/']").Attr("href")
			jobLink, _ := s.Find("a[id='job-card-alias']").Attr("href")
			jobTitle := s.Find("a[id='job-card-alias']").Text()
			timePosted := s.Find("span.font-barlow.text-gray-03").Text()

			// Split the companyLink on '/' and get the last part
			parts := strings.Split(companyLink, "/")
			companyName := parts[len(parts)-1]

			newJob := job.Job{
				Link:    siteABaseUrl + jobLink,
				Company: companyName,
				Title:   jobTitle,
			}

			timePosted = strings.ToLower(timePosted)
			jobTitle = strings.ToLower(jobTitle)

			if strings.Contains(timePosted, "hours ago") || strings.Contains(timePosted, "minutes ago") || strings.Contains(timePosted, "hour ago") {
				recent = true
			}

			titles := []string{"software engineer", "developer", "backend engineer", "backend developer", "backend", "software developer"}

			for _, title := range titles {
				if strings.Contains(jobTitle, title) {
					titleCheck = true
					break
				}
			}

			if recent && titleCheck {
				possibleJobs = append(possibleJobs, newJob)
			}
		})
		// No new jobs found were done
		if currentJobCount == len(possibleJobs) {
			finished = true
		}

		page++
	}

	return job.DeduplicatedLinks(possibleJobs)
}

func GetSiteAJobInfo(jobLink string, proxyUrl string) (string, error) {
	response, err := http.Get(proxyUrl + "/proxy?url=" + jobLink)
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
	jobDescription := doc.Find("div.job-description").Text()

	return jobDescription, nil
}

func ScanNewJobs(siteABaseUrl string, proxyUrl string) []job.Job {
	possibleJobs := scanSiteA(siteABaseUrl)
	fmt.Println("siteA total jobs found", len(possibleJobs))
	interestingJobs := interest.FilterInterest(proxyUrl, possibleJobs, GetSiteAJobInfo)
	fmt.Println("siteA interesting jobs found", len(interestingJobs))
	return interestingJobs
}
