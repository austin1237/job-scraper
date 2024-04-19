package dynamo

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBAPI interface {
	UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

type Table struct {
	Name string
	svc  DynamoDBAPI
}

func NewTable(name string, region string) (*Table, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region), // replace with your region
	})
	if err != nil {
		return nil, err
	}

	svc := dynamodb.New(sess)

	return &Table{Name: name, svc: svc}, nil
}

func (t *Table) ReadItem(company string) (string, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(t.Name),
		Key: map[string]*dynamodb.AttributeValue{
			"company": {
				S: aws.String(strings.ToLower(company)),
			},
		},
	}

	result, err := t.svc.GetItem(input)
	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", nil
	}

	return *result.Item["company"].S, nil
}

func (t *Table) WriteItems(companies []string) {
	// Set the ttl time to 30 days from now
	expirationTime := time.Now().AddDate(0, 1, 0).Unix()

	// Create a wait group
	var wg sync.WaitGroup

	// Write each company to the table in a separate goroutine
	for _, company := range companies {
		wg.Add(1)
		go func(company string) {
			defer wg.Done()

			input := &dynamodb.UpdateItemInput{
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":expirationTime": {
						N: aws.String(strconv.FormatInt(expirationTime, 10)),
					},
				},
				TableName: aws.String(t.Name),
				Key: map[string]*dynamodb.AttributeValue{
					"company": {
						S: aws.String(strings.ToLower(company)),
					},
				},
				ReturnValues:     aws.String("UPDATED_NEW"),
				UpdateExpression: aws.String("set ExpirationTime = :expirationTime"),
			}

			_, err := t.svc.UpdateItem(input)
			if err != nil {
				log.Println("Error writing company to cache", err)
			}
		}(company)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
