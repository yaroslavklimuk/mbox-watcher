package main

import (
	"flag"
	"fmt"
)

func getWatchedFilenames(fileList WatchedFileList) []string {
	filesCount := len(fileList)
	var result = make([]string, filesCount)

	for _, fileInfo := range fileList {
		result = append(result, fileInfo.Filename)
	}
	return result
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yml", "A path to a yaml configuration file.")
	flag.Parse()

	appConfig := ParseAppConfig(configPath)

	stopChannel := make(chan bool)
	watchedFileList, _ := JsonFileLoader(appConfig)
	newFilesCh := LoadNewFilesFromTextFile(appConfig, stopChannel)
	var newFilenamesCh = make(chan string)

	filesToParseCh := make(chan []string)

	go WatchWorkingDir(appConfig.WorkingDir, getWatchedFilenames(watchedFileList),
		newFilenamesCh, filesToParseCh, stopChannel)

	go func() {
		for newFile := range newFilesCh {
			newFilenamesCh <- newFile.Filename
		}
	}()

	for toParse := range filesToParseCh {
		fmt.Printf("to parse: %v\n", toParse)
	}

}
