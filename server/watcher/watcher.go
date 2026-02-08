package watcher

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	rootDir := os.Getenv("WATCH_DIR")
	if rootDir == "" {
		rootDir = "/app/data" // default container path
	}

	if err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			log.Printf("Watching: %s", path)
			return watcher.Add(path)
		}
		return nil
	}); err != nil {
		log.Fatal("Error walking directory:", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() {
						// if directory add to watcher
						log.Printf("New directory detected: %s (watching now)\n", event.Name)
						watcher.Add(event.Name)
					} else {
						log.Printf("New file created: %s\n", event.Name)
					}
				} else if event.Has(fsnotify.Write) {
					log.Printf("File modified: %s\n", event.Name)
				} else if event.Has(fsnotify.Remove) {
					log.Printf("Item removed: %s\n", event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	log.Printf("Recursive watcher watching %s.", rootDir)
	<-done
}
