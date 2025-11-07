package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zautumnz/tg/internal/config"
	"github.com/zautumnz/tg/internal/gui"
	"github.com/zautumnz/tg/internal/termutil"
	"github.com/zautumnz/tg/internal/version"
)

var debugFile string
var showVersion bool

var rootCmd = &cobra.Command{
	Use:          os.Args[0],
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {

		if showVersion {
			fmt.Println(version.Version)
			os.Exit(0)
		}

		var startupErrors []error
		var fileNotFound *config.ErrorFileNotFound

		conf, err := config.LoadConfig()
		if err != nil {
			if !errors.As(err, &fileNotFound) {
				startupErrors = append(startupErrors, err)
			}
			conf = config.DefaultConfig()
		}

		var theme *termutil.Theme

		theme, err = config.LoadTheme(conf)
		if err != nil {
			if !errors.As(err, &fileNotFound) {
				startupErrors = append(startupErrors, err)
			}
			theme, err = config.DefaultTheme(conf)
			if err != nil {
				return fmt.Errorf("failed to load default theme: %w", err)
			}
		}

		termOpts := []termutil.Option{
			termutil.WithTheme(theme),
		}

		if debugFile != "" {
			termOpts = append(termOpts, termutil.WithLogFile(debugFile))
		}

		terminal := termutil.New(termOpts...)

		options := []gui.Option{
			gui.WithFontDPI(conf.Font.DPI),
			gui.WithFontSize(conf.Font.Size),
			gui.WithFontFamily(conf.Font.Family),
		}

		g, err := gui.New(terminal, options...)
		if err != nil {
			return err
		}

		for _, err := range startupErrors {
			g.ShowError(err.Error())
		}

		return g.Run()
	},
}

func Execute() error {
	rootCmd.Flags().BoolVar(&showVersion, "version", showVersion, "Show term version information and exit")
	rootCmd.Flags().StringVar(&debugFile, "log-file", debugFile, "Debug log file")
	return rootCmd.Execute()
}
