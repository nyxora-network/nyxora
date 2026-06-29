package dashboard

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
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

func (t *TUI) render() {
	t.ensureSize()
	t.mu.Lock()
	provider := t.provider
	t.mu.Unlock()

	out := strings.Builder{}
	out.WriteString(HOME)

	if provider == nil {
		out.WriteString(t.center(" NYXORA", t.width))
		out.WriteString("\n\n")
		out.WriteString(t.center("initializing...", t.width))
		fmt.Print(out.String())
		return
	}

	status := provider.Status()

	// ── Header ──
	header := fmt.Sprintf(" %s%sNYXORA%s  %s%s%s  %s%s%s",
		PURPLE+BOLD, DOT, RESET,
		DIM+GRAY, "Adaptive Tunnel Orchestrator", RESET,
		DIM, "v0.1.0", RESET)
	out.WriteString(header + "\n")
	out.WriteString(fmt.Sprintf(" %s%s%s\n", DIM+GRAY, strings.Repeat("━", t.width-2), RESET))

	// ── Status Bar ──
	running, _ := status["running"].(bool)
	connected, _ := status["connected"].(bool)
	_, _ = status["phase"].(string)
	active, _ := status["active_transport"].(string)
	nodeID, _ := status["node_id"].(string)
	uptime, _ := status["uptime"].(string)

	statusIcon := GRAY + "●"
	statusLabel := "idle"
	switch {
	case connected:
		statusIcon = GREEN + "●"
		statusLabel = "connected"
	case running:
		statusIcon = YELLOW + "●"
		statusLabel = "running"
	}

	out.WriteString(fmt.Sprintf(" %s%s %s %s  %s%s %s  %s%s %s\n",
		CYAN+BOLD, "STATUS", RESET,
		statusIcon+" "+statusLabel,
		TEAL, "NODE", RESET,
		truncate(nodeID, 12),
		BLUE, "UP", uptime,
	))

	if connected {
		out.WriteString(fmt.Sprintf(" %s%s %s  %s%s %s\n",
			GREEN, "TUNNEL", RESET,
			active, DIM+GRAY, "(click for details)", RESET,
		))
	}

	// ── Remote Host Info ──
	if remote, ok := status["remote"].(map[string]interface{}); ok {
		hostname, _ := remote["hostname"].(string)
		addr, _ := remote["address"].(string)
		osInfo, _ := remote["os"].(string)
		arch, _ := remote["arch"].(string)

		out.WriteString(fmt.Sprintf("\n %s%sREMOTE HOST%s\n", BOLD, CYAN, RESET))
		out.WriteString(fmt.Sprintf("   %s%s %s%s\n", BOLD, hostname, GRAY, addr, RESET))
		out.WriteString(fmt.Sprintf("   %s%s  %s%s  %s\n", DIM, osInfo, arch, RESET))
	}

	// ── Steps Wizard ──
	if stepsRaw, ok := status["steps"].([]interface{}); ok && len(stepsRaw) > 0 {
		out.WriteString(fmt.Sprintf("\n %s%sSETUP STEPS%s\n", BOLD, MAGENTA, RESET))
		for _, sRaw := range stepsRaw {
			if s, ok := sRaw.(map[string]interface{}); ok {
				name, _ := s["name"].(string)
				stat, _ := s["status"].(string)
				detail, _ := s["detail"].(string)
				done, _ := s["done"].(bool)

				icon := DIM + "○" + RESET
				color := GRAY
				switch stat {
				case "OK":
					icon = GREEN + CHECK + RESET
					color = GREEN
				case "FAILED":
					icon = RED + CROSS + RESET
					color = RED
				case "RUNNING":
					icon = YELLOW + "◉" + RESET
					color = YELLOW
				case "WARN":
					icon = ORANGE + "△" + RESET
					color = ORANGE
				}

				detailStr := ""
				if detail != "" && done {
					detailStr = fmt.Sprintf(" %s%s%s", DIM+GRAY, detail, RESET)
				}
				out.WriteString(fmt.Sprintf("   %s %s%s%s%s\n", icon, color, name, RESET, detailStr))
			}
		}
	}

	// ── Transports Table ──
	if transportsRaw, ok := status["transports"].([]interface{}); ok && len(transportsRaw) > 0 {
		out.WriteString(fmt.Sprintf("\n %s%sTRANSPORTS%s\n", BOLD, CYAN, RESET))
		out.WriteString(fmt.Sprintf(" %s%-12s %-6s %-7s %-6s %-8s %-6s %s%s\n",
			DIM+GRAY, "NAME", "TYPE", "STATUS", "SCORE", "LATENCY", "LOSS", "BAR", RESET))
		out.WriteString(fmt.Sprintf(" %s%s%s\n", DIM+GRAY, strings.Repeat("─", t.width-4), RESET))

		var list []map[string]interface{}
		for _, r := range transportsRaw {
			if t, ok := r.(map[string]interface{}); ok {
				list = append(list, t)
			}
		}
		sort.Slice(list, func(i, j int) bool {
			si, _ := list[i]["score"].(float64)
			sj, _ := list[j]["score"].(float64)
			return si > sj
		})

		for _, tr := range list {
			name, _ := tr["name"].(string)
			typ, _ := tr["type"].(string)
			stat, _ := tr["status"].(string)
			score, _ := tr["score"].(float64)
			latency, _ := tr["latency"].(float64)
			loss, _ := tr["loss"].(float64)

			sColor := GRAY
			sIcon := "○"
			switch stat {
			case "active":
				sColor = GREEN
				sIcon = "●"
			case "failed":
				sColor = RED
				sIcon = "✗"
			default:
				sColor = GRAY
				sIcon = "○"
			}

			scColor := RED
			if score >= 70 {
				scColor = GREEN
			} else if score >= 40 {
				scColor = YELLOW
			}

			barLen := 12
			filled := int((score / 100) * float64(barLen))
			if filled > barLen {
				filled = barLen
			}
			bar := strings.Repeat(BAR_CHAR, filled) + strings.Repeat("─", barLen-filled)

			marker := "  "
			if name == active {
				marker = GREEN + "◀ " + RESET
			}

			out.WriteString(fmt.Sprintf(" %s%-10s %s %-5s %s%s%-5s %s%5.1f %s%6.1fms %s%4.1f%% %s%s%s\n",
				marker,
				BOLD+name+RESET,
				DIM+typ+RESET,
				sColor+sIcon+RESET,
				scColor, fmt.Sprintf("%.1f", score), RESET,
				latency,
				loss,
				scColor, bar, RESET,
			))
		}
	}

	// ── Failover ──
	if failoverRaw, ok := status["failover"].(map[string]interface{}); ok && len(failoverRaw) > 0 {
		out.WriteString(fmt.Sprintf("\n %s%sFAILOVER%s\n", BOLD, MAGENTA, RESET))
		var names []string
		for k := range failoverRaw {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, name := range names {
			val, _ := failoverRaw[name].(string)
			fColor := GREEN
			switch val {
			case "degraded":
				fColor = YELLOW
			case "down":
				fColor = RED
			}
			out.WriteString(fmt.Sprintf("   %s: %s%s%s\n", name, fColor, val, RESET))
		}
	}

	// ── Help Footer ──
	out.WriteString(fmt.Sprintf("\n %s%s", DIM+GRAY, strings.Repeat("━", t.width-2)))
	out.WriteString(fmt.Sprintf("\n %s%s nyxora connect <ip> --user root --port 22 %s", DIM, ARROW, RESET))
	out.WriteString(fmt.Sprintf("\n %s%s ctrl+c to exit %s", DIM, ARROW, RESET))

	fmt.Print(out.String())
}

func (t *TUI) center(s string, width int) string {
	padding := (width - len(s)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + s
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + ".."
}
