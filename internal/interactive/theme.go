package interactive

// Theme defines a complete color palette for the interactive TUI.
// All colors are specified as hex strings (e.g. "#CBA6F7") for TrueColor support.
type Theme struct {
	Primary   string
	Secondary string
	Accent    string
	Bg        string
	Surface   string
	Success   string
	Warning   string
	Error     string
	Info      string
	Border    string
	Highlight string
	Text      string
	TextDim   string
	Muted     string
	GradientA string
	GradientB string
}

type ThemeSet struct {
	Dark  Theme
	Light Theme
	Neon  Theme
}

var themes = map[string]Theme{
	"catppuccin-mocha": {
		Primary:   "#CBA6F7",
		Secondary: "#89B4FA",
		Accent:    "#F5C2E7",
		Bg:        "#1E1E2E",
		Surface:   "#313244",
		Success:   "#A6E3A1",
		Warning:   "#F9E2AF",
		Error:     "#F38BA8",
		Info:      "#89B4FA",
		Border:    "#585B70",
		Highlight: "#F5C2E7",
		Text:      "#CDD6F4",
		TextDim:   "#6C7086",
		Muted:     "#45475A",
		GradientA: "#CBA6F7",
		GradientB: "#89B4FA",
	},
	"tokyo-night": {
		Primary:   "#7AA2F7",
		Secondary: "#BB9AF7",
		Accent:    "#FF9E64",
		Bg:        "#1A1B26",
		Surface:   "#24283B",
		Success:   "#9ECE6A",
		Warning:   "#E0AF68",
		Error:     "#F7768E",
		Info:      "#7DCFFF",
		Border:    "#3B4261",
		Highlight: "#FF9E64",
		Text:      "#A9B1D6",
		TextDim:   "#565F89",
		Muted:     "#2F3346",
		GradientA: "#7AA2F7",
		GradientB: "#BB9AF7",
	},
	"catppuccin-latte": {
		Primary:   "#8839EF",
		Secondary: "#1E66F5",
		Accent:    "#EA76CB",
		Bg:        "#EFF1F5",
		Surface:   "#E6E9EF",
		Success:   "#40A02B",
		Warning:   "#DF8E1D",
		Error:     "#D20F39",
		Info:      "#1E66F5",
		Border:    "#9CA0B0",
		Highlight: "#EA76CB",
		Text:      "#4C4F69",
		TextDim:   "#9CA0B0",
		Muted:     "#CCD0DA",
		GradientA: "#8839EF",
		GradientB: "#1E66F5",
	},
}

var currentTheme = themes["catppuccin-mocha"]
