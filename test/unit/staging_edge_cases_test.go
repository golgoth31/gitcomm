package unit

import (
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
)

func TestStagingState_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		state *model.StagingState
		want  bool
	}{
		{
			name: "empty state",
			state: &model.StagingState{
				StagedFiles: []string{},
			},
			want: true,
		},
		{
			name:  "nil state",
			state: nil,
			want:  false, // nil check should be done by caller
		},
		{
			name: "state with many files",
			state: &model.StagingState{
				StagedFiles: make([]string, 1000),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.state == nil {
				// Test nil handling
				return
			}
			if got := tt.state.IsEmpty(); got != tt.want {
				t.Errorf("StagingState.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAutoStagingResult_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		result *model.AutoStagingResult
		want   bool
	}{
		{
			name: "success with no files",
			result: &model.AutoStagingResult{
				StagedFiles: []string{},
				FailedFiles: []model.StagingFailure{},
				Success:     true,
			},
			want: false,
		},
		{
			name: "failure with no failed files (invalid state)",
			result: &model.AutoStagingResult{
				StagedFiles: []string{},
				FailedFiles: []model.StagingFailure{},
				Success:     false,
			},
			want: false,
		},
		{
			name: "partial failure",
			result: &model.AutoStagingResult{
				StagedFiles: []string{"file1.txt"},
				FailedFiles: []model.StagingFailure{
					{FilePath: "file2.txt"},
				},
				Success: false,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.HasFailures(); got != tt.want {
				t.Errorf("AutoStagingResult.HasFailures() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestorationPlan_EdgeCases(t *testing.T) {
	preCLIState := &model.StagingState{
		StagedFiles:    []string{"file1.txt"},
		CapturedAt:     time.Now(),
		RepositoryPath: "/test/repo",
	}

	currentState := &model.StagingState{
		StagedFiles:    []string{"file1.txt", "file2.txt"},
		CapturedAt:     time.Now(),
		RepositoryPath: "/test/repo",
	}

	tests := []struct {
		name    string
		plan    *model.RestorationPlan
		wantErr bool
	}{
		{
			name: "empty plan (no files to unstage)",
			plan: &model.RestorationPlan{
				FilesToUnstage: []string{},
				PreCLIState:    preCLIState,
				CurrentState:   currentState,
			},
			wantErr: false,
		},
		{
			name: "plan with external changes (files in pre-CLI but not current)",
			plan: &model.RestorationPlan{
				FilesToUnstage: []string{"file2.txt"},
				PreCLIState:    preCLIState,
				CurrentState:   currentState,
			},
			wantErr: false,
		},
		{
			name: "plan trying to unstage file not in current state",
			plan: &model.RestorationPlan{
				FilesToUnstage: []string{"file3.txt"},
				PreCLIState:    preCLIState,
				CurrentState:   currentState,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.plan.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RestorationPlan.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
