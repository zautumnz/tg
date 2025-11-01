package main

import (
	"os"

	"github.com/zautumnz/tg/internal/app/tg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
