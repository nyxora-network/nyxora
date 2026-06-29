package dashboard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	ESC      = "\033["
	BOLD     = "\033[1m"
	DIM      = "\033[2m"
	RESET    = "\033[0m"
	CLEAR    = "\033[2J"
	HOME     = "\033[H"
	HIDE     = "\033[?25l"
	SHOW     = "\033[?25h"

	BLACK   = "\033[30m"
	RED     = "\033[31m"
	GREEN   = "\033[32m"
	YELLOW  = "\033[33m"
	BLUE    = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN    = "\033[36m"
	WHITE   = "\033[37m"

	GRAY   = "\033[90m"
	ORANGE = "\033[38;5;214m"
	PURPLE = "\033[38;5;141m"
	TEAL   = "\033[38;5;80m"

	BAR_CHAR = "━"
	DOT      = "●"
	CHECK    = "✓"
	CROSS    = "✗"
	ARROW    = "➜"
)

type StatusProvider interface {
	Status() map[string]interface{}
}

type TUI struct {
	mu        sync.Mutex
	provider  StatusProvider
	interval  time.Duration
	running   bool
	stopCh    chan struct{}
	width     int
	height    int
	startTime time.Time
}

func NewTUI(intervalSec int) *TUI {
	width, height := 80, 24
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	if output, err := cmd.Output(); err == nil {
		fmt.Sscanf(string(output), "%d %d", &height, &width)
	}
	return &TUI{
		interval:  time.Duration(intervalSec) * time.Second,
		stopCh:    make(chan struct{}),
		width:     width,
		height:    height,
		startTime: time.Now(),
	}
}

func (t *TUI) SetProvider(p StatusProvider) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.provider = p
}

func (t *TUI) Start() error {
	t.mu.Lock()
	if t.running {
		t.mu.Unlock()
		return nil
	}
	t.running = true
	t.mu.Unlock()

	fmt.Print(HIDE + CLEAR)

	go func() {
		for {
			select {
			case <-t.stopCh:
				fmt.Print(SHOW + "\n")
				return
			default:
				t.render()
				time.Sleep(t.interval)
			}
		}
	}()
	return nil
}

func (t *TUI) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		close(t.stopCh)
		fmt.Print(SHOW)
		t.running = false
	}
}

func (t *TUI) ensureSize() {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	if output, err := cmd.Output(); err == nil {
		fmt.Sscanf(string(output), "%d %d", &t.height, &t.width)
	}
}

func center(s string, width int) string {
	padding := (width - len(s)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + s
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + ".."
}

func scoreColor(score float64) string {
	if score >= 70 {
		return GREEN
	} else if score >= 40 {
		return YELLOW
	}
	return RED
}

func scoreBar(score float64, width int) string {
	filled := int((score / 100) * float64(width))
	if filled > width {
		filled = width
	}
	return strings.Repeat(BAR_CHAR, filled) + strings.Repeat("─", width-filled)
}
