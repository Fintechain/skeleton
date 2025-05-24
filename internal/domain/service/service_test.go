package service

import (
	"testing"
)

// Test that constants have expected values
func TestServiceConstants(t *testing.T) {
	// Service status constants
	if StatusStopped != "stopped" {
		t.Errorf("StatusStopped = %v, want stopped", StatusStopped)
	}

	if StatusStarting != "starting" {
		t.Errorf("StatusStarting = %v, want starting", StatusStarting)
	}

	if StatusRunning != "running" {
		t.Errorf("StatusRunning = %v, want running", StatusRunning)
	}

	if StatusStopping != "stopping" {
		t.Errorf("StatusStopping = %v, want stopping", StatusStopping)
	}

	if StatusFailed != "failed" {
		t.Errorf("StatusFailed = %v, want failed", StatusFailed)
	}

	// Service error constants
	if ErrServiceStart != "service.start_failed" {
		t.Errorf("ErrServiceStart = %v, want service.start_failed", ErrServiceStart)
	}

	if ErrServiceStop != "service.stop_failed" {
		t.Errorf("ErrServiceStop = %v, want service.stop_failed", ErrServiceStop)
	}

	if ErrServiceNotFound != "service.not_found" {
		t.Errorf("ErrServiceNotFound = %v, want service.not_found", ErrServiceNotFound)
	}
}
