package interactive

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func primaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Primary)).Bold(true)
}

func accentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Accent)).Bold(true)
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

func mutedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Muted))
}

func textStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Text))
}

func badgeStyle(label string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.Bg)).
		Background(lipgloss.Color(currentTheme.Primary)).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)
}

func boxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(currentTheme.Border)).
		Padding(1, 2).
		Width(58)
}

func subtleBox() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(currentTheme.Muted)).
		Padding(0, 1)
}

func renderGradient(text, colorA, colorB string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return ""
	}

	var result strings.Builder
	for i, r := range runes {
		t := float64(i) / float64(len(runes)-1)
		c := lerpHex(colorA, colorB, t)
		result.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color(c)).
			Bold(true).
			Render(string(r)))
	}
	return result.String()
}

func renderGradientBold(text, colorA, colorB string) string {
	return renderGradient(text, colorA, colorB)
}

func renderProgressBar(filled, width int) string {
	if filled > width {
		filled = width
	}
	fillChar := "█"
	emptyChar := "░"

	bar := strings.Repeat(fillChar, filled) + strings.Repeat(emptyChar, width-filled)

	low := int(float64(filled) / float64(width) * 100)
	style := errorStyle()
	if low >= 70 {
		style = successStyle()
	} else if low >= 40 {
		style = warningStyle()
	}

	return style.Render(bar)
}

func renderGradientBar(score float64, width int, colorA, colorB string) string {
	filled := int((score / 100) * float64(width))
	if filled > width {
		filled = width
	}

	bar := ""
	for i := 0; i < width; i++ {
		t := float64(i) / float64(width-1)
		if i < filled {
			c := lerpHex(colorA, colorB, t)
			bar += lipgloss.NewStyle().
				Foreground(lipgloss.Color(c)).
				Render("█")
		} else {
			bar += mutedStyle().Render("░")
		}
	}
	return bar
}

func lerpHex(a, b string, t float64) string {
	if len(a) == 7 && len(b) == 7 {
		r1, g1, b1 := parseHex(a)
		r2, g2, b2 := parseHex(b)
		r := int(float64(r1)*(1-t) + float64(r2)*t)
		g := int(float64(g1)*(1-t) + float64(g2)*t)
		bv := int(float64(b1)*(1-t) + float64(b2)*t)
		return fmt.Sprintf("#%02X%02X%02X", clamp(r), clamp(g), clamp(bv))
	}
	if t < 0.5 {
		return a
	}
	return b
}

func parseHex(hex string) (int, int, int) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0
	}
	r := hexToInt(hex[1:3])
	g := hexToInt(hex[3:5])
	b := hexToInt(hex[5:7])
	return r, g, b
}

func hexToInt(s string) int {
	val := 0
	for _, c := range s {
		val *= 16
		switch {
		case c >= '0' && c <= '9':
			val += int(c - '0')
		case c >= 'A' && c <= 'F':
			val += int(c-'A') + 10
		case c >= 'a' && c <= 'f':
			val += int(c-'a') + 10
		}
	}
	return val
}

func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}
