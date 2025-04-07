package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestNew(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a new TestRunner
	runner, err := New(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create TestRunner: %v", err)
	}
	defer runner.Stop()

	if runner.workDir != tmpDir {
		t.Errorf("Expected workDir to be %s, got %s", tmpDir, runner.workDir)
	}

	if runner.watcher == nil {
		t.Error("Expected non-nil watcher")
	}

	if runner.done == nil {
		t.Error("Expected non-nil done channel")
	}

	if runner.testCache == nil {
		t.Error("Expected non-nil testCache")
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
		op      fsnotify.Op
		want    bool
	}{
		{
			name:     "Go file write",
			fileName: "test.go",
			op:       fsnotify.Write,
			want:     true,
		},
		{
			name:     "Go file create",
			fileName: "test.go",
			op:       fsnotify.Create,
			want:     true,
		},
		{
			name:     "Non-Go file",
			fileName: "test.txt",
			op:       fsnotify.Write,
			want:     false,
		},
		{
			name:     "Temporary Go file",
			fileName: "test.go~",
			op:       fsnotify.Write,
			want:     false,
		},
		{
			name:     "Hidden Go file",
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
