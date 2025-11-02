package config

type Theme struct {
	Black               string
	Red                 string
	Green               string
	Yellow              string
	Blue                string
	Magenta             string
	Cyan                string
	White               string
	BrightBlack         string
	BrightRed           string
	BrightGreen         string
	BrightYellow        string
	BrightBlue          string
	BrightMagenta       string
	BrightCyan          string
	BrightWhite         string
	Background          string
	Foreground          string
	SelectionBackground string
	SelectionForeground string
	CursorForeground    string
	CursorBackground    string
}

func loadTheme() (*Theme, error) {
	theme := defaultTheme
	return &theme, nil
}
