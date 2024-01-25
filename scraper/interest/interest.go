package interest

import (
	"fmt"
	"scraper/scanner"
	"strings"
	"sync"
)

func checkIfInterested(job scanner.Job) bool {
	keywords := []string{"node", "nodejs", "go ", "golang", "go,", "go)", "(go"}
	// Check if keywords are present in the job's text
	for _, keyword := range keywords {
		if strings.Contains(job.Description, keyword) {
			return true
		}
	}
	return false
}

func FilterInterest(siteAbaseUrl string, proxyUrl string, links []string) []scanner.Job {
	interestingJobs := []scanner.Job{}
	var wg sync.WaitGroup

	// Number of concurrent goroutines
	maxGoroutines := 10
	var goroutineCount int

	for _, link := range links {
		wg.Add(1)
		goroutineCount++

		go func(link string) {
			defer wg.Done()
			jobUrl := siteAbaseUrl + link
			job, err := scanner.GetSiteAJobInfo(jobUrl, proxyUrl)
			if err != nil {
				fmt.Println(err)
			}
			if checkIfInterested(job) {
				interestingJobs = append(interestingJobs, job)
			}
			goroutineCount--
		}(link)

		// Limit the number of concurrent goroutines
		if goroutineCount >= maxGoroutines {
			wg.Wait()
		}
	}

	// Wait for any remaining goroutines to finish
	wg.Wait()
	return interestingJobs
}
