package watcher

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// TestRunner handles continuous test execution
type TestRunner struct {
	watcher    *fsnotify.Watcher
	done       chan bool
	debouncer  *time.Timer
	workDir    string
	testCache  map[string]time.Time
	onSuccess  func()
	onFailure  func(error)
}

// New creates a new TestRunner
func New(workDir string) (*TestRunner, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %v", err)
	}

	return &TestRunner{
		watcher:   watcher,
		done:      make(chan bool),
		workDir:   workDir,
		testCache: make(map[string]time.Time),
	}, nil
}

// Start begins watching for file changes and running tests
func (tr *TestRunner) Start() error {
	// Add all Go files to the watcher
	err := filepath.Walk(tr.workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor and .git directories
		if info.IsDir() && (info.Name() == "vendor" || info.Name() == ".git") {
			return filepath.SkipDir
		}

		// Watch directories for file creation/deletion
		if info.IsDir() {
			return tr.watcher.Add(path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %v", err)
	}

	go tr.watch()
	return nil
}

// watch monitors file changes and triggers test runs
func (tr *TestRunner) watch() {
	tr.debouncer = time.NewTimer(0)
	<-tr.debouncer.C // Consume initial trigger

	for {
		select {
		case event, ok := <-tr.watcher.Events:
			if !ok {
				return
			}

			if tr.shouldRunTests(event) {
				// Reset debouncer
				if !tr.debouncer.Stop() {
					select {
					case <-tr.debouncer.C:
					default:
					}
				}
				tr.debouncer.Reset(500 * time.Millisecond)

				// Run tests after debounce period
				go func() {
					<-tr.debouncer.C
					tr.runTests(event.Name)
				}()
			}

		case err, ok := <-tr.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)

		case <-tr.done:
			return
		}
	}
}

// shouldRunTests determines if tests should be run for a given file event
func (tr *TestRunner) shouldRunTests(event fsnotify.Event) bool {
	// Only run tests on Go file changes
	if !strings.HasSuffix(event.Name, ".go") {
		return false
	}

	// Skip temporary files
	if strings.HasSuffix(event.Name, ".go~") || strings.HasPrefix(filepath.Base(event.Name), ".") {
		return false
	}

	// Run tests on file modifications
	return event.Op&(fsnotify.Write|fsnotify.Create) != 0
}

// runTests executes the test suite
func (tr *TestRunner) runTests(changedFile string) {
	// Determine which packages to test
	pkgs := []string{"./..."}
	if !strings.Contains(changedFile, "_test.go") {
		// If a non-test file changed, find its corresponding test file
		pkg := "./" + filepath.Dir(changedFile)
		pkgs = []string{pkg}
	}

	// Run tests
	cmd := exec.Command("go", append([]string{"test", "-v"}, pkgs...)...)
	cmd.Dir = tr.workDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Tests failed:\n%s", output)
		if tr.onFailure != nil {
			tr.onFailure(fmt.Errorf("tests failed: %v", err))
		}
	} else {
		log.Printf("Tests passed:\n%s", output)
		if tr.onSuccess != nil {
			tr.onSuccess()
		}
	}
}

// OnSuccess sets the callback for successful test runs
func (tr *TestRunner) OnSuccess(callback func()) {
	tr.onSuccess = callback
}

// OnFailure sets the callback for failed test runs
func (tr *TestRunner) OnFailure(callback func(error)) {
	tr.onFailure = callback
}

// Stop terminates the test runner
func (tr *TestRunner) Stop() {
	tr.done <- true
	tr.watcher.Close()
} 
