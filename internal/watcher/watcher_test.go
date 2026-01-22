package watcher

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewFileWatcher(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	fw, err := NewFileWatcher(logger, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("NewFileWatcher failed: %v", err)
	}
	defer fw.Close()

	if fw == nil {
		t.Fatal("NewFileWatcher returned nil")
	}

	if fw.debounceTime != 100*time.Millisecond {
		t.Errorf("Expected debounceTime 100ms, got %v", fw.debounceTime)
	}
}

func TestFileWatcher_PauseResume(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	fw, err := NewFileWatcher(logger, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("NewFileWatcher failed: %v", err)
	}
	defer fw.Close()

	// Initially not paused
	if fw.IsPaused() {
		t.Error("FileWatcher should not be paused initially")
	}

	// Pause
	fw.Pause()
	if !fw.IsPaused() {
		t.Error("FileWatcher should be paused after Pause()")
	}

	// Resume
	fw.Resume()
	if fw.IsPaused() {
		t.Error("FileWatcher should not be paused after Resume()")
	}
}

func TestFileWatcher_WatchFile(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create temp directory and file
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	fw, err := NewFileWatcher(logger, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("NewFileWatcher failed: %v", err)
	}
	defer fw.Close()

	var called bool
	var mu sync.Mutex

	err = fw.Watch(tmpFile, func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("Watch failed: %v", err)
	}

	// Modify the file
	time.Sleep(50 * time.Millisecond) // Let watcher settle
	if err := os.WriteFile(tmpFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Wait for debounce + some extra time
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if !called {
		t.Error("Callback was not called after file modification")
	}
}

func TestFileWatcher_PauseIgnoresEvents(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create temp directory and file
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	fw, err := NewFileWatcher(logger, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("NewFileWatcher failed: %v", err)
	}
	defer fw.Close()

	var callCount int
	var mu sync.Mutex

	err = fw.Watch(tmpFile, func() {
		mu.Lock()
		callCount++
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("Watch failed: %v", err)
	}

	// Pause the watcher
	fw.Pause()

	// Modify the file while paused
	time.Sleep(50 * time.Millisecond)
	if err := os.WriteFile(tmpFile, []byte("modified while paused"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Wait for potential callback
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	count := callCount
	mu.Unlock()

	if count != 0 {
		t.Errorf("Callback was called %d times while paused, expected 0", count)
	}
}

func TestFileWatcher_ResumeAfterPause(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create temp directory and file
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	fw, err := NewFileWatcher(logger, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("NewFileWatcher failed: %v", err)
	}
	defer fw.Close()

	var called bool
	var mu sync.Mutex

	err = fw.Watch(tmpFile, func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("Watch failed: %v", err)
	}

	// Pause and then resume
	fw.Pause()
	time.Sleep(50 * time.Millisecond)
	fw.Resume()

	// Modify the file after resume
	time.Sleep(50 * time.Millisecond)
	if err := os.WriteFile(tmpFile, []byte("modified after resume"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Wait for callback
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if !called {
		t.Error("Callback was not called after resume")
	}
}

func TestFileWatcher_Stop(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create temp directory and file
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	fw, err := NewFileWatcher(logger, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("NewFileWatcher failed: %v", err)
	}
	defer fw.Close()

	var called bool
	var mu sync.Mutex

	err = fw.Watch(tmpFile, func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("Watch failed: %v", err)
	}

	// Stop the watcher
	fw.Stop()

	// Modify the file after stop
	time.Sleep(50 * time.Millisecond)
	if err := os.WriteFile(tmpFile, []byte("modified after stop"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Wait for potential callback
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if called {
		t.Error("Callback was called after Stop()")
	}
}
