package main

import (
	"context"
	"encoding/json"
	"fmt"
	"jobNotifier/cache"
	"jobNotifier/discord"
	"jobNotifier/dynamo"
	"jobNotifier/job"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestBody struct {
	Jobs []job.Job `json:"jobs"`
}

type Results struct {
	Total      int `json:"total"`
	Uncached   int `json:"uncached"`
	Duplicates int `json:"duplicates"`
}

var (
	scraperWebhook string
	dynamoTable    string
)

func init() {
	scraperWebhook = os.Getenv("SCRAPER_WEBHOOK")
	if scraperWebhook == "" {
		log.Fatal("Environment variable SCRAPER_WEBHOOK must be set")
	}

	dynamoTable = os.Getenv("DYNAMO_TABLE")
	if dynamoTable == "" {
		log.Fatal("Environment variable DYNAMO_TABLE must be set")
	}

}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" && os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		offlineHandler()
	} else {
		lambda.Start(handler)
	}
}

func notifyNewJobs(jobs []job.Job) (Results, error) {
	results := Results{}
	results.Total = len(jobs)
	table, err := dynamo.NewTable(dynamoTable, "us-east-1") // replace with your table name
	if err != nil {
		log.Fatal(err)
	}
	cache := cache.NewCache(table)
	unCachedJobs, err := cache.FilterCachedCompanies(jobs)
	results.Uncached = len(unCachedJobs)
	results.Duplicates = results.Total - results.Uncached
	if err != nil {
		return results, err
	}
	errs := discord.SendJobsToDiscord(unCachedJobs, scraperWebhook)
	if len(errs) == 0 {
		cache.WriteCompaniesToCache(unCachedJobs)
	} else {
		return results, fmt.Errorf("error sending to discord %v", errs)
	}
	return results, nil
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	results, err := notifyNewJobs(requestBody.Jobs)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to notify new jobs",
		}, nil
	}

	resultsBytes, err := json.Marshal(results)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to convert results to JSON",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(resultsBytes),
	}, nil
}

func offlineHandler() {
	mockJobs := []job.Job{
		{
			Title:   "Software Engineer",
			Company: "test1",
			Keyword: "Go",
			Link:    "https://testlink1.com",
		},
		{
			Title:   "Data Analyst",
			Company: "test2",
			Keyword: "Go",
			Link:    "https://testlink1.com",
		},
		{
			Title:   "Financial Advisor",
			Company: "test3",
			Keyword: "Go",
			Link:    "https://testlink3.com",
		},
		{
			Title:   "Educational Consultant",
			Company: "test4",
			Keyword: "Go",
			Link:    "https://testlink4.com",
		},
	}
	results, err := notifyNewJobs(mockJobs)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Total", results.Total, "Uncached", results.Uncached, "Duplicates", results.Duplicates)
	}
}
