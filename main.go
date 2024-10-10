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
	w         *watcher.Watcher
	useDocker bool
)

func init() {
	// Flags for both the root and start commands
	startCmd.Flags().StringVarP(&appPath, "app", "a", "", "Path to the Go application to run")
	startCmd.Flags().DurationVarP(&debounce, "debounce", "d", 1*time.Second, "Debounce duration for file changes")
	startCmd.Flags().BoolVarP(&useDocker, "docker", "", false, "Use Docker for restarting the app")

}

var rootCmd = &cobra.Command{
	Use:   "gomon",
	Short: "Gomon is a tool to automatically restart your Go application on file changes",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the watcher",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the appPath is set correctly
		if appPath == "" {
			color.Red("Error: Path to the application must be specified with --app")
			os.Exit(1)
		}

		var err error
		w, err = watcher.NewWatcher(appPath, debounce)
		if err != nil {
			color.Red("Failed to initialize watcher: %v", err)
			os.Exit(1)
		}

		color.Green("Watching for file changes... Press 'rs' and Enter to restart manually.")
		if err := w.WatchAndReload(); err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the watcher",
	Run: func(cmd *cobra.Command, args []string) {
		if w != nil {
			w.Stop() // call the stop method
			color.Green("Watcher stopped.")
		} else {
			color.Red("Watcher is not running.")
		}
	},
}

func main() {
	// Register the start and stop commands to root
	rootCmd.AddCommand(startCmd, stopCmd)
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
