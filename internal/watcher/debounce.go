package watcher

import (
	"sync"
	"time"
)

// Debouncer delays function execution until after a quiet period
type Debouncer struct {
	mu       sync.Mutex
	delay    time.Duration
	timer    *time.Timer
	callback func()
	stopChan chan struct{}
}

// NewDebouncer creates a new debouncer
func NewDebouncer(delay time.Duration, callback func()) *Debouncer {
	return &Debouncer{
		delay:    delay,
		callback: callback,
		stopChan: make(chan struct{}),
	}
}

// Trigger triggers the debounced callback
func (d *Debouncer) Trigger() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Reset timer if it exists
	if d.timer != nil {
		d.timer.Stop()
	}

	// Create new timer
	d.timer = time.AfterFunc(d.delay, func() {
		select {
		case <-d.stopChan:
			return
		default:
			d.callback()
		}
	})
}

// Stop stops the debouncer
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}

	close(d.stopChan)
}

// Flush immediately executes the callback if pending
func (d *Debouncer) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.callback()
	}
}
