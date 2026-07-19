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

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	currentVersion = "0.2.0"
	updateURL      = "https://api.github.com/repos/nyxora-network/nyxora/releases/latest"
	downloadBase   = "https://github.com/nyxora-network/nyxora/releases/download"
	telegramURL    = "https://t.me/NyxoraCore"
)

const timeConstant = time.Millisecond * 80

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type updateModel struct {
	state       string
	latestVer   string
	downloadURL string
	progress    int
	err         string
	width       int
	tick        int
	quitting    bool
	spinner     spinner.Model
}

// RunUpdateChecker launches a Bubble Tea TUI that checks GitHub for the latest release,
// displays version info, and optionally downloads and installs updates.
func RunUpdateChecker() error {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Primary))
	s.Spinner = spinner.Pulse

	p := tea.NewProgram(updateModel{
		state:   "checking",
		spinner: s,
	}, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m updateModel) Init() tea.Cmd {
	return tea.Batch(
		checkForUpdate(),
		tea.Tick(timeConstant, func(t time.Time) tea.Msg {
			return updateTickMsg(t)
		}),
		m.spinner.Tick,
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

		if resp.StatusCode != http.StatusOK {
			return updateErrMsg{err: fmt.Sprintf("GitHub API returned status %d", resp.StatusCode)}
		}

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

	case updateFoundMsg:
		ver := msg.version
		ver = strings.TrimPrefix(ver, "v")

		if ver == currentVersion {
			m.state = "notfound"
		} else {
			m.state = "found"
			m.latestVer = ver
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
		return m, tea.Tick(timeConstant, func(t time.Time) tea.Msg {
			return updateTickMsg(t)
		})

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
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

		if resp.StatusCode != http.StatusOK {
			return updateDownloadErrMsg{err: fmt.Sprintf("download failed with status %d", resp.StatusCode)}
		}

		tmpFile, err := os.CreateTemp("", "nyxora-update-*")
		if err != nil {
			return updateDownloadErrMsg{err: err.Error()}
		}
		defer tmpFile.Close()

		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			return updateDownloadErrMsg{err: err.Error()}
		}

		tmpFile.Chmod(0755)

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
	b.WriteString(renderGradient("  NYXORA Update", currentTheme.GradientA, currentTheme.GradientB))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 50)))
	b.WriteString("\n\n")

	switch m.state {
	case "checking":
		b.WriteString(fmt.Sprintf("  %s %s\n",
			primaryStyle().Render(m.spinner.View()),
			infoStyle().Render("Checking for updates..."),
		))

	case "found":
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			dimStyle().Render("Current version:"),
			dimStyle().Render(currentVersion),
		))
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			dimStyle().Render("Latest version:"),
			successStyle().Render(m.latestVer),
		))
		b.WriteString(fmt.Sprintf("  %s  %s\n\n",
			dimStyle().Render("Status:"),
			warningStyle().Render("Update available!"),
		))
		b.WriteString(fmt.Sprintf("  %s\n", accentStyle().Render("Press Enter to download & install")))

	case "notfound":
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			dimStyle().Render("Current version:"),
			dimStyle().Render(currentVersion),
		))
		b.WriteString(fmt.Sprintf("  %s  %s\n\n",
			dimStyle().Render("Status:"),
			successStyle().Render("You're up to date!"),
		))
		b.WriteString(fmt.Sprintf("  %s\n", accentStyle().Render("Press any key to exit")))

	case "downloading":
		b.WriteString(fmt.Sprintf("  %s v%s...\n\n",
			infoStyle().Render("Downloading"),
			m.latestVer,
		))
			barWidth := 40
		bar := renderGradientBar(float64(m.progress), barWidth, currentTheme.GradientA, currentTheme.GradientB)
		b.WriteString(fmt.Sprintf("  %s %d%%\n", bar, m.progress))

	case "done":
		b.WriteString(fmt.Sprintf("  %s\n\n",
			successStyle().Render("Update installed successfully!"),
		))
		b.WriteString(fmt.Sprintf("  %s v%s\n",
			dimStyle().Render("Restart NYXORA to use"),
			m.latestVer,
		))
		b.WriteString(fmt.Sprintf("  %s\n", accentStyle().Render("Press any key to exit")))

	case "error":
		b.WriteString(fmt.Sprintf("  %s\n\n",
			errorStyle().Render("Update failed:"),
		))
		b.WriteString(fmt.Sprintf("  %s\n\n", dimStyle().Render(m.err)))
		b.WriteString(fmt.Sprintf("  %s\n", accentStyle().Render("Press any key to exit")))
	}

	b.WriteString("\n")
	b.WriteString(dimStyle().Render(strings.Repeat("─", 50)))
	b.WriteString("\n")
	b.WriteString(dimStyle().Render(fmt.Sprintf("  \u200ETelegram: %s", telegramURL)))
	b.WriteString("\n")
	b.WriteString(dimStyle().Italic(true).Render("  enter select  \u2022  q/esc back"))
	b.WriteString("\n")

	return b.String()
}
