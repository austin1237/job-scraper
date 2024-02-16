package interest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckIfInterested(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        bool
	}{
		{
			name:        "contains typescript",
			description: "This is a job for a TypeScript developer",
			want:        true,
		},
		{
			name:        "contains node",
			description: "This is a job for a Node.js developer",
			want:        true,
		},
		{
			name:        "contains go",
			description: "This is a job for a Go developer",
			want:        true,
		},
		{
			name:        "contains go",
			description: "This is a job for go.",
			want:        true,
		},
		{
			name:        "contains django",
			description: "This is a job for a django developer in Chicago",
			want:        false,
		},
		{
			name:        "does not contain keyword",
			description: "This is a job for a Python developer",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CheckIfInterested(tt.description))
		})
	}
}
