package sited

import (
	"log"
	"scraper/cache"
	"scraper/interest"
	"scraper/job"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func siteDJobListParser(siteDaseUrl string, doc *goquery.Document) []job.Job {
	newJobs := []job.Job{}
	doc.Find("a.card.m-0.border-left-0.border-right-0.border-top-0.border-bottom").Each(func(i int, s *goquery.Selection) {
		jobURL, exists := s.Attr("href")
		if exists {
			jobTitle := s.Find(".font-weight-bold.larger").Text()
			postTime := strings.TrimSpace(s.Find(".float-right.d-none.d-md-inline.text-secondary small").Text())
			companyInfo := strings.TrimSpace(s.Find("p.m-0.text-secondary").First().Text())
			company := strings.TrimSpace(strings.Split(companyInfo, "|")[0])

			newJob := job.Job{
				Title:   jobTitle,
				Link:    siteDaseUrl + jobURL,
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
	doc, err := job.GetJobHtml(jobUrl, proxyUrl)
	if err != nil {
		return "", err
	}
	jobInfo := doc.Find("div.job_description").Text()
	return jobInfo, nil
}

func ScanNewJobs(siteDBaseUrl string, proxyUrl string, cache *cache.Cache) ([]job.Job, []job.Job) {
	jobs := job.GetNewJobs(siteDBaseUrl+"/remote-jobs/developer/", proxyUrl, siteDJobListParser)
	log.Println(siteDBaseUrl+" total jobs found", len(jobs))
	unCachedJobs, err := cache.FilterCachedCompanies(jobs)
	if err != nil {
		log.Println("Error filtering cached companies", err)
	}
	log.Println(siteDBaseUrl+" total jobs not found in cache", len(unCachedJobs))
	interestingJobs := interest.FilterInterest(proxyUrl, unCachedJobs, getSiteDJobInfo)
	log.Println(siteDBaseUrl+" interesting jobs", len(interestingJobs))
	return unCachedJobs, interestingJobs
}
