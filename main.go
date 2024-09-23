package main

import (
	"os"
	"time"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"gomon/watcher"
)

var (
	appPath   string
	debounce  time.Duration
)

func init() {
	rootCmd.Flags().StringVarP(&appPath, "app", "a", "", "Path to the Go application to run")
	rootCmd.Flags().DurationVarP(&debounce, "debounce", "d", 1*time.Second, "Debounce duration for file changes")
	rootCmd.MarkFlagRequired("app")
}

var rootCmd = &cobra.Command{
	Use:   "gomon",
	Short: "Gomon is a tool to automatically restart your Go application on file changes",
	Run: func(cmd *cobra.Command, args []string) {
		if appPath == "" {
			color.Red("Error: Path to the application must be specified with --app")
			os.Exit(1)
		}

		w, err := watcher.NewWatcher(appPath, debounce)
		if err != nil {
			color.Red("Failed to initialize watcher: %v", err)
			os.Exit(1)
		}

		color.Green("Watching for file changes... Press 'rs' And Than Enter to restart manually.")
		color.Cyan("If you find Gomon useful, please consider starring it on GitHub: https://github.com/Itzhep/gomon")
		if err := w.WatchAndReload(); err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
