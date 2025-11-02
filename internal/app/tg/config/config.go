package config

import (
	"fmt"
)

type Config struct {
	Font   Font
	Cursor Cursor
}

type Font struct {
	Family    string
	Size      float64
	DPI       float64
	Ligatures bool
}

type Cursor struct {
	Image string
}

type ErrorFileNotFound struct {
	Path string
}

func (e *ErrorFileNotFound) Error() string {
	return fmt.Sprintf("file was not found at '%s'", e.Path)
}

func LoadConfig() (*Config, error) {
	config := defaultConfig
	return &config, nil
}
