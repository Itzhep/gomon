<<<<<<< HEAD
package watcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type BuildConfig struct {
	Flags    []string
	Env      []string
	Commands []string
	Scripts  []string // Add custom scripts
}

type WatchConfig struct {
	Extensions  []string
	ExcludeDirs []string
	IncludeDirs []string
}

type Watcher struct {
	fsWatcher         *fsnotify.Watcher
	cmd               *exec.Cmd
	appPath           string
	dirPath           string
	outputDir         string
	debounce          time.Duration
	lastChange        time.Time
	isRunning         bool
	mutex             sync.Mutex
	buildConfig       BuildConfig
	watchConfig       WatchConfig
	startTime         time.Time
	buildCount        int
	lastBuildDuration time.Duration
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
		return nil, fmt.Errorf("failed to create watcher: %v", err)
	}

	dirPath := filepath.Dir(appPath)
	outputDir := filepath.Join(dirPath, "bin")

	return &Watcher{
		fsWatcher: fsWatcher,
		appPath:   appPath,
		dirPath:   dirPath,
		outputDir: outputDir,
		debounce:  debounce,
		buildConfig: BuildConfig{
			Flags: []string{"-race"},
			Env:   []string{"CGO_ENABLED=1"},
		},
		watchConfig: WatchConfig{
			Extensions:  []string{".go", ".mod", ".sum"},
			ExcludeDirs: []string{"vendor", "node_modules", ".git"},
			IncludeDirs: []string{"."},
		},
		startTime: time.Now(),
	}, nil
}

func (w *Watcher) WatchAndReload() error {
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	if err := w.addWatchDirs(); err != nil {
		return fmt.Errorf("failed to add watch directories: %v", err)
	}

	color.Green("✨ Initial build starting...")
	if err := w.buildAndRun(); err != nil {
		color.Red("❌ Initial build failed: %v", err)
	}

	go w.handleFileChanges()
	go w.handleUserInput()
	startLiveReloadServer()

	return nil
}

func (w *Watcher) addWatchDirs() error {
	for _, dir := range w.watchConfig.IncludeDirs {
		absDir := filepath.Join(w.dirPath, dir)
		err := filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				for _, excludeDir := range w.watchConfig.ExcludeDirs {
					if strings.Contains(path, excludeDir) {
						return filepath.SkipDir
					}
				}
				return w.fsWatcher.Add(path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Watcher) handleFileChanges() {
	for {
		select {
		case event := <-w.fsWatcher.Events:
			if w.shouldTriggerBuild(event) {
				w.mutex.Lock()
				if time.Since(w.lastChange) < w.debounce {
					w.mutex.Unlock()
					continue
				}
				w.lastChange = time.Now()
				w.mutex.Unlock()

				color.Yellow("🔄 File changed: %s", filepath.Base(event.Name))
				if err := w.buildAndRun(); err != nil {
					color.Red("❌ Build failed: %v", err)
				}
			}
		case err := <-w.fsWatcher.Errors:
			color.Red("⚠️ Watcher error: %v", err)
		}
	}
}

func (w *Watcher) shouldTriggerBuild(event fsnotify.Event) bool {
	if event.Op&fsnotify.Write != fsnotify.Write {
		return false
	}
	ext := filepath.Ext(event.Name)
	for _, validExt := range w.watchConfig.Extensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func (w *Watcher) buildAndRun() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	buildStart := time.Now()
	w.buildCount++

	if w.isRunning {
		color.Yellow("⏹️ Stopping previous process...")
		if err := w.stop(); err != nil {
			return fmt.Errorf("failed to stop process: %v", err)
		}
	}

	color.Blue("🏗️ Building...")

	outputFile := filepath.Join(w.outputDir, "app")
	if runtime.GOOS == "windows" {
		outputFile += ".exe"
	}

	// Run custom scripts
	for _, script := range w.buildConfig.Scripts {
		cmd := exec.Command("sh", "-c", script)
		cmd.Dir = w.dirPath
		cmd.Env = append(os.Environ(), w.buildConfig.Env...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("script failed: %v\n%s", err, output)
		}
	}

	buildCmd := exec.Command("go", append([]string{"build", "-o", outputFile}, w.buildConfig.Flags...)...)
	buildCmd.Dir = w.dirPath
	buildCmd.Env = append(os.Environ(), w.buildConfig.Env...)

	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build failed: %v\n%s", err, output)
	}

	w.lastBuildDuration = time.Since(buildStart)
	color.Green("✅ Build successful (took %v)", w.lastBuildDuration)

	w.cmd = exec.Command(outputFile)
	w.cmd.Stdout = os.Stdout
	w.cmd.Stderr = os.Stderr
	w.cmd.Env = append(os.Environ(), w.buildConfig.Env...)

	if err := w.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %v", err)
	}

	w.isRunning = true
	color.HiGreen("🚀 Process started (build #%d)", w.buildCount)

	go func() {
		w.cmd.Wait()
		w.isRunning = false
	}()

	return nil
}

func (w *Watcher) stop() error {
	if w.cmd != nil && w.cmd.Process != nil {
		if runtime.GOOS == "windows" {
			return w.cmd.Process.Kill()
		}
		return w.cmd.Process.Signal(os.Interrupt)
	}
	return nil
}

func (w *Watcher) handleUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		switch command {
		case "rs":
			color.Yellow("🔄 Manual restart triggered")
			if err := w.buildAndRun(); err != nil {
				color.Red("❌ Restart failed: %v", err)
			}
		case "stats":
			w.printStats()
		}
	}
}

func (w *Watcher) printStats() {
	color.HiCyan("\n📊 Watcher Statistics")
	color.HiCyan("-------------------")
	color.HiWhite("Total Runtime: %v", time.Since(w.startTime))
	color.HiWhite("Total Builds: %d", w.buildCount)
	color.HiWhite("Last Build Duration: %v", w.lastBuildDuration)
	color.HiWhite("Watching Extensions: %v", w.watchConfig.Extensions)
	color.HiWhite("Excluded Directories: %v", w.watchConfig.ExcludeDirs)
}

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
=======
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
>>>>>>> 4a08f121f7a71f3eb64a4ac6bbf605c6dedd6bfd
