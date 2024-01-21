package interest

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Job struct {
	Text string `json:"text"`
	Name string `json:"name"`
	Link string
}

func getJobInfo(siteUrl string, proxyUrl string) Job {
	// Make an HTTP GET request
	// query string the url
	job := Job{
		Link: siteUrl,
	}

	siteUrl = url.QueryEscape(siteUrl)
	response, err := http.Get(proxyUrl + "/counter?url=" + siteUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := errors.New("source HTTP request failed with status: " + response.Status)
		log.Fatal(err)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return job
	}

	// Parse JSON data into a Job struct
	err = json.Unmarshal(body, &job)
	if err != nil {
		return job
	}

	return job
}

func checkIfInterested(job Job) bool {
	keywords := []string{"node", "nodejs", "go ", "golang", "go,", "go)", "(go"}
	// Check if keywords are present in the job's text
	for _, keyword := range keywords {
		if strings.Contains(job.Text, keyword) {
			return true
		}
	}
	return false
}

func FilterInterest(siteAbaseUrl string, proxyUrl string, links []string) []Job {
	interestingJobs := []Job{}
	var wg sync.WaitGroup

	// Number of concurrent goroutines
	maxGoroutines := 10
	var goroutineCount int

	for _, link := range links {
		wg.Add(1)
		goroutineCount++

		go func(link string) {
			defer wg.Done()
			siteUrl := siteAbaseUrl + link
			job := getJobInfo(siteUrl, proxyUrl)
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
