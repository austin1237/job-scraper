package interest

import (
	"fmt"
	"regexp"
	"scraper/job"
	"strings"
	"sync"
)

func CheckIfInterested(description string) bool {
	keywords := []string{"typescript", "node", "nodejs", "node.js", "go", "golang"}
	// Check if keywords are present in the job's text
	descriptionToLower := strings.ToLower(description)
	for _, keyword := range keywords {
		pattern := "\\b" + keyword + "\\b"
		match, _ := regexp.MatchString(pattern, descriptionToLower)
		if match {
			return true
		}
	}
	return false
}

type JobInfoGetter func(string, string) (string, error)

func FilterInterest(proxyUrl string, possibleJobs []job.Job, jobInfoGetter JobInfoGetter) []job.Job {
	interestingJobs := []job.Job{}
	possibleJobs = job.DeduplicatedLinks(possibleJobs)
	var wg sync.WaitGroup

	// Number of concurrent goroutines
	maxGoroutines := 10
	var goroutineCount int

	for _, possibleJob := range possibleJobs {
		wg.Add(1)
		goroutineCount++

		go func(possibleJob job.Job) {
			defer wg.Done()
			description, err := jobInfoGetter(possibleJob.Link, proxyUrl)
			if err != nil {
				fmt.Println(err)
			}
			if CheckIfInterested(description) {
				interestingJobs = append(interestingJobs, possibleJob)
			}
			goroutineCount--
		}(possibleJob)

		// Limit the number of concurrent goroutines
		if goroutineCount >= maxGoroutines {
			wg.Wait()
		}
	}

	// Wait for any remaining goroutines to finish
	wg.Wait()
	return interestingJobs
}
