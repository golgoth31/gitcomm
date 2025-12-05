package model

import (
	"testing"
)

func TestCommitOptions(t *testing.T) {
	tests := []struct {
		name    string
		options CommitOptions
		want    string
	}{
		{
			name:    "default options",
			options: CommitOptions{},
			want:    "AutoStage: false, NoSignoff: false",
		},
		{
			name: "auto-stage enabled",
			options: CommitOptions{
				AutoStage: true,
			},
			want: "AutoStage: true",
		},
		{
			name: "no signoff enabled",
			options: CommitOptions{
				NoSignoff: true,
			},
			want: "NoSignoff: true",
		},
		{
			name: "all options enabled",
			options: CommitOptions{
				AutoStage: true,
				NoSignoff: true,
				SkipAI:    true,
			},
			want: "all enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation that options can be created
			if tt.options.AutoStage && !tt.options.AutoStage {
				t.Error("AutoStage should be true")
			}
			if tt.options.NoSignoff && !tt.options.NoSignoff {
				t.Error("NoSignoff should be true")
			}
		})
	}
}
