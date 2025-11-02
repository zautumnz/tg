package termutil

import (
	"os"
)

type Option func(t *Terminal)

func WithLogFile(path string) Option {
	return func(t *Terminal) {
		if path == "-" {
			t.logFile = os.Stdout
			return
		}
		t.logFile, _ = os.Create(path)
	}
}

func WithTheme(theme *Theme) Option {
	return func(t *Terminal) {
		t.theme = theme
	}
}

func WithWindowManipulator(m WindowManipulator) Option {
	return func(t *Terminal) {
		t.windowManipulator = m
	}
}
