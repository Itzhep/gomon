package watcher

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fsWatcher  *fsnotify.Watcher
	cmd        *exec.Cmd
	appPath    string
	dirPath    string
	outputDir  string // Directory to store the built binary
	debounce   time.Duration
	lastChange time.Time
	isRunning   bool // Flag to check if the app is currently running
}

func isRunningInDocker() bool {
	data, err := ioutil.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "docker")
}

func NewWatcher(appPath string, debounce time.Duration) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Dir(appPath)
	outputDir := filepath.Join(dirPath, "output") // Specify the output directory

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, err
	}

	return &Watcher{
		fsWatcher: fsWatcher,
		appPath:   appPath,
		dirPath:   dirPath,
		outputDir: outputDir,
		debounce:  debounce,
		isRunning: false, // Initialize the running flag
	}, nil
}

func (w *Watcher) BuildBinary() error {
	binaryPath := filepath.Join(w.outputDir, "app.exe") // Specify the binary path

	// Build the application
	cmd := exec.Command("go", "build", "-o", binaryPath, w.appPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		color.Red("Build error: %v", err)
		return err // Return the error for further handling
	}

	color.Green("Built binary: %s", binaryPath)
	return nil // Return nil if the build was successful
}

func (w *Watcher) StartApp() error {
	if w.cmd != nil && w.cmd.Process != nil {
		_ = w.cmd.Process.Kill() // Kill previous instance if running
	}
	color.Cyan("Starting app...")

	// Run the .exe file in the output directory
	w.cmd = exec.Command(filepath.Join(w.outputDir, "app.exe")) // Run the compiled application
	w.cmd.Stdout = os.Stdout
	w.cmd.Stderr = os.Stderr

	if err := w.cmd.Start(); err != nil {
		return err // Return error if starting fails
	}
	w.isRunning = true // Set running flag to true
	return nil // Return nil if the start was successful
}

func (w *Watcher) WatchAndReload() error {
	defer w.fsWatcher.Close()

	done := make(chan bool)

	ignoredDirs := []string{".git", "vendor", "node_modules", "output"} // Directories to ignore

	err := filepath.Walk(w.dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if running in Docker
		if isRunningInDocker() {
			color.Cyan("Detected Docker environment, running docker-compose up...")
			cmd := exec.Command("docker-compose", "up", "-d") // Restart via docker-compose
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run() // Run docker-compose up
		}

		// Skip ignored directories
		for _, dir := range ignoredDirs {
			if info.IsDir() && info.Name() == dir {
				return filepath.SkipDir
			}
		}

		if info.IsDir() {
			return w.fsWatcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := w.BuildBinary(); err != nil {
		return err
	}

	if err := w.StartApp(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-w.fsWatcher.Events:
				if !ok {
					return
				}

				// Only restart the app when a .go file in the active set is changed
				if time.Since(w.lastChange) < w.debounce {
					continue
				}

				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 &&
					filepath.Ext(event.Name) == ".go" {
					w.lastChange = time.Now()
					color.Yellow("File changed: %s", event.Name)

					// Check if the app is already running
					if w.isRunning {
						color.Cyan("Restarting app...")
						if err := w.BuildBinary(); err != nil {
							color.Red("Failed to build the app: %v", err)
							continue
						}

						if err := w.StartApp(); err != nil {
							color.Red("Failed to restart the app: %v", err)
						}
					}
				}

			case err, ok := <-w.fsWatcher.Errors:
				if !ok {
					return
				}
				color.Red("Error: %v", err)
			}
		}
	}()

	go w.handleUserInput()

	<-done
	return nil
}

func (w *Watcher) Stop() {
	if w.cmd != nil && w.cmd.Process != nil {
		_ = w.cmd.Process.Kill() // Kill the app process
	}
	w.isRunning = false // Reset running flag
	_ = w.fsWatcher.Close() // Close the file watcher
}

func (w *Watcher) handleUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		if command == "rs" {
			color.Cyan("Restarting app...")
			if err := w.BuildBinary(); err != nil {
				color.Red("Failed to build the app: %v", err)
				continue
			}
			if err := w.StartApp(); err != nil {
				color.Red("Failed to restart the app: %v", err)
			}
		} else {
			color.Red("Unknown command: %s", command)
		}
	}
	if err := scanner.Err(); err != nil {
		color.Red("Error reading from stdin: %v", err)
	}
}
