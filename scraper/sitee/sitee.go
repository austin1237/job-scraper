package sitee

import (
	"log"
	"scraper/interest"
	"scraper/job"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func siteEJobListParser(baseURL string, doc *goquery.Document) []job.Job {
	newJobs := []job.Job{}
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		recent := false
		link, _ := s.Find("a").First().Attr("href")
		company := strings.TrimSpace(s.Find("a").Eq(1).Text())
		s.Find("span").Each(func(i int, span *goquery.Selection) {
			text := strings.ToLower(span.Text())
			if strings.Contains(text, "new job") {
				recent = true
			}
		})
		if recent {
			newJob := job.Job{
				Company: company,
				Link:    baseURL + link,
			}
			newJobs = append(newJobs, newJob)
		}
	})

	return newJobs

}

func getSiteEJobInfo(jobUrl string, proxyUrl string) (string, error) {
	doc, err := job.GetJobHtml(jobUrl, proxyUrl)
	if err != nil {
		return "", err
	}
	jobInfo := ""
	doc.Find("div.mb-6.prose.break-words.prose-md.max-w-none").Each(func(i int, s *goquery.Selection) {
		jobInfo += s.Find("*").Text() + " "
	})
	return jobInfo, nil
}

func ScanNewJobs(baseURL string, proxyURL string) []job.Job {
	jobs := job.GetNewJobs(baseURL+"/category/development", proxyURL, siteEJobListParser)
	log.Println(baseURL+" total jobs found", len(jobs))
	interestingJobs := interest.FilterInterest(proxyURL, jobs, getSiteEJobInfo)
	log.Println(baseURL+" total interesting jobs found", len(interestingJobs))
	return interestingJobs
}
