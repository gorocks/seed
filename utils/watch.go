package utils

import (
	"github.com/Guazi-inc/seed/logger"
	"github.com/Guazi-inc/seed/logger/color"
	"regexp"
	"strings"
	"github.com/fsnotify/fsnotify"
)

var (
	watchExts           = []string{".go", ".json"}
	watchExtsStatic     = []string{".html", ".tpl", ".js", ".css"}
	ignoredFilesRegExps = []string{
		`.#(\w+).go`,
		`.(\w+).go.swp`,
		`(\w+).go~`,
		`(\w+).tmp`,
	}
)

//在监控文件的时候，如果文件发生改变执行的方法
type DoWatch interface {
	Exec(paths []string, files []string, name string)
}

// NewWatcher starts an fsnotify Watcher on the specified paths
func NewWatcher(paths []string, files []string, do DoWatch) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Log.Fatalf("Failed to create watcher: %s", err)
	}
	go func() {
		for {
			select {
			case e := <-watcher.Events:
				// Skip ignored files
				if shouldIgnoreFile(e.Name) {
					continue
				}
				if !shouldWatchFileWithExtension(e.Name) {
					continue
				}
				do.Exec(paths, files, e.Name)
			case err := <-watcher.Errors:
				logger.Log.Warnf("Watcher error: %s", err.Error()) // No need to exit here
			}
		}
	}()

	logger.Log.Info("Initializing watcher...")
	for _, path := range paths {
		logger.Log.Hintf(colors.Bold("Watching: ")+"%s", path)
		err = watcher.Add(path)
		if err != nil {
			logger.Log.Fatalf("Failed to watch directory: %s", err)
		}
	}
}

// shouldIgnoreFile ignores filenames generated by Emacs, Vim or SublimeText.
// It returns true if the file should be ignored, false otherwise.
func shouldIgnoreFile(filename string) bool {
	for _, regex := range ignoredFilesRegExps {
		r, err := regexp.Compile(regex)
		if err != nil {
			logger.Log.Fatalf("Could not compile regular expression: %s", err)
		}
		if r.MatchString(filename) {
			return true
		}
		continue
	}
	return false
}

// shouldWatchFileWithExtension returns true if the name of the file
// hash a suffix that should be watched.
func shouldWatchFileWithExtension(name string) bool {
	for _, s := range watchExts {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

func ifStaticFile(filename string) bool {
	for _, s := range watchExtsStatic {
		if strings.HasSuffix(filename, s) {
			return true
		}
	}
	return false
}
