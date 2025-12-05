package model

import (
	"testing"
	"time"
)

func TestStagingState_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		state *StagingState
		want  bool
	}{
		{
			name: "empty state",
			state: &StagingState{
				StagedFiles: []string{},
			},
			want: true,
		},
		{
			name: "non-empty state",
			state: &StagingState{
				StagedFiles: []string{"file1.txt"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsEmpty(); got != tt.want {
				t.Errorf("StagingState.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStagingState_Contains(t *testing.T) {
	state := &StagingState{
		StagedFiles: []string{"file1.txt", "file2.txt"},
	}

	tests := []struct {
		name string
		file string
		want bool
	}{
		{
			name: "file exists",
			file: "file1.txt",
			want: true,
		},
		{
			name: "file does not exist",
			file: "file3.txt",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := state.Contains(tt.file); got != tt.want {
				t.Errorf("StagingState.Contains(%q) = %v, want %v", tt.file, got, tt.want)
			}
		})
	}
}

func TestStagingState_Diff(t *testing.T) {
	tests := []struct {
		name  string
		this  *StagingState
		other *StagingState
		want  []string
	}{
		{
			name: "files in this but not other",
			this: &StagingState{
				StagedFiles: []string{"file1.txt", "file2.txt", "file3.txt"},
			},
			other: &StagingState{
				StagedFiles: []string{"file1.txt"},
			},
			want: []string{"file2.txt", "file3.txt"},
		},
		{
			name: "no difference",
			this: &StagingState{
				StagedFiles: []string{"file1.txt"},
			},
			other: &StagingState{
				StagedFiles: []string{"file1.txt"},
			},
			want: []string{},
		},
		{
			name: "other is nil",
			this: &StagingState{
				StagedFiles: []string{"file1.txt"},
			},
			other: nil,
			want:  []string{"file1.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.this.Diff(tt.other)
			if len(got) != len(tt.want) {
				t.Errorf("StagingState.Diff() length = %d, want %d", len(got), len(tt.want))
				return
			}

			gotSet := make(map[string]bool)
			for _, f := range got {
				gotSet[f] = true
			}
			for _, f := range tt.want {
				if !gotSet[f] {
					t.Errorf("StagingState.Diff() missing file %q", f)
				}
			}
		})
	}
}

func TestAutoStagingResult_HasFailures(t *testing.T) {
	tests := []struct {
		name   string
		result *AutoStagingResult
		want   bool
	}{
		{
			name: "has failures",
			result: &AutoStagingResult{
				FailedFiles: []StagingFailure{
					{FilePath: "file1.txt"},
				},
			},
			want: true,
		},
		{
			name: "no failures",
			result: &AutoStagingResult{
				FailedFiles: []StagingFailure{},
			},
			want: false,
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

func TestAutoStagingResult_GetFailedFilePaths(t *testing.T) {
	result := &AutoStagingResult{
		FailedFiles: []StagingFailure{
			{FilePath: "file1.txt"},
			{FilePath: "file2.txt"},
		},
	}

	got := result.GetFailedFilePaths()
	want := []string{"file1.txt", "file2.txt"}

	if len(got) != len(want) {
		t.Errorf("AutoStagingResult.GetFailedFilePaths() length = %d, want %d", len(got), len(want))
		return
	}

	for i, f := range want {
		if got[i] != f {
			t.Errorf("AutoStagingResult.GetFailedFilePaths()[%d] = %q, want %q", i, got[i], f)
		}
	}
}

func TestRestorationPlan_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		plan *RestorationPlan
		want bool
	}{
		{
			name: "empty plan",
			plan: &RestorationPlan{
				FilesToUnstage: []string{},
			},
			want: true,
		},
		{
			name: "non-empty plan",
			plan: &RestorationPlan{
				FilesToUnstage: []string{"file1.txt"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.plan.IsEmpty(); got != tt.want {
				t.Errorf("RestorationPlan.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestorationPlan_Validate(t *testing.T) {
	preCLIState := &StagingState{
		StagedFiles:    []string{"file1.txt"},
		CapturedAt:     time.Now(),
		RepositoryPath: "/test/repo",
	}

	currentState := &StagingState{
		StagedFiles:    []string{"file1.txt", "file2.txt"},
		CapturedAt:     time.Now(),
		RepositoryPath: "/test/repo",
	}

	tests := []struct {
		name    string
		plan    *RestorationPlan
		wantErr bool
	}{
		{
			name: "valid plan",
			plan: &RestorationPlan{
				FilesToUnstage: []string{"file2.txt"},
				PreCLIState:    preCLIState,
				CurrentState:   currentState,
			},
			wantErr: false,
		},
		{
			name: "nil pre-CLI state",
			plan: &RestorationPlan{
				FilesToUnstage: []string{"file2.txt"},
				PreCLIState:    nil,
				CurrentState:   currentState,
			},
			wantErr: true,
		},
		{
			name: "nil current state",
			plan: &RestorationPlan{
				FilesToUnstage: []string{"file2.txt"},
				PreCLIState:    preCLIState,
				CurrentState:   nil,
			},
			wantErr: true,
		},
		{
			name: "file not currently staged",
			plan: &RestorationPlan{
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
