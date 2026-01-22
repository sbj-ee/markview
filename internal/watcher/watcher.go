package watcher

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

// FileWatcher watches a file for changes and triggers callbacks
type FileWatcher struct {
	watcher      *fsnotify.Watcher
	logger       *zap.Logger
	filePath     string
	onChange     func()
	debouncer    *Debouncer
	mu           sync.Mutex
	running      bool
	stopChan     chan struct{}
	debounceTime time.Duration
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(logger *zap.Logger, debounceTime time.Duration) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		watcher:      watcher,
		logger:       logger,
		debounceTime: debounceTime,
		stopChan:     make(chan struct{}),
	}

	return fw, nil
}

// Watch starts watching a file for changes
func (fw *FileWatcher) Watch(filePath string, onChange func()) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// Stop previous watch if any
	if fw.running {
		fw.stopWatch()
	}

	// Create new stop channel for this watch session
	fw.stopChan = make(chan struct{})

	fw.filePath = filePath
	fw.onChange = onChange

	// Create debouncer
	fw.debouncer = NewDebouncer(fw.debounceTime, onChange)

	// Watch the parent directory (some editors replace files atomically)
	dir := filepath.Dir(filePath)
	err := fw.watcher.Add(dir)
	if err != nil {
		return err
	}

	fw.running = true
	go fw.watchLoop()

	fw.logger.Info("Started watching file", zap.String("path", filePath))
	return nil
}

// Stop stops watching the file
func (fw *FileWatcher) Stop() {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.running {
		fw.stopWatch()
	}
}

// stopWatch stops the watch loop (must be called with lock held)
func (fw *FileWatcher) stopWatch() {
	close(fw.stopChan)
	fw.running = false

	if fw.debouncer != nil {
		fw.debouncer.Stop()
	}

	fw.logger.Info("Stopped watching file", zap.String("path", fw.filePath))
}

// watchLoop is the main event loop for file watching
func (fw *FileWatcher) watchLoop() {
	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// Only process events for the watched file
			if event.Name == fw.filePath || filepath.Base(event.Name) == filepath.Base(fw.filePath) {
				fw.handleEvent(event)
			}

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			fw.logger.Error("File watcher error", zap.Error(err))

		case <-fw.stopChan:
			return
		}
	}
}

// handleEvent processes a file system event
func (fw *FileWatcher) handleEvent(event fsnotify.Event) {
	fw.logger.Debug("File event",
		zap.String("event", event.Op.String()),
		zap.String("file", event.Name),
	)

	// Trigger debounced callback on write or create events
	if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
		fw.debouncer.Trigger()
	}
}

// Close closes the file watcher and releases resources
func (fw *FileWatcher) Close() error {
	fw.Stop()
	return fw.watcher.Close()
}
