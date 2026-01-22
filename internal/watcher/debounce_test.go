package watcher

import (
	"sync"
	"testing"
	"time"
)

func TestNewDebouncer(t *testing.T) {
	callback := func() {
		// Test callback
	}

	debouncer := NewDebouncer(100*time.Millisecond, callback)

	if debouncer == nil {
		t.Fatal("NewDebouncer returned nil")
	}

	if debouncer.delay != 100*time.Millisecond {
		t.Errorf("Expected delay 100ms, got %v", debouncer.delay)
	}

	if debouncer.callback == nil {
		t.Error("Callback is nil")
	}

	// Clean up
	debouncer.Stop()
}

func TestDebouncer_SingleTrigger(t *testing.T) {
	var called bool
	var mu sync.Mutex

	callback := func() {
		mu.Lock()
		called = true
		mu.Unlock()
	}

	debouncer := NewDebouncer(50*time.Millisecond, callback)
	defer debouncer.Stop()

	debouncer.Trigger()

	// Wait for debounce delay plus a bit extra
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if !called {
		t.Error("Callback was not called after debounce delay")
	}
}

func TestDebouncer_MultipleTriggers(t *testing.T) {
	var callCount int
	var mu sync.Mutex

	callback := func() {
		mu.Lock()
		callCount++
		mu.Unlock()
	}

	debouncer := NewDebouncer(50*time.Millisecond, callback)
	defer debouncer.Stop()

	// Trigger multiple times rapidly
	for i := 0; i < 5; i++ {
		debouncer.Trigger()
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for debounce delay
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	// Should only be called once due to debouncing
	if callCount != 1 {
		t.Errorf("Expected callback to be called once, but was called %d times", callCount)
	}
}

func TestDebouncer_Stop(t *testing.T) {
	var called bool
	var mu sync.Mutex

	callback := func() {
		mu.Lock()
		called = true
		mu.Unlock()
	}

	debouncer := NewDebouncer(50*time.Millisecond, callback)

	debouncer.Trigger()
	debouncer.Stop()

	// Wait to see if callback gets called (it shouldn't)
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if called {
		t.Error("Callback was called after Stop()")
	}
}

func TestDebouncer_Flush(t *testing.T) {
	var called bool
	var mu sync.Mutex

	callback := func() {
		mu.Lock()
		called = true
		mu.Unlock()
	}

	debouncer := NewDebouncer(1*time.Second, callback) // Long delay
	defer debouncer.Stop()

	debouncer.Trigger()

	// Flush immediately without waiting for delay
	debouncer.Flush()

	mu.Lock()
	defer mu.Unlock()

	if !called {
		t.Error("Callback was not called after Flush()")
	}
}

func TestDebouncer_FlushWithoutTrigger(t *testing.T) {
	var callCount int
	var mu sync.Mutex

	callback := func() {
		mu.Lock()
		callCount++
		mu.Unlock()
	}

	debouncer := NewDebouncer(100*time.Millisecond, callback)
	defer debouncer.Stop()

	// Flush without triggering
	debouncer.Flush()

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if callCount != 0 {
		t.Errorf("Expected callback not to be called, but was called %d times", callCount)
	}
}

func TestDebouncer_ConcurrentTriggers(t *testing.T) {
	var callCount int
	var mu sync.Mutex

	callback := func() {
		mu.Lock()
		callCount++
		mu.Unlock()
	}

	debouncer := NewDebouncer(50*time.Millisecond, callback)
	defer debouncer.Stop()

	// Trigger from multiple goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			debouncer.Trigger()
		}()
	}

	wg.Wait()

	// Wait for debounce delay
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	// Should only be called once despite concurrent triggers
	if callCount != 1 {
		t.Errorf("Expected callback to be called once, but was called %d times", callCount)
	}
}

func TestDebouncer_ResetTimer(t *testing.T) {
	var called bool
	var mu sync.Mutex

	callback := func() {
		mu.Lock()
		called = true
		mu.Unlock()
	}

	debouncer := NewDebouncer(100*time.Millisecond, callback)
	defer debouncer.Stop()

	debouncer.Trigger()
	time.Sleep(50 * time.Millisecond)

	// Trigger again - should reset timer
	debouncer.Trigger()
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	notCalledYet := !called
	mu.Unlock()

	if !notCalledYet {
		t.Error("Callback was called before second debounce delay expired")
	}

	// Wait for second delay to complete
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if !called {
		t.Error("Callback was not called after second debounce delay")
	}
}
