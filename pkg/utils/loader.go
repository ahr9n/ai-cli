package utils

import (
	"fmt"
	"sync"
	"time"
)

type SpinnerStyle struct {
	Frames  []string
	Message string
	Color   string
	Speed   time.Duration
}

var (
	Dots = SpinnerStyle{
		Frames:  []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
		Speed:   100 * time.Millisecond,
		Message: "Thinking",
		Color:   "\033[32m",
	}
)

type Loader struct {
	stop    chan struct{}
	stopped bool
	mu      sync.Mutex
	style   SpinnerStyle
}

func NewLoader(style SpinnerStyle) *Loader {
	return &Loader{
		stop:  make(chan struct{}),
		style: style,
	}
}

func (l *Loader) SetMessage(message string) {
	l.mu.Lock()
	l.style.Message = message
	l.mu.Unlock()
}

func (l *Loader) Start() {
	l.mu.Lock()
	if l.stopped {
		l.mu.Unlock()
		return
	}
	l.mu.Unlock()

	resetColor := "\033[0m"

	go func() {
		ticker := time.NewTicker(l.style.Speed)
		defer ticker.Stop()

		i := 0
		for {
			select {
			case <-l.stop:
				return
			case <-ticker.C:
				l.mu.Lock()
				message := l.style.Message
				l.mu.Unlock()

				fmt.Print("\r\033[K")
				fmt.Printf("%s%s %s%s",
					l.style.Color,
					message,
					l.style.Frames[i%len(l.style.Frames)],
					resetColor,
				)
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
