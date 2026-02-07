// Copyright 2026 Grand Valley State University
// Copyright 2020 Trey Dockendorf
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/containerd/cgroups/v3/cgroup2"
)

func TestGetNamev2WithSlurmPath(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	tests := []struct {
		name     string
		pidPath  string
		path     string
		expected []string
	}{
		{
			name:    "Full SLURM task path",
			pidPath: "/system.slice/slurmstepd.scope/job_123/step_0/user/task_1",
			path:    "/system.slice/slurmstepd.scope",
			expected: []string{
				"/system.slice/slurmstepd.scope/job_123/step_0/user/task_1",
				"/system.slice/slurmstepd.scope/job_123",
			},
		},
		{
			name:    "SLURM job level path",
			pidPath: "/system.slice/slurmstepd.scope/job_456",
			path:    "/system.slice/slurmstepd.scope",
			expected: []string{
				"/system.slice/slurmstepd.scope/job_456",
				"/system.slice/slurmstepd.scope/job_456",
			},
		},
		{
			name:    "SLURM with step but no task",
			pidPath: "/system.slice/slurmstepd.scope/job_789/step_2",
			path:    "/system.slice/slurmstepd.scope",
			expected: []string{
				"/system.slice/slurmstepd.scope/job_789/step_2",
				"/system.slice/slurmstepd.scope/job_789",
			},
		},
		{
			name:     "Non-SLURM path",
			pidPath:  "/user.slice/user-1000.slice/session-1.scope",
			path:     "/user.slice",
			expected: []string{"/user.slice/user-1000.slice/session-1.scope"},
		},
		{
			name:     "SLURM path without job_ prefix",
			pidPath:  "/system.slice/slurmstepd.scope/nojob",
			path:     "/system.slice/slurmstepd.scope",
			expected: []string{"/system.slice/slurmstepd.scope/nojob"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNamev2(tt.pidPath, tt.path, logger)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d names, got %d", len(tt.expected), len(result))
				return
			}
			for i, name := range result {
				if name != tt.expected[i] {
					t.Errorf("Expected name[%d] = %s, got %s", i, tt.expected[i], name)
				}
			}
		})
	}
}

func TestGetInfov2WithValidSlurmPath(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	tests := []struct {
		name          string
		cgroupName    string
		expectedJob   bool
		expectedJobID string
		expectedStep  string
		expectedTask  string
	}{
		{
			name:          "Full SLURM path with task",
			cgroupName:    "/system.slice/slurmstepd.scope/job_100/step_0/user/task_5",
			expectedJob:   true,
			expectedJobID: "100",
			expectedStep:  "0",
			expectedTask:  "5",
		},
		{
			name:          "SLURM path without task",
			cgroupName:    "/system.slice/slurmstepd.scope/job_200/step_1",
			expectedJob:   true,
			expectedJobID: "200",
			expectedStep:  "1",
			expectedTask:  "",
		},
		{
			name:          "SLURM path with special task",
			cgroupName:    "/system.slice/slurmstepd.scope/job_300/step_0/user/task_special",
			expectedJob:   true,
			expectedJobID: "300",
			expectedStep:  "0",
			expectedTask:  "special",
		},
		{
			name:          "SLURM path job only",
			cgroupName:    "/system.slice/slurmstepd.scope/job_400",
			expectedJob:   true,
			expectedJobID: "400",
			expectedStep:  "",
			expectedTask:  "",
		},
		{
			name:          "Non-SLURM path",
			cgroupName:    "/user.slice/user-1000.slice",
			expectedJob:   false,
			expectedJobID: "",
			expectedStep:  "",
			expectedTask:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &CgroupMetric{}
			// Use empty pids slice for pattern matching test
			getInfov2(tt.cgroupName, []int{}, metric, logger)

			if metric.job != tt.expectedJob {
				t.Errorf("Expected job=%v, got %v", tt.expectedJob, metric.job)
			}
			if metric.jobid != tt.expectedJobID {
				t.Errorf("Expected jobid=%s, got %s", tt.expectedJobID, metric.jobid)
			}
			if metric.step != tt.expectedStep {
				t.Errorf("Expected step=%s, got %s", tt.expectedStep, metric.step)
			}
			if metric.task != tt.expectedTask {
				t.Errorf("Expected task=%s, got %s", tt.expectedTask, metric.task)
			}
		})
	}
}

func TestGetStatv2EdgeCases(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()

	// Test file with valid format
	validFile := filepath.Join(tmpDir, "valid.stat")
	validContent := `anon 1024
file 2048
swapcached 512
`
	if err := os.WriteFile(validFile, []byte(validContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}

	// Test file with invalid format (only one field)
	invalidFile := filepath.Join(tmpDir, "invalid.stat")
	invalidContent := "singlevalue"
	if err := os.WriteFile(invalidFile, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}

	// Test file with non-numeric value
	nonNumericFile := filepath.Join(tmpDir, "nonnumeric.stat")
	nonNumericContent := "anon notanumber"
	if err := os.WriteFile(nonNumericFile, []byte(nonNumericContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}

	// Test file with empty content
	emptyFile := filepath.Join(tmpDir, "empty.stat")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}

	tests := []struct {
		name      string
		statName  string
		path      string
		expectErr bool
		expected  float64
	}{
		{
			name:      "Valid stat retrieval",
			statName:  "anon",
			path:      validFile,
			expectErr: false,
			expected:  1024,
		},
		{
			name:      "Valid stat - different key",
			statName:  "file",
			path:      validFile,
			expectErr: false,
			expected:  2048,
		},
		{
			name:      "Missing key",
			statName:  "nonexistent",
			path:      validFile,
			expectErr: true,
		},
		{
			name:      "Invalid format - single field",
			statName:  "singlevalue",
			path:      invalidFile,
			expectErr: true,
		},
		{
			name:      "Non-numeric value",
			statName:  "anon",
			path:      nonNumericFile,
			expectErr: true,
		},
		{
			name:      "Empty file",
			statName:  "anon",
			path:      emptyFile,
			expectErr: true,
		},
		{
			name:      "Non-existent file",
			statName:  "anon",
			path:      "/does/not/exist",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getStatv2(tt.statName, tt.path)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestGetMetricsv2ErrorHandling(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	exporter := NewExporter([]string{"/dne"}, logger)

	tests := []struct {
		name       string
		cgroupName string
		expectErr  bool
	}{
		{
			name:       "Non-existent cgroup",
			cgroupName: "/does/not/exist",
			expectErr:  true,
		},
		{
			name:       "Invalid path",
			cgroupName: "/invalid/cgroup/path",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use default opts
			opts := cgroup2.WithMountpoint(*CgroupRoot)
			metric, err := exporter.getMetricsv2(tt.cgroupName, []int{}, opts)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if !metric.err {
					t.Errorf("Expected metric.err to be true")
				}
			}
		})
	}
}

func TestCollectv2WithMultiplePaths(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// Test with multiple paths, including invalid ones
	paths := []string{
		"/system.slice/slurmstepd.scope",
		"/invalid/path",
		"/another/invalid/path",
	}

	exporter := NewExporter(paths, logger)
	metrics, err := exporter.collectv2()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Should have at least some metrics (even if errors)
	if len(metrics) == 0 {
		t.Errorf("Expected at least some metrics")
	}

	// Check that error metrics are properly set
	errorCount := 0
	for _, m := range metrics {
		if m.err {
			errorCount++
		}
	}

	if errorCount == 0 {
		t.Logf("Warning: Expected some error metrics for invalid paths")
	}
}

func TestGetNamev2EmptyPath(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	result := getNamev2("", "", logger)
	if len(result) != 1 {
		t.Errorf("Expected 1 name for empty path, got %d", len(result))
	}
	if result[0] != "" {
		t.Errorf("Expected empty string, got %s", result[0])
	}
}

func TestGetInfov2WithProcErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// Save original PidGroupPath
	originalPidGroupPath := PidGroupPath
	defer func() { PidGroupPath = originalPidGroupPath }()

	// Mock PidGroupPath to return test path
	PidGroupPath = func(pid int) (string, error) {
		return "/system.slice/slurmstepd.scope/job_999/step_0/user/task_0", nil
	}

	metric := &CgroupMetric{}
	// Use invalid PIDs to trigger proc errors
	getInfov2("/system.slice/slurmstepd.scope/job_999/step_0/user/task_0", []int{999999}, metric, logger)

	// Should have job info extracted from pattern
	if !metric.job {
		t.Errorf("Expected job to be true")
	}
	if metric.jobid != "999" {
		t.Errorf("Expected jobid=999, got %s", metric.jobid)
	}
}

func TestGetStatv2MultipleKeysInFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create file with multiple keys
	testFile := filepath.Join(tmpDir, "multi.stat")
	content := `key1 100
key2 200
key3 300
key4 400
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}

	// Test retrieving different keys
	tests := []struct {
		key      string
		expected float64
	}{
		{"key1", 100},
		{"key2", 200},
		{"key3", 300},
		{"key4", 400},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result, err := getStatv2(tt.key, testFile)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetNamev2SlurmEdgeCases(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	tests := []struct {
		name        string
		pidPath     string
		path        string
		expectedLen int
		description string
	}{
		{
			name:        "Multiple job_ occurrences",
			pidPath:     "/system.slice/job_old/slurmstepd.scope/job_123",
			path:        "/system.slice/job_old/slurmstepd.scope",
			expectedLen: 2,
			description: "Should find first job_ occurrence",
		},
		{
			name:        "job_ at root",
			pidPath:     "job_456/step_0",
			path:        "slurm",
			expectedLen: 2,
			description: "job_ at beginning of path",
		},
		{
			name:        "SLURM in path, no job_",
			pidPath:     "/system.slice/slurmstepd.scope/other",
			path:        "/system.slice/slurmstepd.scope",
			expectedLen: 1,
			description: "SLURM path without job_ pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNamev2(tt.pidPath, tt.path, logger)
			if len(result) != tt.expectedLen {
				t.Errorf("%s: Expected %d names, got %d", tt.description, tt.expectedLen, len(result))
			}
		})
	}
}

func TestGetInfov2PatternEdgeCases(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	tests := []struct {
		name       string
		cgroupName string
		expectJob  bool
		expectID   string
	}{
		{
			name:       "Multi-digit job ID",
			cgroupName: "/job_123456789",
			expectJob:  true,
			expectID:   "123456789",
		},
		{
			name:       "Single digit job ID",
			cgroupName: "/job_1",
			expectJob:  true,
			expectID:   "1",
		},
		{
			name:       "Job with special characters in step",
			cgroupName: "/job_100/step_batch",
			expectJob:  true,
			expectID:   "100",
		},
		{
			name:       "Almost matching pattern",
			cgroupName: "/job_/step_0",
			expectJob:  false,
			expectID:   "",
		},
		{
			name:       "job without underscore",
			cgroupName: "/job123",
			expectJob:  false,
			expectID:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &CgroupMetric{}
			getInfov2(tt.cgroupName, []int{}, metric, logger)

			if metric.job != tt.expectJob {
				t.Errorf("Expected job=%v, got %v", tt.expectJob, metric.job)
			}
			if metric.jobid != tt.expectID {
				t.Errorf("Expected jobid=%s, got %s", tt.expectID, metric.jobid)
			}
		})
	}
}
