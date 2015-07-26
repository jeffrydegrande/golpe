package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

func Watch() error {
	watcher, err := fsnotify.NewWatcher()
	check(err)
	defer watcher.Close()

	lastBuildTime := time.Now()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}

				timeAgo := time.Now().Sub(lastBuildTime)
				fmt.Printf("Last build was %f seconds ago\n", timeAgo.Seconds())
				if event.Name != "public" && timeAgo > 1 {
					var m = fmt.Sprintf("Changed detected in %s, rebuilding\n", event.Name)
					say(m)
					BuildAll()
				}
				lastBuildTime = time.Now()
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

	err = watcher.Add("stylesheets")
	check(err)

	<-done

	return nil
}
