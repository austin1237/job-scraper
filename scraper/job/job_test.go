package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
