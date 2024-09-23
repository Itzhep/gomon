package watcher

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"github.com/fsnotify/fsnotify"
	"github.com/fatih/color"
)

type Watcher struct {
	fsWatcher *fsnotify.Watcher
	cmd       *exec.Cmd
	appPath   string
	dirPath   string
	debounce  time.Duration
	lastChange time.Time
}

func NewWatcher(appPath string, debounce time.Duration) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Dir(appPath)

	return &Watcher{
		fsWatcher: fsWatcher,
		appPath:   appPath,
		dirPath:   dirPath,
		debounce:  debounce,
	}, nil
}

func (w *Watcher) StartApp() error {
	if w.cmd != nil && w.cmd.Process != nil {
		_ = w.cmd.Process.Kill() // Kill previous instance if running
	}
	color.Cyan("Starting app...")

	w.cmd = exec.Command("go", "run", w.appPath)
	w.cmd.Stdout = os.Stdout
	w.cmd.Stderr = os.Stderr

	return w.cmd.Start() // Start the application
}

func (w *Watcher) WatchAndReload() error {
	defer w.fsWatcher.Close()

	done := make(chan bool)

	err := filepath.Walk(w.dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.fsWatcher.Add(path)
		}
		return nil
	})
	if err != nil {
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

				if time.Since(w.lastChange) < w.debounce {
					continue
				}

				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 &&
					filepath.Ext(event.Name) == ".go" {
					w.lastChange = time.Now()
					color.Yellow("File changed: %s", event.Name)
					if err := w.StartApp(); err != nil {
						color.Red("Failed to restart the app: %v", err)
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

func (w *Watcher) handleUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		if command == "rs" {
			color.Cyan("Restarting app...")
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
