package net

import (
	"context"
	"sync"
	"time"
)

// IPWatcher monitors network changes and detects IP address changes
type IPWatcher struct {
	CurrentIP string
	Interval  time.Duration
	OnChange  func(oldIP, newIP string)

	mu      sync.RWMutex
	stopCh  chan struct{}
	stopped bool
}

// NewIPWatcher creates a new IP watcher with the specified check interval
// Default interval is 5 seconds if interval is 0
func NewIPWatcher(interval time.Duration) *IPWatcher {
	if interval == 0 {
		interval = 5 * time.Second
	}

	return &IPWatcher{
		Interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins monitoring for IP address changes
// It uses the provided context for graceful shutdown
func (w *IPWatcher) Start(ctx context.Context) error {
	w.mu.Lock()
	if w.stopped {
		w.mu.Unlock()
		return nil
	}
	w.mu.Unlock()

	// Detect initial IP
	netInfo, err := DetectLocalIP()
	if err != nil {
		return err
	}

	w.mu.Lock()
	w.CurrentIP = netInfo.IP
	w.mu.Unlock()

	// Start monitoring loop
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.stopCh:
			return nil
		case <-ticker.C:
			if err := w.checkIPChange(); err != nil {
				// Continue monitoring even if detection fails
				continue
			}
		}
	}
}

// Stop stops the IP watcher
func (w *IPWatcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.stopped {
		return
	}

	w.stopped = true
	close(w.stopCh)
}

// checkIPChange detects if the IP address has changed and triggers the callback
func (w *IPWatcher) checkIPChange() error {
	netInfo, err := DetectLocalIP()
	if err != nil {
		return err
	}

	w.mu.Lock()
	oldIP := w.CurrentIP
	newIP := netInfo.IP
	w.mu.Unlock()

	if oldIP != newIP {
		w.mu.Lock()
		w.CurrentIP = newIP
		w.mu.Unlock()

		if w.OnChange != nil {
			w.OnChange(oldIP, newIP)
		}
	}

	return nil
}

// GetCurrentIP returns the current IP address (thread-safe)
func (w *IPWatcher) GetCurrentIP() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.CurrentIP
}
