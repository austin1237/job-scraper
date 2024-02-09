package sitea

import (
	"log"
	"scraper/interest"
	"scraper/job"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func siteAJobListParser(baseURL string, doc *goquery.Document) []job.Job {
	newJobs := []job.Job{}
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
			Link:    baseURL + jobLink,
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
			newJobs = append(newJobs, newJob)
		}
	})
	return newJobs
}

func GetSiteAJobInfo(jobLink string, proxyUrl string) (string, error) {
	doc, err := job.GetJobHtml(jobLink, proxyUrl)
	if err != nil {
		return "", err
	}
	jobDescription := doc.Find("div.job-description").Text()

	return jobDescription, nil
}

func ScanNewJobs(siteABaseUrl string, proxyUrl string) []job.Job {
	possibleJobs := []job.Job{}
	finished := false
	page := 1

	for !finished || page > 15 {
		currentJobCount := len(possibleJobs)
		pageStr := strconv.Itoa(page)
		url := siteABaseUrl + "/jobs/remote/nationwide/dev-engineering?page=" + pageStr
		jobs := job.GetNewJobs(url, proxyUrl, siteAJobListParser)
		possibleJobs = append(possibleJobs, jobs...)
		// No new jobs found were done
		if currentJobCount == len(possibleJobs) {
			finished = true
		}
		page++
	}

	log.Println(siteABaseUrl+" total jobs found", len(possibleJobs))
	interestingJobs := interest.FilterInterest(proxyUrl, possibleJobs, GetSiteAJobInfo)
	log.Println(siteABaseUrl+" interesting jobs found", len(interestingJobs))
	return interestingJobs
}
