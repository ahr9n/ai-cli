package utils

import (
	"fmt"
	"sync"
	"time"
)

type Loader struct {
	stop    chan struct{}
	stopped bool
	mu      sync.Mutex
}

func NewLoader() *Loader {
	return &Loader{
		stop: make(chan struct{}),
	}
}

func (l *Loader) Start() {
	l.mu.Lock()
	if l.stopped {
		l.mu.Unlock()
		return
	}
	l.mu.Unlock()

	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()

		i := 0
		fmt.Print("\r")
		for {
			select {
			case <-l.stop:
				return
			case <-ticker.C:
				fmt.Print("\r\033[K")
				fmt.Printf("thinking %s", frames[i%len(frames)])
				i++
			}
		}
	}()
}

func (l *Loader) Stop() {
	l.mu.Lock()
	if !l.stopped {
		l.stopped = true
		close(l.stop)
	}
	l.mu.Unlock()
	fmt.Print("\r\033[K")
}
