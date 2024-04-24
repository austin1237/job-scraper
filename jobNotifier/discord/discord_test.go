package discord

import (
	"jobNotifier/job"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMessages(t *testing.T) {
	jobs := []job.Job{
		{Link: "http://example.com/job1", Company: "Company1", Keyword: "Go"},
		{Link: "http://example.com/job2", Company: "Company2"},
		{Link: "http://example.com/job3", Company: "Company3"},
		// Add more jobs to test the 2000 character limit
	}

	messages := generateMessages(jobs)

	// Check that each message is less than or equal to 2000 characters
	for _, message := range messages {
		assert.True(t, len(message) <= 2000, "Message length should be less than or equal to 2000 characters")
	}

	// Check that all jobs are included in the messages
	for _, job := range jobs {
		jobLine := job.Link + ", " + job.Company
		if job.Keyword != "" {
			jobLine += ", " + job.Keyword
		}
		found := false
		for _, message := range messages {
			if strings.Contains(message, jobLine) {
				found = true
				break
			}
		}
		assert.True(t, found, "All jobs should be included in the messages")
	}
}

func TestGenerateMessages_MultipleMessages(t *testing.T) {
	// Create a job with a link and company name that together are 200 characters long
	newJob := job.Job{
		Link:    strings.Repeat("a", 100), // = 100
		Company: strings.Repeat("b", 97),  // ", " and the ending "\n" is 3 characters, so 97 + 3 = 100
	}

	// Create 11 jobs, which should result in a total length of 2200 of job text characters
	jobs := make([]job.Job, 11)
	for i := range jobs {
		jobs[i] = newJob
	}

	messages := generateMessages(jobs)

	// Check that multiple messages were created
	assert.True(t, len(messages) == 2, "Multiple messages should be created when the total length of the jobs exceeds 2000 characters")
	// The addional 6 characters are the "```" and "```" characters at the start and end of the message
	assert.True(t, len(messages[0]) == 1806, "The first message should be 1806 characters long")
	assert.True(t, len(messages[1]) == 406, "The second message should be 406 characters long")
}
