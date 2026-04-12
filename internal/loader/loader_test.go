package loader

import (
	"testing"
	"time"
)

func TestNewReturnsLoader(t *testing.T) {
	l := New(Dots)
	if l == nil {
		t.Fatal("New returned nil")
	}
	if l.stopped {
		t.Error("new loader should not be stopped")
	}
	if l.started {
		t.Error("new loader should not be started")
	}
}

func TestStopWithoutStart(t *testing.T) {
	l := New(Dots)
	// Stop without Start should not panic.
	l.Stop()
	if !l.stopped {
		t.Error("loader should be stopped after Stop()")
	}
}

func TestDoubleStop(t *testing.T) {
	l := New(Dots)
	l.Stop()
	l.Stop() // must not panic or double-close channel
	if !l.stopped {
		t.Error("loader should still be stopped")
	}
}

func TestStartStop(t *testing.T) {
	l := New(SpinnerStyle{
		Frames: []string{"."},
		Speed:  10 * time.Millisecond,
	})
	// Start may be a no-op if not a terminal, but must not panic.
	l.Start()
	time.Sleep(30 * time.Millisecond)
	l.Stop()

	if !l.stopped {
		t.Error("loader should be stopped after Stop()")
	}
}

func TestSetMessage(t *testing.T) {
	l := New(Dots)
	l.SetMessage("Loading")
	if l.style.Message != "Loading" {
		t.Errorf("message = %q, want %q", l.style.Message, "Loading")
	}
}

func TestStartAfterStopIsNoop(t *testing.T) {
	l := New(Dots)
	l.Stop()
	l.Start() // should be a no-op since already stopped
	// Verify no goroutine leak by checking stopped is still true.
	if !l.stopped {
		t.Error("Start after Stop should remain stopped")
	}
}
