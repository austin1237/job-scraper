package job

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeduplicatedLinks(t *testing.T) {
	jobs := []Job{
		{Link: "http://example.com/job1"},
		{Link: "http://example.com/job2"},
		{Link: "http://example.com/job1"},
	}

	deduplicated := DeduplicatedLinks(jobs)

	assert.Equal(t, 2, len(deduplicated), "Expected 2 jobs")
	assert.Equal(t, "http://example.com/job1", deduplicated[0].Link, "Expected http://example.com/job1")
	assert.Equal(t, "http://example.com/job2", deduplicated[1].Link, "Expected http://example.com/job2")
}

func TestGetJobHtml(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test that the proxy URL is correctly appended to the request URL
		if strings.HasPrefix(req.URL.String(), "/proxy?url=") {
			rw.Write([]byte("<html><body>Hello, World!</body></html>"))
		} else {
			http.Error(rw, "Invalid request URL", http.StatusBadRequest)
		}
	}))
	defer server.Close()

	// Test the GetJobHtml function
	doc, err := GetJobHtml("https://example.com", server.URL)
	require.NoError(t, err, "GetJobHtml returned an error")

	// Test that the returned *goquery.Document has the correct HTML
	require.Equal(t, "Hello, World!", doc.Find("body").Text(), "doc.Find(\"body\").Text() returned incorrect text")

	// Test the GetJobHtml function with an invalid URL
	_, err = GetJobHtml("invalid url", server.URL)
	require.Error(t, err, "GetJobHtml should return an error for invalid URLs")
}

// MockJobParser is a mock job parser.
type MockJobParser struct {
	mock.Mock
}

// Parse is a mock method that returns a slice of jobs.
func (m *MockJobParser) Parse(baseURL string, doc *goquery.Document) []Job {
	args := m.Called(baseURL, doc)
	return args.Get(0).([]Job)
}

func TestGetNewJobs(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("<html><body>Hello, World!</body></html>"))
	}))
	defer server.Close()

	// Create a mock job parser
	mockJobParser := new(MockJobParser)
	mockJobParser.On("Parse", mock.Anything, mock.Anything).Return([]Job{{Title: "Test Job"}})

	// Test the GetNewJobs function
	jobs := GetNewJobs("https://example.com", server.URL, mockJobParser.Parse)
	require.Len(t, jobs, 1, "GetNewJobs should return 1 job")
	require.Equal(t, "Test Job", jobs[0].Title, "jobs[0].Title should be 'Test Job'")
}
