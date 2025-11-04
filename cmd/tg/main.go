package main

import (
	"os"

	"github.com/zautumnz/tg/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
