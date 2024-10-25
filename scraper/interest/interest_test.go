package interest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckIfInterested(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{
			name:        "contains typescript",
			description: "This is a job for a TypeScript developer",
			want:        "typescript",
		},
		{
			name:        "contains node",
			description: "This is a job for a Node.js developer",
			want:        "node",
		},
		{
			name:        "contains go",
			description: "This is a job for a Go developer",
			want:        "go",
		},
		{
			name:        "contains go",
			description: "This is a job for go.",
			want:        "go",
		},
		{
			name:        "contains deno",
			description: "This is a job for a Deno developer",
			want:        "deno",
		},
		{
			name:        "contains bun",
			description: "This is a job for a Bun developer",
			want:        "bun",
		},
		{
			name:        "contains django",
			description: "This is a job for a django developer in Chicago",
			want:        "",
		},
		{
			name:        "does not contain keyword",
			description: "This is a job for a Python developer",
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CheckIfInterested(tt.description))
		})
	}
}
