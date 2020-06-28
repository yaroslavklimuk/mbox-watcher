package main

import "os"

func getFilesToParse(files []os.FileInfo, watchedFiles []string) []string {
	filesCount := len(files)
	var result = make([]string, filesCount)

	for _, fileInfo := range files {
		if fileInfo.Size() > 0 && stringInSlice(fileInfo.Name(), watchedFiles) {
			result = append(result, fileInfo.Name())
		}
	}
	return result
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func WatchWorkingDir(workDir string, watchedFiles []string, newFileCh chan string,
	parseFilesCh chan []string, stopChannel chan bool) {

	workDirItem, openErr := os.Open(workDir)
	if openErr != nil {
		panic(openErr)
	}
	defer workDirItem.Close()
	defer close(parseFilesCh)

	for {
		select {
		case stopSignal := <-stopChannel:
			if stopSignal == true {
				return
			}
		default:
			dirContents, readErr := workDirItem.Readdir(-1)
			if readErr != nil {
				panic(readErr)
			}
			toParse := getFilesToParse(dirContents, watchedFiles)
			parseFilesCh <- toParse
		}

	}
}
