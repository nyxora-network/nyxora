package interactive

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type connectWizard struct {
	step        int
	addrInput   textinput.Model
	userInput   textinput.Model
	passInput   textinput.Model
	mode        string
	transports  string
	width       int
	height      int
	cursor      int
	quitting    bool
	tick        int
	modeOptions []modeOption
}

type modeOption struct {
	name string
	desc string
	req  string
}

func newConnectWizard() *connectWizard {
	addr := textinput.New()
	addr.Placeholder = "192.168.1.100"
	addr.Focus()
	addr.CharLimit = 100
	addr.Width = 30
	addr.Prompt = "  > "

	user := textinput.New()
	user.Placeholder = "root"
	user.CharLimit = 50
	user.Width = 30
	user.Prompt = "  > "

	pass := textinput.New()
	pass.Placeholder = "password"
	pass.EchoMode = textinput.EchoPassword
	pass.EchoCharacter = '•'
	pass.CharLimit = 100
	pass.Width = 30
	pass.Prompt = "  > "

	return &connectWizard{
		step:      0,
		addrInput: addr,
		userInput: user,
		passInput: pass,
		mode:      "auto",
		modeOptions: []modeOption{
			{"full", "All 11 tunnels", "2GB+ RAM"},
			{"lite", "Lightweight", "512MB-2GB"},
			{"minimal", "SSH + SS only", "<512MB"},
		},
	}
}

func (m *connectWizard) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.Tick(timeConstant, func(t time.Time) tea.Msg {
		return tickMsg(t)
	}))
}

func (m *connectWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return nil, tea.Quit
		case "q":
			if m.step == 0 {
				m.quitting = true
				return nil, tea.Quit
			}
		case "esc":
			if m.step > 0 {
				m.step--
				m.focusCurrent()
			} else {
				m.quitting = true
				return nil, tea.Quit
			}
			return nil, nil
		case "enter":
			return m.advanceStep()
		case "tab":
			if m.step < 4 {
				m.step++
				m.focusCurrent()
			}
			return nil, nil
		case "up", "k":
			if m.step == 3 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.modeOptions) - 1
				}
			}
		case "down", "j":
			if m.step == 3 {
				m.cursor++
				if m.cursor >= len(m.modeOptions) {
					m.cursor = 0
				}
			}
		}

	case tickMsg:
		m.tick++
		return nil, tea.Tick(timeConstant, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch m.step {
	case 0:
		m.addrInput, cmd = m.addrInput.Update(msg)
		cmds = append(cmds, cmd)
	case 1:
		m.userInput, cmd = m.userInput.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		m.passInput, cmd = m.passInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return nil, tea.Batch(cmds...)
}

func (m *connectWizard) focusCurrent() {
	m.addrInput.Blur()
	m.userInput.Blur()
	m.passInput.Blur()
	switch m.step {
	case 0:
		m.addrInput.Focus()
	case 1:
		m.userInput.Focus()
	case 2:
		m.passInput.Focus()
	}
}

func (m *connectWizard) advanceStep() (tea.Model, tea.Cmd) {
	switch m.step {
	case 0:
		if m.addrInput.Value() == "" {
			return nil, nil
		}
		m.step = 1
		m.focusCurrent()
	case 1:
		if m.userInput.Value() == "" {
			m.userInput.SetValue("root")
		}
		m.step = 2
		m.focusCurrent()
	case 2:
		m.step = 3
	case 3:
		m.mode = m.modeOptions[m.cursor].name
		m.step = 4
	case 4:
		return nil, tea.Quit
	}
	return nil, nil
}

func (m *connectWizard) View() string {
	if m.quitting {
		return "\n  Cancelled.\n\n"
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(renderGradient("  NYXORA Connect", currentTheme.GradientA, currentTheme.GradientB))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 54)))
	b.WriteString("\n\n")

	steps := []struct {
		label string
		done  bool
		curr  bool
	}{
		{"Address", m.step > 0, m.step == 0},
		{"User", m.step > 1, m.step == 1},
		{"Password", m.step > 2, m.step == 2},
		{"Mode", m.step > 3, m.step == 3},
		{"Go!", m.step > 4, m.step == 4},
	}

	for _, s := range steps {
		var icon string
		var style lipgloss.Style
		if s.done {
			icon = statusIcons["active"]
			style = successStyle()
		} else if s.curr {
			icon = statusIcons["testing"]
			style = accentStyle()
		} else {
			icon = statusIcons["idle"]
			style = dimStyle()
		}
		b.WriteString(fmt.Sprintf("  %s %s\n", style.Render(icon), style.Render(s.label)))
	}

	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 54)))
	b.WriteString("\n\n")

	switch m.step {
	case 0:
		b.WriteString(accentStyle().Render("  Remote server address:"))
		b.WriteString("\n")
		b.WriteString(m.addrInput.View())
		b.WriteString("\n")
	case 1:
		b.WriteString(accentStyle().Render("  SSH username:"))
		b.WriteString("\n")
		b.WriteString(m.userInput.View())
		b.WriteString("\n")
	case 2:
		b.WriteString(accentStyle().Render("  SSH password:"))
		b.WriteString("\n")
		b.WriteString(m.passInput.View())
		b.WriteString("\n")
	case 3:
		b.WriteString(accentStyle().Render("  Server mode:"))
		b.WriteString("\n\n")
		for i, md := range m.modeOptions {
			cur := "  "
			style := dimStyle()
			if m.cursor == i {
				cur = successStyle().Render("▸ ")
				style = lipgloss.NewStyle().
					Foreground(lipgloss.Color(currentTheme.Highlight)).
					Bold(true)
			}
			icon := transportIcons[md.name]
			if icon == "" {
				icon = "●"
			}
			b.WriteString(fmt.Sprintf("  %s%s %s %s\n", cur, icon, style.Render(md.name), dimStyle().Render(fmt.Sprintf("(%s) [%s]", md.desc, md.req))))
		}
	case 4:
		b.WriteString(successStyle().Render("  Ready to connect!"))
		b.WriteString("\n\n")
		summaryStyle := textStyle()
		b.WriteString(fmt.Sprintf("  %s %s\n", dimStyle().Render("Address:"), summaryStyle.Render(m.addrInput.Value())))
		b.WriteString(fmt.Sprintf("  %s %s\n", dimStyle().Render("User:"), summaryStyle.Render(m.userInput.Value())))
		b.WriteString(fmt.Sprintf("  %s %s\n", dimStyle().Render("Password:"), dimStyle().Render(strings.Repeat("•", len(m.passInput.Value())))))
		b.WriteString(fmt.Sprintf("  %s %s\n", dimStyle().Render("Mode:"), summaryStyle.Render(m.mode)))
		b.WriteString(fmt.Sprintf("\n  %s", successStyle().Render("Press Enter to connect!")))
	}

	b.WriteString("\n\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 54)))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  esc back  \u23CE next  \u21B9 skip  q cancel"))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render("  \u200Ehttps://t.me/NyxoraCore"))
	b.WriteString("\n")

	return b.String()
}
