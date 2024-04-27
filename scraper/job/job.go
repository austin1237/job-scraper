package job

import (
	"bytes"
	"encoding/json"
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

type Response struct {
	Total      int `json:"total"`
	Uncached   int `json:"uncached"`
	Duplicates int `json:"duplicates"`
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

func GetJobHtml(siteUrl string, proxyURL string, optionalRoute ...string) (*goquery.Document, error) {
	var route string
	if len(optionalRoute) > 0 {
		route = optionalRoute[0]
	} else {
		route = "proxy" // default mode
	}
	response, err := http.Get(proxyURL + "/" + route + "?url=" + siteUrl)
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

func GetNewJobs(siteUrl string, proxyURL string, jobParser parser, optionalMode ...string) []Job {
	var mode string
	if len(optionalMode) > 0 {
		mode = optionalMode[0]
	} else {
		mode = "proxy" // default mode
	}
	u, err := url.Parse(siteUrl)
	baseURL := u.Scheme + "://" + u.Host
	if err != nil {
		log.Println("Failed to parse url", err)
		return []Job{}
	}
	doc, err := GetJobHtml(siteUrl, proxyURL, mode)
	if err != nil {
		return []Job{}
	}
	return jobParser(baseURL, doc)
}

func SendJobs(jobURL string, jobs []Job) (Response, error) {
	var response Response
	jsonData, err := json.Marshal(map[string][]Job{"jobs": jobs})
	if err != nil {
		return response, err
	}

	resp, err := http.Post(jobURL+"/job", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}
