package watcher

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "watcher-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test creating a new TestRunner
	runner, err := New(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create TestRunner: %v", err)
	}
	defer runner.Stop()

	// Verify the TestRunner was initialized correctly
	if runner.workDir != tmpDir {
		t.Errorf("Expected workDir to be %s, got %s", tmpDir, runner.workDir)
	}
	if runner.watcher == nil {
		t.Error("Expected watcher to be initialized")
	}
	if runner.done == nil {
		t.Error("Expected done channel to be initialized")
	}
	if runner.testCache == nil {
		t.Error("Expected testCache to be initialized")
	}
	if runner.analyzer == nil {
		t.Error("Expected analyzer to be initialized")
	}
}

func TestShouldRunTests(t *testing.T) {
	runner, err := New(".")
	if err != nil {
		t.Fatalf("Failed to create TestRunner: %v", err)
	}
	defer runner.Stop()

	tests := []struct {
		name     string
		fileName string
		op       fsnotify.Op
		want     bool
	}{
		{
			name:     "go file write",
			fileName: "test.go",
			op:       fsnotify.Write,
			want:     true,
		},
		{
			name:     "non-go file",
			fileName: "test.txt",
			op:       fsnotify.Write,
			want:     false,
		},
		{
			name:     "temp go file",
			fileName: "test.go~",
			op:       fsnotify.Write,
			want:     false,
		},
		{
			name:     "hidden go file",
			fileName: ".test.go",
			op:       fsnotify.Write,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := fsnotify.Event{
				Name: tt.fileName,
				Op:   tt.op,
			}
			if got := runner.shouldRunTests(event); got != tt.want {
				t.Errorf("shouldRunTests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartAndStop(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(testFile, []byte("package test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create a new TestRunner
	runner, err := New(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create TestRunner: %v", err)
	}

	// Start watching
	if err := runner.Start(); err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}

	// Set up success callback
	successCalled := make(chan bool)
	runner.OnSuccess(func() {
		successCalled <- true
	})

	// Modify the test file
	time.Sleep(100 * time.Millisecond) // Wait for watcher to start
	if err := os.WriteFile(testFile, []byte("package test // modified"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Wait for callback or timeout
	select {
	case <-successCalled:
		// Success callback was called
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for success callback")
	}

	// Stop the watcher
	runner.Stop()
} 
