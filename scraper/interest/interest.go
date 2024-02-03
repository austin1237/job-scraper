package interest

import (
	"fmt"
	"scraper/job"
	"strings"
	"sync"
)

func checkIfInterested(description string) bool {
	keywords := []string{"typescript", "node", "nodejs", "node.js", "go ", "golang", "go,", "go)", "(go"}
	// Check if keywords are present in the job's text
	descriptionToLower := strings.ToLower(description)
	for _, keyword := range keywords {
		if strings.Contains(descriptionToLower, keyword) {
			return true
		}
	}
	return false
}

type JobInfoGetter func(string, string) (string, error)

func FilterInterest(proxyUrl string, Jobs []job.Job, jobInfoGetter JobInfoGetter) []job.Job {
	interestingJobs := []job.Job{}
	var wg sync.WaitGroup

	// Number of concurrent goroutines
	maxGoroutines := 10
	var goroutineCount int

	for _, possibleJob := range Jobs {
		wg.Add(1)
		goroutineCount++

		go func(possibleJob job.Job) {
			defer wg.Done()
			description, err := jobInfoGetter(possibleJob.Link, proxyUrl)
			if err != nil {
				fmt.Println(err)
			}
			if checkIfInterested(description) {
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
