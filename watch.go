package main

import (
	"fmt"
	"log"

	"github.com/go-fsnotify/fsnotify"
)

func watch() error {
	watcher, err := fsnotify.NewWatcher()
	check(err)
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}

				var m = fmt.Sprintf("Changed detected in %s, rebuilding\n", event.Name)
				say(m)
				buildAll()
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(".")
	check(err)

	err = watcher.Add("javascripts")
	check(err)

	err = watcher.Add("javascripts/components")
	check(err)

	<-done

	return nil
}
