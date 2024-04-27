package sitec

import (
	"scraper/interest"
	"scraper/job"

	"github.com/PuerkitoBio/goquery"
)

func siteCJobListParser(baseURL string, doc *goquery.Document) []job.Job {
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
					Link:    baseURL + jobURL,
				}
				newJobs = append(newJobs, newJob)
			}
		}
	})

	return newJobs
}

func getSiteCJobInfo(jobUrl string, proxyUrl string) (string, error) {
	doc, err := job.GetJobHtml(jobUrl, proxyUrl)
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
		fullStack := job.GetNewJobs(sitecBaseUrl+"/categories/remote-full-stack-programming-jobs#job-listings", proxyUrl, siteCJobListParser)
		jobChannel <- fullStack
	}()

	go func() {
		backEnd := job.GetNewJobs(sitecBaseUrl+"/categories/remote-back-end-programming-jobs#job-listings", proxyUrl, siteCJobListParser)
		jobChannel <- backEnd
	}()

	for i := 0; i < 2; i++ {
		jobs = append(jobs, <-jobChannel...)
	}

	interestingJobs := interest.FilterInterest(proxyUrl, jobs, getSiteCJobInfo)
	return interestingJobs
}
