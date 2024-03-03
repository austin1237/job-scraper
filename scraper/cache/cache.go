package cache

import (
	"scraper/job"
)

type Table interface {
	ReadItem(company string) (string, error)
	WriteItems(companies []string)
}

type Cache struct {
	table Table
}

func NewCache(table Table) *Cache {
	return &Cache{table: table}
}

func (c *Cache) FilterCachedCompanies(jobs []job.Job) ([]job.Job, error) {
	notInCache := make([]job.Job, 0)
	errChan := make(chan error, len(jobs))
	notFoundChan := make(chan job.Job, len(jobs))
	foundChan := make(chan job.Job, len(jobs))

	for _, newJob := range jobs {
		go func(newJob job.Job) {
			result, err := c.table.ReadItem(newJob.Company)
			if result == "" {
				// company is not in the cache
				notFoundChan <- newJob
			} else {
				foundChan <- newJob
			}

			if err != nil {
				errChan <- err
			}

		}(newJob)
	}

	// Collect results from the goroutines
	for range jobs {
		select {
		case job := <-notFoundChan:
			notInCache = append(notInCache, job)
		case <-foundChan:
			// do nothing
		case err := <-errChan:
			return nil, err
		}

	}

	return notInCache, nil
}

func (c *Cache) WriteCompaniesToCache(jobs []job.Job) {
	companies := make([]string, 0, len(jobs))
	for _, job := range jobs {
		companies = append(companies, job.Company)
	}
	c.table.WriteItems(companies)
}
