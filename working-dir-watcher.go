package main

import "os"

func getFilesToParse(files []os.FileInfo, watchedFiles *WatchedFileList) WatchedFileList {
	filesCount := len(files)
	result := WatchedFileList{}
	for _, fileInfo := range files {
		if fileInfo.Size() > 0 {
			watchedFile := getWatchedFileByName(fileInfo.Name(), watchedFiles)
			if watchedFile != nil {
				result = append(result, watchedFile)
			}
		}
	}
	return result
}

func getWatchedFileByName(name string, watchedFiles WatchedFileList) WatchedFile {
	for _, watchedFile := range watchedFiles {
		if watchedFile.Filename == name {
			return watchedFile
		}
	}
	return nil
}

func WatchWorkingDir(workDir string, watchedFiles *WatchedFileList, stopChannel chan bool) (chan WatchedFileList, error) {
	workDirItem, openErr := os.Open(workDir)
	if openErr != nil {
		return nil, openErr
	}
	defer workDirItem.Close()

	filesToParseCh := make(chan WatchedFileList)
	defer close(filesToParseCh)

	go func() {
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
				filesToParseCh <- toParse
			}
		}
	}()

	return filesToParseCh
}
