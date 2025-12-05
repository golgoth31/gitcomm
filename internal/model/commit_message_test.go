package model

import (
	"testing"
)

func TestCommitMessage_IsEmpty(t *testing.T) {
	tests := []struct {
		name    string
		message CommitMessage
		want    bool
	}{
		{
			name:    "empty message",
			message: CommitMessage{},
			want:    true,
		},
		{
			name: "message with type only",
			message: CommitMessage{
				Type: "feat",
			},
			want: true,
		},
		{
			name: "message with type and subject",
			message: CommitMessage{
				Type:    "feat",
				Subject: "add new feature",
			},
			want: false,
		},
		{
			name: "message with subject only",
			message: CommitMessage{
				Subject: "add new feature",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.message.IsEmpty(); got != tt.want {
				t.Errorf("CommitMessage.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
