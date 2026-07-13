package dashboard

import (
	"fmt"
	"sort"
	"strings"
)

func (t *TUI) render() {
	t.ensureSize()
	t.mu.Lock()
	provider := t.provider
	t.mu.Unlock()

	out := strings.Builder{}
	out.WriteString(HOME + CLEARLN)

	if provider == nil {
		out.WriteString(topBorder(t.width))
		out.WriteString(center(catppuccinBase+BOLD+" NYXORA"+RESET, t.width) + "\n")
		out.WriteString(center(catppuccinSub+DIM+"initializing..."+RESET, t.width) + "\n")
		out.WriteString(bottomBorder(t.width))
		fmt.Print(out.String())
		return
	}

	status := provider.Status()

	out.WriteString(topBorder(t.width))
	out.WriteString(t.renderHeader(status))
	out.WriteString(sepLine(t.width))
	out.WriteString(t.renderStatusBar(status))
	out.WriteString(t.renderRemoteHost(status))
	out.WriteString(t.renderSteps(status))
	out.WriteString(t.renderTransports(status))
	out.WriteString(t.renderFailover(status))
	out.WriteString(bottomBorder(t.width))

	fmt.Print(out.String())
}

func (t *TUI) renderHeader(status map[string]interface{}) string {
	header := fmt.Sprintf(" %s%s●%s  %s%s%s %s%s%s",
		catppuccinBase+BOLD, DOT, RESET,
		catppuccinBase+BOLD, "NYXORA", RESET,
		catppuccinSub+DIM, "Adaptive Tunnel Orchestrator", RESET)
	return header + "\n"
}

func (t *TUI) renderStatusBar(status map[string]interface{}) string {
	running, _ := status["running"].(bool)
	connected, _ := status["connected"].(bool)
	active, _ := status["active_transport"].(string)
	nodeID, _ := status["node_id"].(string)
	uptime, _ := status["uptime"].(string)

	statusIcon := catppuccinSub + "●"
	statusLabel := "idle"
	switch {
	case connected:
		statusIcon = catppuccinGreen + "●"
		statusLabel = "connected"
	case running:
		statusIcon = catppuccinYellow + "●"
		statusLabel = "running"
	}

	result := fmt.Sprintf(" %s%s %s %s  %s%s %s%s  %s%s %s%s\n",
		catppuccinTeal+BOLD, "STATUS", RESET,
		statusIcon+" "+statusLabel,
		catppuccinMauve, "NODE", RESET,
		truncateStr(nodeID, 12),
		BLUE+BOLD, "UP", RESET,
		uptime,
	)

	if connected {
		result += fmt.Sprintf(" %s%s %s  %s%s %s%s\n",
			catppuccinGreen, "TUNNEL", RESET,
			active, catppuccinSub+DIM, "(details)", RESET,
		)
	}

	return result
}

func (t *TUI) renderRemoteHost(status map[string]interface{}) string {
	remote, ok := status["remote"].(map[string]interface{})
	if !ok {
		return ""
	}
	hostname, _ := remote["hostname"].(string)
	addr, _ := remote["address"].(string)
	osInfo, _ := remote["os"].(string)
	arch, _ := remote["arch"].(string)

	result := fmt.Sprintf("\n %s%sREMOTE HOST%s\n", BOLD, catppuccinMauve, RESET)
	result += fmt.Sprintf("   %s%s %s%s%s\n", BOLD, hostname, catppuccinSub, addr, RESET)
	result += fmt.Sprintf("   %s%s  %s%s\n", DIM, osInfo, arch, RESET)
	return result
}

func (t *TUI) renderSteps(status map[string]interface{}) string {
	stepsRaw, ok := status["steps"].([]interface{})
	if !ok || len(stepsRaw) == 0 {
		return ""
	}

	result := fmt.Sprintf("\n %s%sSETUP STEPS%s\n", BOLD, catppuccinMauve, RESET)
	for _, sRaw := range stepsRaw {
		s, ok := sRaw.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := s["name"].(string)
		stat, _ := s["status"].(string)
		detail, _ := s["detail"].(string)
		done, _ := s["done"].(bool)

		icon := catppuccinSub + "○" + RESET
		color := catppuccinSub
		switch stat {
		case "OK":
			icon = catppuccinGreen + CHECK + RESET
			color = catppuccinGreen
		case "FAILED":
			icon = catppuccinRed + CROSS + RESET
			color = catppuccinRed
		case "RUNNING":
			icon = catppuccinYellow + "◉" + RESET
			color = catppuccinYellow
		case "WARN":
			icon = ORANGE + "△" + RESET
			color = ORANGE
		}

		detailStr := ""
		if detail != "" && done {
			detailStr = fmt.Sprintf(" %s%s%s", DIM+catppuccinSub, detail, RESET)
		}
		result += fmt.Sprintf("   %s %s%s%s%s\n", icon, color, name, RESET, detailStr)
	}
	return result
}

func (t *TUI) renderTransports(status map[string]interface{}) string {
	transportsRaw, ok := status["transports"].([]interface{})
	if !ok || len(transportsRaw) == 0 {
		return ""
	}

	active, _ := status["active_transport"].(string)

	result := fmt.Sprintf("\n %s%sTRANSPORTS%s\n", BOLD, catppuccinMauve, RESET)
	result += fmt.Sprintf(" %s%-12s %-6s %-7s %-6s %-8s %-6s %s%s\n",
		DIM+catppuccinSub, "NAME", "TYPE", "STATUS", "SCORE", "LATENCY", "LOSS", "BAR", RESET)
	result += fmt.Sprintf(" %s%s%s\n", DIM+catppuccinSub, strings.Repeat("─", t.width-4), RESET)

	var list []map[string]interface{}
	for _, r := range transportsRaw {
		if m, ok := r.(map[string]interface{}); ok {
			list = append(list, m)
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

		sColor := catppuccinSub
		sIcon := "○"
		switch stat {
		case "active":
			sColor = catppuccinGreen
			sIcon = "●"
		case "failed":
			sColor = catppuccinRed
			sIcon = "✗"
		case "testing":
			sColor = catppuccinYellow
			sIcon = "◉"
		}

		bar := scoreBar(score, 12)
		sc := scoreColor(score)

		marker := "  "
		if name == active {
			marker = catppuccinGreen + "◀ " + RESET
		}

		result += fmt.Sprintf(" %s%-10s %s %-5s %s%5.1f %6.1fms %4.1f%% %s%s%s\n",
			marker,
			BOLD+name+RESET,
			DIM+typ+RESET,
			sColor+sIcon+RESET,
			sc, score,
			latency, loss,
			sc, bar, RESET,
		)
	}
	return result
}

func (t *TUI) renderFailover(status map[string]interface{}) string {
	failoverRaw, ok := status["failover"].(map[string]interface{})
	if !ok || len(failoverRaw) == 0 {
		return ""
	}

	result := fmt.Sprintf("\n %s%sFAILOVER%s\n", BOLD, catppuccinMauve, RESET)
	var names []string
	for k := range failoverRaw {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, name := range names {
		val, _ := failoverRaw[name].(string)
		fColor := catppuccinGreen
		switch val {
		case "degraded":
			fColor = catppuccinYellow
		case "down":
			fColor = catppuccinRed
		}
		result += fmt.Sprintf("   %s: %s%s%s\n", name, fColor, val, RESET)
	}
	return result
}
