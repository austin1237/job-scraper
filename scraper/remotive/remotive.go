package remotive

import (
	"encoding/json"
	"log"
	"net/http"
	"scraper/cache"
	"scraper/interest"
	"scraper/job"
	"time"
)

type remotiveJob struct {
	ID                        int    `json:"id"`
	URL                       string `json:"url"`
	Title                     string `json:"title"`
	CompanyName               string `json:"company_name"`
	CompanyLogo               string `json:"company_logo"`
	Category                  string `json:"category"`
	JobType                   string `json:"job_type"`
	PublicationDate           string `json:"publication_date"`
	CandidateRequiredLocation string `json:"candidate_required_location"`
	Salary                    string `json:"salary"`
	Description               string `json:"description"`
}

type JobsResponse struct {
	JobCount int           `json:"job-count"`
	Jobs     []remotiveJob `json:"jobs"`
}

func callApi(site string) []remotiveJob {
	newJobs := []remotiveJob{}
	resp, err := http.Get(site)

	if err != nil {
		log.Println("error calling remotive api", err)
		return newJobs
	}
	defer resp.Body.Close()
	var jobsResponse JobsResponse
	err = json.NewDecoder(resp.Body).Decode(&jobsResponse)

	if err != nil {
		log.Println("error decoding response", err)
		return newJobs
	}
	yesterday := time.Now().Add(-24 * time.Hour)

	for _, newJob := range jobsResponse.Jobs {
		recent := false
		locationMatch := false
		pubDate, err := time.Parse("2006-01-02T15:04:05", newJob.PublicationDate)
		if err != nil {
			log.Println("error parsing date", err)
			continue
		}

		if pubDate.After(yesterday) {
			recent = true
		}

		if newJob.CandidateRequiredLocation == "USA" || newJob.CandidateRequiredLocation == "Worldwide" || newJob.CandidateRequiredLocation == "" {
			locationMatch = true
		}

		if recent && locationMatch {
			newJobs = append(newJobs, newJob)
		}

	}
	return newJobs
}

func ScanNewJobs(baseURL string, proxyUrl string, cache *cache.Cache) ([]job.Job, []job.Job) {
	remotiveJobs := callApi(baseURL + "/api/remote-jobs?category=software-dev&limit=100")
	log.Println("Remotive total jobs found", len(remotiveJobs))
	var interestingJobs []job.Job
	var newJobs []job.Job

	for _, newJob := range remotiveJobs {
		newJobs = append(newJobs, job.Job{
			Title:   newJob.Title,
			Link:    newJob.URL,
			Company: newJob.CompanyName,
		})
	}

	unCachedJobs, err := cache.FilterCachedCompanies(newJobs)
	if err != nil {
		log.Println("Error filtering cached companies", err)
	}
	log.Println(baseURL+" total jobs not found in cache", len(unCachedJobs))

	for _, newJob := range remotiveJobs {
		if interest.CheckIfInterested(newJob.Description) {
			interestingJobs = append(interestingJobs, job.Job{
				Title:   newJob.Title,
				Link:    newJob.URL,
				Company: newJob.CompanyName,
			})
		}
	}
	log.Println("Remotive interesting jobs", len(interestingJobs))
	return unCachedJobs, interestingJobs
}
