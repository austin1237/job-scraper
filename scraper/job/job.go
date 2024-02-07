package job

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Job struct {
	Title   string
	Company string
	Link    string
}

func DeduplicatedLinks(jobs []Job) []Job {
	seen := make(map[string]bool)
	deduplicated := []Job{}

	for _, possibleJob := range jobs {
		if !seen[possibleJob.Link] {
			seen[possibleJob.Link] = true
			deduplicated = append(deduplicated, possibleJob)
		}
	}
	return deduplicated
}

func GetJobHtml(siteUrl string, proxyURL string) (*goquery.Document, error) {
	response, err := http.Get(proxyURL + "/proxy?url=" + siteUrl)
	if err != nil {
		log.Println(siteUrl+": Failed to get site", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := errors.New(siteUrl + ": HTTP request failed with status: " + response.Status)
		log.Println(err.Error())
		return nil, err
	}

	// Parse the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println(siteUrl+": Failed to parse site", err)
		return nil, err
	}

	return doc, nil

}

type parser func(string, *goquery.Document) []Job

func GetNewJobs(siteUrl string, proxyURL string, jobParser parser) []Job {
	u, err := url.Parse(siteUrl)
	baseURL := u.Scheme + "://" + u.Host
	if err != nil {
		log.Println("Failed to parse url", err)
		return []Job{}
	}
	doc, err := GetJobHtml(siteUrl, proxyURL)
	if err != nil {
		return []Job{}
	}
	return jobParser(baseURL, doc)
}
