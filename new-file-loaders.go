package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
)

type NewFileToWatch struct {
	Filename  string `json:"filename"`
	RulesetId string `json:"ruleset-id"`
	RealEmail string `json:"real-email"`
}

func LoadNewFilesFromTextFile(newFilesSource string, stopChannel chan bool) chan NewFileToWatch {
	// we make a channel to send info about every new file to it
	watchedFilesCh := make(chan NewFileToWatch)
	var newFileToWatch NewFileToWatch

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		defer close(watchedFilesCh)

		watcher.Add(newFilesSource)
		var prevContents = make([]byte, 0)
		for {
			select {
			case event, ok := <-watcher.Events:
				if ok {
					if event.Op == fsnotify.Write {
						jsonContents, _ := ioutil.ReadFile(newFilesSource)
						// fmt.Printf("new:%s\nold:%s\n", string(jsonContents), string(prevContents))
						comp := bytes.Compare(prevContents, jsonContents)
						if comp != 0 {
							json.Unmarshal(jsonContents, &newFileToWatch)
							watchedFilesCh <- newFileToWatch
						}
						prevContents = jsonContents
					}
				}
			case stopSignal := <-stopChannel:
				if stopSignal == true {
					return
				}
			}
		}
	}()

	return watchedFilesCh
}
