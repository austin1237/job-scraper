package scanner

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Post struct {
	Title string
	Link  string
	Time  string
}

func scanSiteA(siteABaseUrl string) []string {
	var links []string
	finished := false
	page := 1
	endTime := "yesterday"

	for !finished {
		pageStr := strconv.Itoa(page)
		url := siteABaseUrl + "/jobs/remote/nationwide/dev-engineering?page=" + pageStr
		// Make an HTTP GET request
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Fatalf("HTTP request failed with status: %s", response.Status)
		}

		// Parse the HTML document using goquery
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Find the job link and title /job/engineer/...
		doc.Find("div#search-results-bottom a#job-card-alias").Each(func(index int, item *goquery.Selection) {
			// Extract and print the href attribute
			href, exists := item.Attr("href")
			title := item.Text()
			title = strings.ToLower(title)
			// Only grab releveant jobs
			if (strings.Contains(title, "software engineer") || strings.Contains(title, "developer")) && exists {
				links = append(links, href)
			}
		})

		// Find the job link and title /job/engineer/...
		doc.Find("div#search-results-top a#job-card-alias").Each(func(index int, item *goquery.Selection) {
			// Extract and print the href attribute
			href, exists := item.Attr("href")
			title := item.Text()
			title = strings.ToLower(title)
			// Only grab releveant jobs
			if (strings.Contains(title, "software engineer") || strings.Contains(title, "developer")) && exists {
				links = append(links, href)
			}
		})

		// check all elements and see if see we no longer need to paginate
		doc.Find("div#search-results-bottom").Each(func(_ int, s *goquery.Selection) {
			// Check the text of each element
			s.Find("*").Each(func(_ int, e *goquery.Selection) {
				// fmt.Println(e.Text())
				if strings.Contains(strings.ToLower(e.Text()), endTime) {
					finished = true
				}
			})
		})

		// check all elements and see if see we no longer need to paginate
		doc.Find("div#search-results-top").Each(func(_ int, s *goquery.Selection) {
			// Check the text of each element
			s.Find("*").Each(func(_ int, e *goquery.Selection) {
				if strings.Contains(strings.ToLower(e.Text()), endTime) {
					finished = true
				}
			})
		})
		fmt.Println("page: " + pageStr)
		page++
	}
	links = deduplicatedLinks(links)
	return links
}

func deduplicatedLinks(links []string) []string {
	// Create a map to track unique links
	uniqueLinks := make(map[string]struct{})

	// Create a new deduplicated array
	deduplicatedLinks := []string{}

	// Iterate through the original array and add unique links to deduplicatedData
	for _, link := range links {
		if _, exists := uniqueLinks[link]; !exists {
			uniqueLinks[link] = struct{}{}
			deduplicatedLinks = append(deduplicatedLinks, link)
		}
	}

	return deduplicatedLinks
}

func ScanNewJobs(siteABaseUrl string) []string {
	links := scanSiteA(siteABaseUrl)
	return links
}
