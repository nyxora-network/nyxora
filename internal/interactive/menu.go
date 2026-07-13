package interactive

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type systemInfoMsg struct {
	cpuLoad    float64
	ramUsed    uint64
	ramTotal   uint64
	goroutines int
}

type viewState int

const (
	viewMenu viewState = iota
	viewConnect
	viewDashboard
	viewServerInfo
	viewInstall
	viewUpdate
	viewDisconnect
	viewHelp
	viewTopology
	viewQuit
)

type model struct {
	state         viewState
	choices       []string
	cursor        int
	width         int
	height        int
	quitting      bool
	booting       bool
	bootStep      int
	bootDone      bool
	tick          int
	spinner       spinner.Model
	theme         string
	showStatus    bool
	cpuLoad       float64
	ramUsed       uint64
	ramTotal      uint64
	goroutines    int
	activeTunnels int
	totalTunnels  int
	bestTunnel    string
	bestScore     float64
	pingMs        float64
	lossPercent   float64
	notification  int
	keyHint       string
	connectWizard *connectWizard
	helpModel     *helpModel
	topology      *topologyModel
}

func initialModel() model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Primary))
	s.Spinner = spinner.Moon

	return model{
		choices: []string{
			"C  Connect to Server",
			"D  Dashboard",
			"I  Server Info",
			"N  Install",
			"U  Check for Updates",
			"X  Disconnect",
			"T  Tunnel Topology",
			"H  Help",
			"Q  Exit",
		},
		cursor:        0,
		booting:       true,
		bootStep:      0,
		spinner:       s,
		theme:         "catppuccin-mocha",
		activeTunnels: 5,
		totalTunnels:  11,
		bestTunnel:    "hysteria",
		bestScore:     68.5,
		pingMs:        45.2,
		lossPercent:   0.2,
		notification:  1,
		keyHint:       "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Millisecond*60, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
		tickSystemCmd(),
		tea.EnterAltScreen,
		m.spinner.Tick,
	)
}

func tickSystemCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		load := 0.0
		if data, err := os.ReadFile("/proc/loadavg"); err == nil {
			fields := strings.Fields(string(data))
			if len(fields) >= 1 {
				load, _ = strconv.ParseFloat(fields[0], 64)
			}
		}

		var totalRAM uint64
		if data, err := os.ReadFile("/proc/meminfo"); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						kb, _ := strconv.ParseUint(fields[1], 10, 64)
						totalRAM = kb / 1024
					}
				}
			}
		}

		return systemInfoMsg{
			cpuLoad:    load,
			ramUsed:    mem.Alloc / 1024 / 1024,
			ramTotal:   totalRAM,
			goroutines: runtime.NumGoroutine(),
		}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == viewConnect && m.connectWizard != nil {
			return m.connectWizard.Update(msg)
		}
		if m.state == viewHelp && m.helpModel != nil {
			return m.helpModel.Update(msg)
		}
		if m.state == viewTopology && m.topology != nil {
			return m.topology.Update(msg)
		}
		return m.handleKeyMsg(msg)

	case tickMsg:
		m.tick++
		if m.booting && !m.bootDone {
			m.bootStep++
			if m.bootStep > 24 {
				m.bootDone = true
			}
		}
		return m, tea.Tick(time.Millisecond*60, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case systemInfoMsg:
		m.cpuLoad = msg.cpuLoad
		m.ramUsed = msg.ramUsed
		m.ramTotal = msg.ramTotal
		m.goroutines = msg.goroutines
		return m, tickSystemCmd()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.state == viewMenu {
			m.quitting = true
			return m, tea.Quit
		}
		m.state = viewMenu
		return m, nil

	case "enter":
		if m.booting && !m.bootDone {
			return m, nil
		}
		return m.handleMenuSelect()

	case "up", "k":
		if !m.booting || m.bootDone {
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
			m.updateKeyHint()
		}

	case "down", "j":
		if !m.booting || m.bootDone {
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}
			m.updateKeyHint()
		}

	case "1":
		m.theme = "catppuccin-mocha"
		currentTheme = themes["catppuccin-mocha"]
	case "2":
		m.theme = "tokyo-night"
		currentTheme = themes["tokyo-night"]
	case "3":
		m.theme = "catppuccin-latte"
		currentTheme = themes["catppuccin-latte"]

	case "s":
		m.showStatus = !m.showStatus

	case "?":
		m.state = viewHelp
		m.helpModel = newHelpModel()

	case "t", "T":
		m.state = viewTopology
		m.topology = newTopologyModel()
	}

	return m, nil
}

func (m *model) handleMenuSelect() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0:
		m.state = viewConnect
		m.connectWizard = newConnectWizard()
	case 1:
		m.state = viewDashboard
		return m, tea.Quit
	case 2:
		m.state = viewServerInfo
	case 3:
		m.state = viewInstall
	case 4:
		m.state = viewUpdate
	case 5:
		m.state = viewDisconnect
	case 6:
		m.state = viewTopology
		m.topology = newTopologyModel()
	case 7:
		m.state = viewHelp
		m.helpModel = newHelpModel()
	case 8:
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *model) updateKeyHint() {
	keyMap := map[int]string{
		0: "Connect to a remote server",
		1: "Open live dashboard",
		2: "View server info",
		3: "Install dependencies",
		4: "Check for updates",
		5: "Disconnect all tunnels",
		6: "View tunnel topology",
		7: "Help & keyboard shortcuts",
		8: "Exit NYXORA",
	}
	m.keyHint = keyMap[m.cursor]
}

func (m model) View() string {
	if m.quitting {
		return m.quittingView()
	}

	switch m.state {
	case viewConnect:
		if m.connectWizard != nil {
			return m.connectWizard.View()
		}
	case viewHelp:
		if m.helpModel != nil {
			return m.helpModel.View()
		}
	case viewTopology:
		if m.topology != nil {
			return m.topology.View()
		}
	}

	if m.booting && !m.bootDone {
		return m.bootView()
	}
	return m.menuView()
}

func (m model) quittingView() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(renderGradient("  NYXORA", currentTheme.GradientA, currentTheme.GradientB))
	b.WriteString("\n\n")
	b.WriteString("  " + dimStyle().Render("Thank you for using NYXORA"))
	b.WriteString("\n")
	b.WriteString("  " + dimStyle().Render("https://t.me/NyxoraCore"))
	b.WriteString("\n\n")
	return b.String()
}

func (m model) bootView() string {
	var b strings.Builder
	b.WriteString("\n")

	logoGradient := renderGradientBold(logo[0], currentTheme.GradientA, currentTheme.GradientB)
	b.WriteString("  " + logoGradient + "\n")
	for _, line := range logo[1:] {
		b.WriteString("  " + renderGradient(line, currentTheme.GradientA, currentTheme.GradientB) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(primaryStyle().Render("  NYXORA"))
	b.WriteString("  " + dimStyle().Render("Adaptive Tunnel Orchestrator v0.2.0"))
	b.WriteString("\n\n")

	barWidth := 42
	bar := renderGradientBar(float64(m.bootStep)*100/24, barWidth, currentTheme.GradientA, currentTheme.GradientB)
	pct := m.bootStep * 100 / 24

	steps := []string{
		"Initializing system...",
		"Loading transport modules...",
		"Checking dependencies...",
		"Configuring WireGuard...",
		"Setting up scoring engine...",
		"Preparing multipath scheduler...",
		"Configuring failover engine...",
		"Loading dashboard...",
		"Optimizing routes...",
		"Ready!",
	}

	stepIdx := m.bootStep * len(steps) / 24
	if stepIdx >= len(steps) {
		stepIdx = len(steps) - 1
	}

	spinnerChar := m.spinner.View()
	b.WriteString(fmt.Sprintf("  %s  %s\n", bar, dimStyle().Render(fmt.Sprintf("%d%%", pct))))
	b.WriteString(fmt.Sprintf("  %s %s %s\n",
		primaryStyle().Render(spinnerChar),
		warningStyle().Render("▸"),
		steps[stepIdx]))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  \u200Ehttps://t.me/NyxoraCore"))
	b.WriteString("\n")

	return b.String()
}

func (m model) menuView() string {
	var b strings.Builder

	title := renderGradient("  NYXORA", currentTheme.GradientA, currentTheme.GradientB)
	version := dimStyle().Render("v0.2.0")

	header := fmt.Sprintf("%s %s", title, version)
	separator := dimStyle().Render(strings.Repeat("─", 52))

	content := header + "\n" + separator + "\n\n"

	if m.showStatus {
		content += m.renderStatusBar() + "\n"
	}

	content += m.renderSystemInfo() + "\n"

	for i, choice := range m.choices {
		content += m.renderMenuItem(i, choice)
	}

	content += "\n" + dimStyle().Render(strings.Repeat("─", 52)) + "\n"

	if m.keyHint != "" {
		hintBox := subtleBox().
			Width(50).
			Render("  " + m.keyHint)
		content += hintBox + "\n"
	}

	navHelp := dimStyle().Render("  \u2191\u2193 navigate  \u23CE select  1/2/3 theme  s status  ? help  q quit")
	telegram := dimStyle().Render("  \u200Ehttps://t.me/NyxoraCore")

	content += navHelp + "\n" + telegram + "\n"

	menuBox := boxStyle().
		Width(m.width - 2).
		Render(content)

	b.WriteString(menuBox)
	b.WriteString("\n")

	return b.String()
}

func (m model) renderStatusBar() string {
	var b strings.Builder

	lossStyle := successStyle()
	if m.lossPercent > 5 {
		lossStyle = warningStyle()
	}
	if m.lossPercent > 20 {
		lossStyle = errorStyle()
	}

	b.WriteString(fmt.Sprintf("  %s %s  %s  %s %s\n",
		infoStyle().Render("Active:"),
		successStyle().Render(fmt.Sprintf("%d/%d", m.activeTunnels, m.totalTunnels)),
		dimStyle().Render("│"),
		infoStyle().Render("Best:"),
		successStyle().Render(m.bestTunnel),
	))
	b.WriteString(fmt.Sprintf("  %s %s  %s\n",
		successStyle().Render(fmt.Sprintf("Score: %.1f", m.bestScore)),
		dimStyle().Render("│"),
		lossStyle.Render(fmt.Sprintf("Ping: %.0fms  Loss: %.1f%%", m.pingMs, m.lossPercent)),
	))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderSystemInfo() string {
	ramPercent := 0.0
	if m.ramTotal > 0 {
		ramPercent = float64(m.ramUsed) / float64(m.ramTotal) * 100
	}

	cpuStyle := successStyle()
	if m.cpuLoad > 2.0 {
		cpuStyle = warningStyle()
	}
	if m.cpuLoad > 4.0 {
		cpuStyle = errorStyle()
	}

	ramStyle := successStyle()
	if ramPercent > 70 {
		ramStyle = warningStyle()
	}
	if ramPercent > 90 {
		ramStyle = errorStyle()
	}

	cpuBar := renderProgressBar(int(m.cpuLoad/8*20), 20)
	ramBar := renderProgressBar(int(ramPercent/100*20), 20)

	return fmt.Sprintf("  %s %s\n  %s %s\n  %s %s\n",
		dimStyle().Render("CPU:"),
		cpuStyle.Render(fmt.Sprintf("%.1f", m.cpuLoad)),
		dimStyle().Render("    "),
		cpuBar,
		dimStyle().Render("RAM:"),
		ramStyle.Render(fmt.Sprintf("%.0f%%  ", ramPercent)),
	) + fmt.Sprintf("  %s %s  %s  %s %s\n",
		dimStyle().Render("RAM:"),
		ramBar,
		dimStyle().Render("│"),
		dimStyle().Render("Go:"),
		dimStyle().Render(strconv.Itoa(m.goroutines)),
	)
}

func (m model) renderMenuItem(i int, choice string) string {
	isSelected := m.cursor == i

	cursor := "  "

	if isSelected {
		cursor = successStyle().Render("▸ ")
	}

	if len(choice) < 3 {
		return ""
	}

	key := choice[0]
	label := choice[3:]

	keyBadge := badgeStyle(string(key)).
		Background(lipgloss.Color(currentTheme.Muted)).
		Foreground(lipgloss.Color(currentTheme.Text))

	if isSelected {
		keyBadge = badgeStyle(string(key)).
			Background(lipgloss.Color(currentTheme.Highlight)).
			Foreground(lipgloss.Color(currentTheme.Bg))
	}

	keyBadgeStr := keyBadge.Render(string(key))
	labelStr := dimStyle().Render(label)

	if isSelected {
		labelStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color(currentTheme.Highlight)).
			Bold(true).
			Render(label)
	}

	return fmt.Sprintf("%s%s %s\n",
		cursor,
		keyBadgeStr,
		labelStr,
	)
}

// RunTransportStatus displays an interactive table of transport statuses in a Bubble Tea TUI.
func RunTransportStatus(transports []TransportStatus) error {
	p := tea.NewProgram(transportStatusModel{transports: transports}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type transportStatusModel struct {
	transports []TransportStatus
	cursor     int
	quitting   bool
	tick       int
}

type TransportStatus struct {
	Name    string
	Port    int
	Status  string
	Score   float64
	Latency float64
	Loss    float64
}

func (m transportStatusModel) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	)
}

func (m transportStatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.transports) - 1
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.transports) {
				m.cursor = 0
			}
		case "1", "2", "3":
			switch msg.String() {
			case "1":
				currentTheme = themes["catppuccin-mocha"]
			case "2":
				currentTheme = themes["tokyo-night"]
			case "3":
				currentTheme = themes["catppuccin-latte"]
			}
		}
	case tickMsg:
		m.tick++
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

func (m transportStatusModel) View() string {
	if m.quitting {
		return "\n"
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(renderGradient("  NYXORA Transports", currentTheme.GradientA, currentTheme.GradientB))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 62)))
	b.WriteString("\n\n")

	header := fmt.Sprintf("  %-2s %-12s %-6s %-8s %-6s %-8s %-6s %s",
		"#", "NAME", "PORT", "STATUS", "SCORE", "LATENCY", "LOSS", "BAR")
	b.WriteString(dimStyle().Bold(true).Render(header))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  " + strings.Repeat("─", 60)))
	b.WriteString("\n")

	for i, t := range m.transports {
		b.WriteString(m.renderTransportRow(i, t))
	}

	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 62)))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  \u2191\u2195 navigate  q/esc back  1/2/3 theme"))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  \u200Ehttps://t.me/NyxoraCore"))
	b.WriteString("\n")

	return b.String()
}

func (m transportStatusModel) renderTransportRow(i int, t TransportStatus) string {
	cursor := "  "
	nameStyle := dimStyle()
	if m.cursor == i {
		cursor = successStyle().Render("▸ ")
		nameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(currentTheme.Highlight)).
			Bold(true)
	}

	statusIcon := statusIcons["idle"]
	switch t.Status {
	case "active":
		statusIcon = statusIcons["active"]
	case "testing":
		statusIcon = statusIcons["testing"]
	case "failed":
		statusIcon = statusIcons["failed"]
	}

	scoreStyle := errorStyle()
	if t.Score >= 70 {
		scoreStyle = successStyle()
	} else if t.Score >= 40 {
		scoreStyle = warningStyle()
	}

	gradientBar := renderGradientBar(t.Score, 15, currentTheme.GradientA, currentTheme.GradientB)

	return fmt.Sprintf("%s%s %-12s %6d %s%-8s %s   %6.1fms %4.1f%% %s\n",
		cursor,
		dimStyle().Render(fmt.Sprintf("%2d", i+1)),
		nameStyle.Render(t.Name),
		t.Port,
		statusIcon+" ",
		t.Status,
		scoreStyle.Render(fmt.Sprintf("%5.1f", t.Score)),
		t.Latency,
		t.Loss,
		gradientBar,
	)
}

// RunMenu launches the full interactive Bubble Tea TUI.
// Returns the selected menu index and any error encountered.
func RunMenu() (int, error) {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return -1, err
	}
	result := m.(model)
	return result.cursor, nil
}
