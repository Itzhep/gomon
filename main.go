package main

import (
	"fmt"
	"os"
	"time"
	"net/http"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"gomon/watcher"
)

var (
	appPath     string
	debounce    time.Duration
	w           *watcher.Watcher
	useDocker   bool
	verbose     bool
	excludeDirs []string
)

var upgrader = websocket.Upgrader{}

func liveReloadHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        color.Red("❌ Failed to upgrade connection: %v", err)
        return
    }
    defer conn.Close()

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break
        }
    }
}

func startLiveReloadServer() {
    http.HandleFunc("/livereload", liveReloadHandler)
    go http.ListenAndServe(":35729", nil)
    color.HiGreen("🔄 Live reload server started on port 35729")
}

func init() {
	// Root command flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	// Start command flags
	startCmd.Flags().StringVarP(&appPath, "app", "a", "", "Path to the Go application to run")
	startCmd.Flags().DurationVarP(&debounce, "debounce", "d", 1*time.Second, "Debounce duration for file changes")
	startCmd.Flags().BoolVarP(&useDocker, "docker", "", false, "Use Docker for restarting the app")
	startCmd.Flags().StringSliceVarP(&excludeDirs, "exclude", "e", []string{".git", "vendor", "node_modules"}, "Directories to exclude from watching")
	
	// Mark required flags
	startCmd.MarkFlagRequired("app")
}

var rootCmd = &cobra.Command{
	Use:   "gomon",
	Short: "Gomon is a tool to automatically restart your Go application on file changes",
	Long: color.GreenString(`
🛠️  Gomon - Go File Watcher & Auto-Reloader
----------------------------------------
Automatically rebuilds and restarts your Go application when file changes are detected.
For more information, visit: https://github.com/Itzhep/gomon
	`),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		printBanner()
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start watching and auto-reloading your application",
	Example: `  gomon start --app main.go
  gomon start -a main.go -d 2s
  gomon start -a main.go --docker
  gomon start -a main.go -e "tmp,logs,test"`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		w, err = watcher.NewWatcher(appPath, debounce)
		if err != nil {
			color.Red("❌ Failed to initialize watcher: %v", err)
			os.Exit(1)
		}

		color.HiCyan("👀 Watching for file changes...")
		color.HiYellow("💡 Type 'rs' + Enter to manually restart")
		color.HiYellow("💡 Press Ctrl+C to exit")

		if useDocker {
			color.HiBlue("🐳 Docker mode enabled")
		}

		if verbose {
			color.HiWhite("📝 Verbose logging enabled")
			color.HiWhite("📂 Watching directory: %s", appPath)
			color.HiWhite("⏱️  Debounce time: %v", debounce)
			color.HiWhite("🚫 Excluded directories: %v", excludeDirs)
		}

		startLiveReloadServer()

		if err := w.WatchAndReload(); err != nil {
			color.Red("❌ Error: %v", err)
			os.Exit(1)
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the watcher",
	Run: func(cmd *cobra.Command, args []string) {
		if w != nil {
			//w.Stop() // call the stop method
			color.Green("Watcher stopped.")
		} else {
			color.Red("Watcher is not running.")
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		color.HiCyan("Gomon v2.0.0")
	},
}

func printBanner() {
	banner := color.HiCyanString(`
   ______                          
  / ____/___  ____ ___  ____  ____ 
 / / __/ __ \/ __ '__ \/ __ \/ __ \
/ /_/ / /_/ / / / / / / /_/ / / / /
\____/\____/_/ /_/ /_/\____/_/ /_/ 
`)
	fmt.Println(banner)
}

func main() {
    startLiveReloadServer()
    rootCmd.AddCommand(startCmd, stopCmd, versionCmd)
    if err := rootCmd.Execute(); err != nil {
        color.Red("❌ Error: %v", err)
        os.Exit(1)
    }
}


