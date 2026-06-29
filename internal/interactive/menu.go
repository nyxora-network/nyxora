package interactive

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Theme colors
type Theme struct {
	Primary    string
	Secondary  string
	Accent     string
	Bg         string
	Success    string
	Warning    string
	Error      string
	Info       string
	Border     string
	Highlight  string
	Text       string
	TextDim    string
	GradientA  string
	GradientB  string
}

var themes = map[string]Theme{
	"dark": {
		Primary: "141", Secondary: "99", Accent: "212",
		Success: "46", Warning: "214", Error: "196", Info: "75",
		Border: "241", Highlight: "212", Text: "255", TextDim: "241",
		GradientA: "99", GradientB: "141",
	},
	"neon": {
		Primary: "201", Secondary: "51", Accent: "226",
		Success: "82", Warning: "208", Error: "196", Info: "45",
		Border: "141", Highlight: "226", Text: "255", TextDim: "240",
		GradientA: "201", GradientB: "51",
	},
	"light": {
		Primary: "57", Secondary: "93", Accent: "141",
		Success: "22", Warning: "130", Error: "160", Info: "62",
		Border: "250", Highlight: "57", Text: "234", TextDim: "250",
		GradientA: "57", GradientB: "93",
	},
}

var currentTheme Theme = themes["dark"]

type tickMsg time.Time

type systemInfoMsg struct {
	cpuLoad    float64
	ramUsed    uint64
	ramTotal   uint64
	goroutines int
}

type model struct {
	choices       []string
	cursor        int
	width         int
	height        int
	quitting      bool
	booting       bool
	bootStep      int
	bootDone      bool
	tick          int
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
}

func initialModel() model {
	return model{
		choices: []string{
			"C  Connect to Server",
			"D  Dashboard",
			"I  Server Info",
			"N  Install",
			"U  Check for Updates",
			"X  Disconnect",
			"Q  Exit",
		},
		cursor:        0,
		booting:       true,
		bootStep:      0,
		theme:         "dark",
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
	return tea.Batch(tickCmd(), tickSystemCmd(), tea.EnterAltScreen)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
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
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.booting && !m.bootDone {
				return m, nil
			}
			return m, tea.Quit
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
		case "1", "2", "3":
			m.cycleTheme()
		case "f1":
			m.theme = "dark"
			currentTheme = themes["dark"]
		case "f2":
			m.theme = "neon"
			currentTheme = themes["neon"]
		case "f3":
			m.theme = "light"
			currentTheme = themes["light"]
		case "s":
			m.showStatus = !m.showStatus
		}

	case tickMsg:
		m.tick++
		if m.booting && !m.bootDone {
			m.bootStep++
			if m.bootStep > 20 {
				m.bootDone = true
			}
		}
		return m, tickCmd()

	case systemInfoMsg:
		m.cpuLoad = msg.cpuLoad
		m.ramUsed = msg.ramUsed
		m.ramTotal = msg.ramTotal
		m.goroutines = msg.goroutines
		return m, tickSystemCmd()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *model) updateKeyHint() {
	keyMap := map[int]string{
		0: "Press Enter to connect to a remote server",
		1: "Press Enter to open live dashboard",
		2: "Press Enter to view server info",
		3: "Press Enter to install dependencies",
		4: "Press Enter to check for updates",
		5: "Press Enter to disconnect all tunnels",
		6: "Press Enter to exit NYXORA",
	}
	m.keyHint = keyMap[m.cursor]
}

func (m *model) cycleTheme() {
	themeNames := []string{"dark", "neon", "light"}
	for i, t := range themeNames {
		if t == m.theme {
			m.theme = themeNames[(i+1)%len(themeNames)]
			currentTheme = themes[m.theme]
			return
		}
	}
}

func (m model) View() string {
	if m.quitting {
		return m.quittingView()
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

	b.WriteString(renderGradient("  NYXORA", currentTheme.GradientA, currentTheme.GradientB))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  Adaptive Tunnel Orchestrator v0.2.0"))
	b.WriteString("\n\n")

	barWidth := 40
	progress := m.bootStep * barWidth / 20
	bar := renderProgressBar(progress, barWidth)
	pct := m.bootStep * 100 / 20

	steps := []string{
		"Initializing system...",
		"Loading transport modules...",
		"Checking dependencies...",
		"Setting up scoring engine...",
		"Preparing multipath scheduler...",
		"Configuring failover engine...",
		"Loading dashboard...",
		"Ready!",
	}

	stepIdx := m.bootStep * len(steps) / 20
	if stepIdx >= len(steps) {
		stepIdx = len(steps) - 1
	}

	b.WriteString(fmt.Sprintf("  %s %s\n", bar, dimStyle().Render(fmt.Sprintf("%d%%", pct))))
	b.WriteString(fmt.Sprintf("  %s %s\n", warningStyle().Render("▸"), steps[stepIdx]))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  📱 https://t.me/NyxoraCore"))
	b.WriteString("\n")

	return b.String()
}

func (m model) menuView() string {
	var b strings.Builder
	b.WriteString("\n")

	b.WriteString(renderGradient("  NYXORA", currentTheme.GradientA, currentTheme.GradientB))
	b.WriteString("  " + dimStyle().Render("v0.2.0"))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 56)))
	b.WriteString("\n\n")

	if m.showStatus {
		b.WriteString(m.renderStatusBar())
		b.WriteString("\n")
	}

	b.WriteString(m.renderSystemInfo())
	b.WriteString("\n")

	for i, choice := range m.choices {
		b.WriteString(m.renderMenuItem(i, choice))
	}

	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 56)))
	b.WriteString("\n")

	if m.keyHint != "" {
		b.WriteString(infoStyle().Render("  "+m.keyHint))
		b.WriteString("\n")
	}

	b.WriteString(dimStyle().Render("  ↑↓ navigate • enter select • 1/2/3 theme • s status • q quit"))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  📱 https://t.me/NyxoraCore"))
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

	b.WriteString(fmt.Sprintf("  %s %s %s %s\n",
		infoStyle().Render("Active:"),
		successStyle().Render(fmt.Sprintf("%d/%d", m.activeTunnels, m.totalTunnels)),
		dimStyle().Render("|"),
		infoStyle().Render("Best:"),
	))
	b.WriteString(fmt.Sprintf("  %s %s %s %s\n",
		successStyle().Render(m.bestTunnel),
		dimStyle().Render(fmt.Sprintf("(%.1f)", m.bestScore)),
		dimStyle().Render("|"),
		lossStyle.Render(fmt.Sprintf("Ping: %.0fms Loss: %.1f%%", m.pingMs, m.lossPercent)),
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

	return fmt.Sprintf("  %s %s %s %s %s %s %s %s\n",
		dimStyle().Render("CPU"),
		cpuStyle.Render(fmt.Sprintf("%.1f", m.cpuLoad)),
		dimStyle().Render("|"),
		dimStyle().Render("RAM"),
		ramStyle.Render(fmt.Sprintf("%.0f%%", ramPercent)),
		dimStyle().Render("|"),
		dimStyle().Render("Go"),
		dimStyle().Render(strconv.Itoa(m.goroutines)),
	)
}

func (m model) renderMenuItem(i int, choice string) string {
	cursor := "  "
	style := menuItemStyle()

	if m.cursor == i {
		cursor = successStyle().Render("▸ ")
		style = menuSelectedStyle()
	}

	icon := choice[:2]
	label := choice[3:]

	key := ""
	if len(choice) > 4 {
		parts := strings.SplitN(choice, "  ", 2)
		if len(parts) == 2 {
			icon = parts[0]
			key = parts[1][:1]
			label = parts[1][2:]
		}
	}

	if m.cursor == i {
		keyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(currentTheme.Highlight)).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(currentTheme.Highlight))
		return fmt.Sprintf("%s%s %s %s\n",
			cursor,
			style.Render(icon),
			style.Render(label),
			keyStyle.Render(key),
		)
	}

	return fmt.Sprintf("%s%s %s\n",
		cursor,
		dimStyle().Render(icon),
		dimStyle().Render(label),
	)
}

func renderGradient(text, colorA, colorB string) string {
	runes := []rune(text)
	n := len(runes)
	if n == 0 {
		return ""
	}

	var result strings.Builder
	for i, r := range runes {
		t := float64(i) / float64(n-1)
		r1, _ := strconv.Atoi(colorA)
		r2, _ := strconv.Atoi(colorB)
		mid := r1 + int(t*float64(r2-r1))
		if mid > 255 {
			mid = 255
		}
		if mid < 0 {
			mid = 0
		}
		result.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color(strconv.Itoa(mid))).
			Bold(true).
			Render(string(r)))
	}
	return result.String()
}

func renderProgressBar(filled, width int) string {
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return successStyle().Render(bar)
}

func successStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Success)).Bold(true)
}

func warningStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Warning)).Bold(true)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Error)).Bold(true)
}

func infoStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Info))
}

func dimStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.TextDim))
}

func menuItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1)
}

func menuSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.Highlight)).
		Bold(true).
		Padding(0, 1)
}

func RunMenu() (int, error) {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return -1, err
	}
	result := m.(model)
	return result.cursor, nil
}

type connectWizard struct {
	step       int
	addr       string
	user       string
	password   string
	mode       string
	transports string
	ports      string
	width      int
	height     int
	cursor     int
	quitting   bool
	tick       int
}

func connectModel() connectWizard {
	return connectWizard{step: 0, user: "root", mode: "auto"}
}

func (m connectWizard) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tickCmd())
}

func (m connectWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.step > 0 {
				m.step--
			} else {
				m.quitting = true
				return m, tea.Quit
			}
		case "enter":
			if m.step < 4 {
				m.step++
			} else {
				return m, tea.Quit
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 2
			}
		case "down", "j":
			m.cursor++
			if m.cursor > 2 {
				m.cursor = 0
			}
		}
	case tickMsg:
		m.tick++
		return m, tickCmd()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m connectWizard) View() string {
	if m.quitting {
		return "\n  Cancelled.\n\n"
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(renderGradient("  NYXORA Connect", currentTheme.GradientA, currentTheme.GradientB))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 50)))
	b.WriteString("\n\n")

	steps := []string{"Address", "User", "Password", "Mode", "Go!"}
	for i, s := range steps {
		icon := dimStyle().Render("○")
		style := dimStyle()
		if i < m.step {
			icon = successStyle().Render("●")
			style = successStyle()
		} else if i == m.step {
			icon = warningStyle().Render("◉")
			style = menuSelectedStyle()
		}
		b.WriteString(fmt.Sprintf("  %s %s\n", icon, style.Render(s)))
	}
	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 50)))
	b.WriteString("\n\n")

	switch m.step {
	case 0:
		b.WriteString(menuSelectedStyle().Render("  Remote server address:"))
		b.WriteString("\n\n")
		if m.addr == "" {
			b.WriteString(dimStyle().Render("  > _"))
		} else {
			b.WriteString(fmt.Sprintf("  > %s", m.addr))
		}
	case 1:
		b.WriteString(menuSelectedStyle().Render("  SSH username:"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("  > %s", m.user))
	case 2:
		b.WriteString(menuSelectedStyle().Render("  SSH password:"))
		b.WriteString("\n\n")
		b.WriteString("  > " + strings.Repeat("*", max(len(m.password), 8)))
	case 3:
		b.WriteString(menuSelectedStyle().Render("  Server mode:"))
		b.WriteString("\n\n")
		modes := []struct {
			name string
			desc string
			req  string
		}{
			{"full", "All 11 tunnels", "2GB+ RAM"},
			{"lite", "Lightweight", "512MB-2GB"},
			{"minimal", "SSH + SS only", "<512MB"},
		}
		for i, md := range modes {
			cursor := "  "
			style := dimStyle()
			if m.cursor == i {
				cursor = successStyle().Render("▸ ")
				style = menuSelectedStyle()
			}
			b.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, style.Render(md.name), dimStyle().Render(fmt.Sprintf("(%s) [%s]", md.desc, md.req))))
		}
	case 4:
		b.WriteString(successStyle().Render("  Ready to connect!"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("  Address:    %s\n", m.addr))
		b.WriteString(fmt.Sprintf("  User:       %s\n", m.user))
		b.WriteString(fmt.Sprintf("  Password:   %s\n", strings.Repeat("*", len(m.password))))
		b.WriteString(fmt.Sprintf("  Mode:       %s\n", m.mode))
	}

	b.WriteString("\n\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 50)))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  esc back • enter next • q cancel"))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  📱 https://t.me/NyxoraCore"))
	b.WriteString("\n")

	return b.String()
}

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
	return tea.Batch(tickCmd())
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
		}
	case tickMsg:
		m.tick++
		return m, tickCmd()
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

	header := fmt.Sprintf("  %-12s %-6s %-8s %-6s %-8s %-6s %s",
		"NAME", "PORT", "STATUS", "SCORE", "LATENCY", "LOSS", "BAR")
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
	b.WriteString(dimStyle().Render("  ↑↓ navigate • q/esc back • 1/2/3 theme"))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  📱 https://t.me/NyxoraCore"))
	b.WriteString("\n")

	return b.String()
}

func (m transportStatusModel) renderTransportRow(i int, t TransportStatus) string {
	cursor := "  "
	nameStyle := dimStyle()
	if m.cursor == i {
		cursor = successStyle().Render("▸ ")
		nameStyle = menuSelectedStyle()
	}

	statusIcon := dimStyle().Render("○")
	statusStyle := dimStyle()
	switch t.Status {
	case "active":
		statusIcon = successStyle().Render("●")
		statusStyle = successStyle()
	case "testing":
		statusIcon = warningStyle().Render("◉")
		statusStyle = warningStyle()
	case "failed":
		statusIcon = errorStyle().Render("✗")
		statusStyle = errorStyle()
	}

	scoreStyle := errorStyle()
	if t.Score >= 70 {
		scoreStyle = successStyle()
	} else if t.Score >= 40 {
		scoreStyle = warningStyle()
	}

	animatedBar := m.renderAnimatedBar(t.Score, t.Name)

	return fmt.Sprintf("%s%-12s %6d %s%-8s %s   %6.1fms %4.1f%% %s\n",
		cursor,
		nameStyle.Render(t.Name),
		t.Port,
		statusStyle.Render(statusIcon+" "),
		t.Status,
		scoreStyle.Render(fmt.Sprintf("%5.1f", t.Score)),
		t.Latency,
		t.Loss,
		animatedBar,
	)
}

func (m transportStatusModel) renderAnimatedBar(score float64, name string) string {
	width := 15
	filled := int((score / 100) * float64(width))
	if filled > width {
		filled = width
	}

	offset := 0
	for _, c := range name {
		offset += int(c)
	}
	phase := (m.tick + offset) % 4

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			if i == filled-1 && phase < 2 {
				bar += "▓"
			} else {
				bar += "█"
			}
		} else {
			bar += "░"
		}
	}

	style := errorStyle()
	if score >= 70 {
		style = successStyle()
	} else if score >= 40 {
		style = warningStyle()
	}

	return style.Render(bar)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(x float64) float64 {
	return math.Abs(x)
}
