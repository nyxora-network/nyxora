package interactive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	currentVersion = "0.2.0"
	updateURL      = "https://api.github.com/repos/nyxorammd-lgtm/nyxora/releases/latest"
	downloadBase   = "https://github.com/nyxorammd-lgtm/nyxora/releases/download"
	telegramURL    = "https://t.me/NyxoraCore"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type updateModel struct {
	state       string // checking, found, notfound, downloading, done, error
	latestVer   string
	downloadURL string
	progress    int
	err         string
	width       int
	tick        int
	quitting    bool
}

func RunUpdateChecker() error {
	p := tea.NewProgram(updateModel{state: "checking"}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m updateModel) Init() tea.Cmd {
	return tea.Batch(
		checkForUpdate(),
		tickUpdateCmd(),
	)
}

func checkForUpdate() tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(updateURL)
		if err != nil {
			return updateErrMsg{err: err.Error()}
		}
		defer resp.Body.Close()

		var release githubRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return updateErrMsg{err: err.Error()}
		}

		return updateFoundMsg{
			version: release.TagName,
			assets:  release.Assets,
		}
	}
}

type updateFoundMsg struct {
	version string
	assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	}
}

type updateErrMsg struct {
	err string
}

type updateTickMsg time.Time

func tickUpdateCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return updateTickMsg(t)
	})
}

func (m updateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.state == "found" && m.downloadURL != "" {
				m.state = "downloading"
				return m, downloadUpdate(m.downloadURL)
			}
			if m.state == "done" || m.state == "notfound" || m.state == "error" {
				return m, tea.Quit
			}
		}

	case updateFoundMsg:
		// Strip leading 'v' if present
		ver := msg.version
		ver = strings.TrimPrefix(ver, "v")

		if ver == currentVersion {
			m.state = "notfound"
		} else {
			m.state = "found"
			m.latestVer = ver
			// Find binary for current OS/arch
			arch := runtime.GOARCH
			osName := runtime.GOOS
			for _, asset := range msg.assets {
				name := strings.ToLower(asset.Name)
				if strings.Contains(name, osName) && strings.Contains(name, arch) {
					m.downloadURL = asset.BrowserDownloadURL
					break
				}
			}
			if m.downloadURL == "" {
				m.downloadURL = fmt.Sprintf("%s/v%s/nyxora_%s_%s", downloadBase, ver, osName, arch)
			}
		}

	case updateErrMsg:
		m.state = "error"
		m.err = msg.err

	case updateDownloadedMsg:
		m.state = "done"

	case updateDownloadErrMsg:
		m.state = "error"
		m.err = msg.err

	case updateTickMsg:
		m.tick++
		if m.state == "downloading" {
			m.progress = (m.progress + 2) % 100
		}
		return m, tickUpdateCmd()

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

type updateDownloadedMsg struct{}
type updateDownloadErrMsg struct{ err string }

func downloadUpdate(url string) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(url)
		if err != nil {
			return updateDownloadErrMsg{err: err.Error()}
		}
		defer resp.Body.Close()

		tmpFile, err := os.CreateTemp("", "nyxora-update-*")
		if err != nil {
			return updateDownloadErrMsg{err: err.Error()}
		}
		defer tmpFile.Close()

		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			return updateDownloadErrMsg{err: err.Error()}
		}

		// Make executable and replace current binary
		tmpFile.Chmod(0755)
		tmpFile.Close()

		execPath, err := os.Executable()
		if err != nil {
			return updateDownloadErrMsg{err: err.Error()}
		}

		if err := os.Rename(tmpFile.Name(), execPath); err != nil {
			return updateDownloadErrMsg{err: err.Error()}
		}

		return updateDownloadedMsg{}
	}
}

func (m updateModel) View() string {
	if m.quitting {
		return "\n"
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("141")).Bold(true).Render("  NYXORA Update Checker"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(strings.Repeat("─", 50)))
	b.WriteString("\n\n")

	switch m.state {
	case "checking":
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		idx := m.tick % len(spinner)
		b.WriteString(fmt.Sprintf("  %s Checking for updates...\n", lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(spinner[idx])))

	case "found":
		b.WriteString(fmt.Sprintf("  Current version:  %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(currentVersion)))
		b.WriteString(fmt.Sprintf("  Latest version:   %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true).Render(m.latestVer)))
		b.WriteString(fmt.Sprintf("  Status:           %s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("Update available!")))
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render("  Press Enter to download & install"))
		b.WriteString("\n")

	case "notfound":
		b.WriteString(fmt.Sprintf("  Current version:  %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(currentVersion)))
		b.WriteString(fmt.Sprintf("  Status:           %s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("You're up to date!")))
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render("  Press any key to exit"))
		b.WriteString("\n")

	case "downloading":
		b.WriteString(fmt.Sprintf("  Downloading v%s...\n\n", m.latestVer))
		barWidth := 40
		progress := m.progress * barWidth / 100
		bar := strings.Repeat("█", progress) + strings.Repeat("░", barWidth-progress)
		b.WriteString(fmt.Sprintf("  %s %d%%\n", lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(bar), m.progress))

	case "done":
		b.WriteString(fmt.Sprintf("  %s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true).Render("Update installed successfully!")))
		b.WriteString(fmt.Sprintf("  Restart NYXORA to use v%s\n", m.latestVer))
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render("  Press any key to exit"))
		b.WriteString("\n")

	case "error":
		b.WriteString(fmt.Sprintf("  %s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Update failed:")))
		b.WriteString(fmt.Sprintf("  %s\n\n", m.err))
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render("  Press any key to exit"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(strings.Repeat("─", 50)))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(fmt.Sprintf("  📱 Telegram: %s", telegramURL)))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true).Render("  enter select  •  q/esc back"))
	b.WriteString("\n")

	return b.String()
}
