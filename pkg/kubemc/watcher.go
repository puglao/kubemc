package kubemc

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func WatchMCEvents(dirname string, e chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) > 0 {
					// log.Println("file changed:", event.Name)
					e <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(dirname)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
