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
	ITALIC   = "\033[3m"
	RESET    = "\033[0m"
	CLEARLN  = "\033[2K"
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

func trueColor(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

func trueColorBG(r, g, b int) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

var (
	catppuccinBase   = trueColor(203, 166, 247)
	catppuccinMauve  = trueColor(137, 180, 250)
	catppuccinGreen  = trueColor(166, 227, 161)
	catppuccinYellow = trueColor(249, 226, 175)
	catppuccinRed    = trueColor(243, 139, 168)
	catppuccinTeal   = trueColor(148, 226, 213)
	catppuccinText   = trueColor(205, 214, 244)
	catppuccinSub    = trueColor(108, 112, 134)
	catppuccinSurf   = trueColor(49, 50, 68)
	catppuccinMantl  = trueColor(30, 30, 46)
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
		_, _ = fmt.Sscanf(string(output), "%d %d", &height, &width)
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

	fmt.Print(HIDE + HOME + CLEARLN)

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
		_, _ = fmt.Sscanf(string(output), "%d %d", &t.height, &t.width)
	}
}

func center(s string, width int) string {
	clean := stripANSICodes(s)
	padding := (width - len(clean)) / 2
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
		return catppuccinGreen
	} else if score >= 40 {
		return catppuccinYellow
	}
	return catppuccinRed
}

func scoreBar(score float64, width int) string {
	filled := int((score / 100) * float64(width))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat(BAR_CHAR, filled) + strings.Repeat("─", width-filled)
	return scoreColor(score) + bar + RESET
}

func stripANSICodes(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' && i+1 < len(s) && s[i+1] == '[' {
			for j := i + 2; j < len(s); j++ {
				if s[j] == 'm' || s[j] == 'H' || s[j] == 'J' || s[j] == 'K' || s[j] == 'l' || s[j] == 'h' {
					i = j
					break
				}
			}
		} else {
			result.WriteByte(s[i])
		}
	}
	return result.String()
}

func topBorder(width int) string {
	return catppuccinSub + "┌" + strings.Repeat("─", width-2) + "┐" + RESET + "\n"
}

func bottomBorder(width int) string {
	return catppuccinSub + "└" + strings.Repeat("─", width-2) + "┘" + RESET + "\n"
}

func sepLine(width int) string {
	return catppuccinSub + "├" + strings.Repeat("─", width-2) + "┤" + RESET + "\n"
}
