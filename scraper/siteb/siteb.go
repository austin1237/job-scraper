package siteb

import (
	"scraper/interest"
	"strings"

	"scraper/job"

	"github.com/PuerkitoBio/goquery"
)

func siteBJobListParser(siteBBaseUrl string, doc *goquery.Document) []job.Job {
	newJobs := []job.Job{}
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
	doc, err := job.GetJobHtml(jobUrl, proxyUrl)
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
	jobs := job.GetNewJobs(sitebBaseUrl+"/jobs", proxyUrl, siteBJobListParser)
	interestingJobs := interest.FilterInterest(proxyUrl, jobs, getSiteBJobInfo)
	return interestingJobs
}
