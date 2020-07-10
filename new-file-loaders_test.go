package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"testing"
	"time"
)

func TestLoadNewFilesFromTextFile(t *testing.T) {

	type LoadNewFilesFromTextFileTestCase struct {
		NewFiles []NewFileToWatch `json:"files"`
	}

	testTable := make([]LoadNewFilesFromTextFileTestCase, 0)
	data, err := ioutil.ReadFile("testcases/load-new-file-from-textfile-test-table.json")
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(data, &testTable)

	for ind, tcase := range testTable {
		t.Run(string(ind), func(t *testing.T) {
			stopChannel := make(chan bool)
			newFileChannel := LoadNewFilesFromTextFile("new-file-to-watch.txt", stopChannel)

			var wg sync.WaitGroup
			var receivedNewFiles []NewFileToWatch

			wg.Add(1)
			go func() {
				defer wg.Done()
				for readNewFile := range newFileChannel {
					receivedNewFiles = append(receivedNewFiles, readNewFile)
				}
			}()

			for _, writeNewFile := range tcase.NewFiles {
				jsonToWrite, _ := json.Marshal(writeNewFile)
				locErr := ioutil.WriteFile("new-file-to-watch.txt", jsonToWrite, 0644)
				if locErr != nil {
					log.Println(locErr)
				}
				time.Sleep(1 * time.Second)
			}
			stopChannel <- true
			wg.Wait()

			inputNewFilesJSON, _ := json.Marshal(tcase.NewFiles)
			receivedNewFilesJSON, outpErr := json.Marshal(receivedNewFiles)

			if outpErr != nil {
				t.Errorf("bad output format: %v", outpErr)
			}

			comp := bytes.Compare(inputNewFilesJSON, receivedNewFilesJSON)
			if comp != 0 {
				if string(inputNewFilesJSON) != string(receivedNewFilesJSON) {
					t.Logf("want: %s\ngot: %s\n", string(inputNewFilesJSON), string(receivedNewFilesJSON))
					t.Fail()
				}
			}
			receivedNewFiles = nil
		})
	}
}
