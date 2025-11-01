package cmd

import (
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/spf13/cobra"
	"github.com/zautumnz/tg/internal/app/tg/config"
	"github.com/zautumnz/tg/internal/app/tg/gui"
	"github.com/zautumnz/tg/internal/app/tg/termutil"
	"github.com/zautumnz/tg/internal/app/tg/version"
)

var rewriteConfig bool
var debugFile string
var initialCommand string
var shell string
var themePath string
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

		if rewriteConfig {
			if _, err := conf.Save(); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}
			fmt.Println("Config written.")
			return nil
		}

		var theme *termutil.Theme

		if themePath != "" {
			theme, err = config.LoadThemeFromPath(conf, themePath)
			if err != nil {
				return fmt.Errorf("failed to load theme: %s", err)
			}
		} else {
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
		}

		termOpts := []termutil.Option{
			termutil.WithTheme(theme),
		}

		if debugFile != "" {
			termOpts = append(termOpts, termutil.WithLogFile(debugFile))
		}
		if shell != "" {
			termOpts = append(termOpts, termutil.WithShell(shell))
		}
		if initialCommand != "" {
			termOpts = append(termOpts, termutil.WithInitialCommand(initialCommand))
		}

		terminal := termutil.New(termOpts...)

		options := []gui.Option{
			gui.WithFontDPI(conf.Font.DPI),
			gui.WithFontSize(conf.Font.Size),
			gui.WithFontFamily(conf.Font.Family),
			gui.WithOpacity(conf.Opacity),
			gui.WithLigatures(conf.Font.Ligatures),
		}

		if conf.Cursor.Image != "" {
			img, err := getImageFromFilePath(conf.Cursor.Image)
			if err != nil {
				startupErrors = append(startupErrors, err)
			} else {
				options = append(options, gui.WithCursorImage(img))
			}
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

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func Execute() error {
	rootCmd.Flags().BoolVar(&showVersion, "version", showVersion, "Show term version information and exit")
	rootCmd.Flags().BoolVar(&rewriteConfig, "rewrite-config", rewriteConfig, "Write the resultant config after parsing config files and merging with defauls back to the config file")
	rootCmd.Flags().StringVar(&debugFile, "log-file", debugFile, "Debug log file")
	rootCmd.Flags().StringVarP(&shell, "shell", "s", shell, "Shell to launch terminal with - defaults to configured user shell")
	rootCmd.Flags().StringVarP(&initialCommand, "command", "c", initialCommand, "Command to run when shell starts - use this with caution")
	rootCmd.Flags().StringVar(&themePath, "theme-path", themePath, "Path to a theme file to use instead of the default")
	return rootCmd.Execute()
}
