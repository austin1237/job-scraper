package sitea

import (
	"log"
	"scraper/interest"
	"scraper/job"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func siteAJobListParser(baseURL string, doc *goquery.Document) []job.Job {
	newJobs := []job.Job{}
	doc.Find("div[id^='job-card-']").Each(func(i int, s *goquery.Selection) {
		recent := false
		titleCheck := false
		companyLink, _ := s.Find("a[href^='/company/']").Attr("href")
		jobLink, _ := s.Find("a[class='card-alias-after-overlay hover-underline link-visited-color text-break']").Attr("href")
		jobTitle := s.Find("a[class='card-alias-after-overlay hover-underline link-visited-color text-break']").Text()
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
	jobDescription := doc.Find("div[id^='job-post-body-']").Text()

	return jobDescription, nil
}

func ScanNewJobs(siteABaseUrl string, proxyUrl string) []job.Job {
	var wg sync.WaitGroup
	jobsChan := make(chan []job.Job)

	fetchJobs := func(url string) {
		defer wg.Done()
		finished := false
		page := 1
		for !finished && page <= 15 {
			pageStr := strconv.Itoa(page)
			url := url + "?page=" + pageStr
			jobs := job.GetNewJobs(url, proxyUrl, siteAJobListParser)
			jobsChan <- jobs
			// No new jobs found were done
			if len(jobs) == 0 {
				finished = true
			}
			page++
		}
	}

	wg.Add(2)
	go fetchJobs(siteABaseUrl + "/jobs/remote/nationwide/dev-engineering")
	// non fully remote Denver, CO from the last 24 hours
	go fetchJobs(siteABaseUrl + "/jobs/hybrid/office/dev-engineering?search=Software+Engineer&daysSinceUpdated=1&location=Denver-CO-USA&longitude=-104.98485&latitude=39.73845&searcharea=25mi")

	go func() {
		wg.Wait()
		close(jobsChan)
	}()

	possibleJobs := []job.Job{}
	for jobs := range jobsChan {
		possibleJobs = append(possibleJobs, jobs...)
	}
	log.Println(siteABaseUrl+", Total jobs found before interest check: ", len(possibleJobs))

	interestingJobs := interest.FilterInterest(proxyUrl, possibleJobs, GetSiteAJobInfo)
	return interestingJobs
}
