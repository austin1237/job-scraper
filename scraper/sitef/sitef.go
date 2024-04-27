package sitef

import (
	"scraper/interest"
	"scraper/job"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func siteFJobListParser(baseURL string, doc *goquery.Document) []job.Job {
	newJobs := []job.Job{}
	doc.Find("div.job-wrapper").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("h4 a").Attr("href")
		company := strings.TrimSpace(s.Find("div.company a").First().Text())
		timePosted := strings.TrimSpace(s.Find("div.date").Last().Text())
		if strings.Contains(timePosted, "hour") || strings.Contains(timePosted, "minute") {
			newJob := job.Job{
				Company: company,
				Link:    baseURL + link,
			}
			newJobs = append(newJobs, newJob)
		}
	})

	return newJobs
}

func getSiteFJobInfo(jobUrl string, proxyUrl string) (string, error) {
	doc, err := job.GetJobHtml(jobUrl, proxyUrl)
	if err != nil {
		return "", err
	}
	text := ""
	doc.Find("div.job").Each(func(i int, s *goquery.Selection) {
		text += s.Text() + " "
	})
	return text, nil
}

func ScanNewJobs(baseURL string, proxyURL string) []job.Job {
	subUrl := "/jobs?category=development&location=north-america&positionType=full-time"
	jobs := job.GetNewJobs(baseURL+subUrl, proxyURL, siteFJobListParser, "headless")
	interestingJobs := interest.FilterInterest(proxyURL, jobs, getSiteFJobInfo)
	return interestingJobs
}
