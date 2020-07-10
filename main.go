package main

import (
	"flag"
	"fmt"
	"time"
)

func watchMboxFiles(watchedFileList WatchedFileList, allRules map[string]FileRuleSets, 
	newFilesToWatchCh chan NewFileToWatch, mailSender MailSender, stopChannel chan bool){

	// add new files to watch from (channel NewFileToWatch)
	go func() {
		for newFile := range newFilesToWatchCh {
			ruleSet := &allRules[newFile.RulesetId]
			newWatchedFile := WatchedFile{
				Filename: newFile.Filename
				RealEmail: newFile.RealEmail
				Rules: ruleSet
			}
			watchedFileList = append(watchedFileList, newWatchedFile)
		}
	}()

	filesToParseCh, err := WatchWorkingDir(appConfig.WorkingDir, &watchedFileList, stopChannel)
	mboxReader := &MboxReader{}

	for filesToParse := range filesToParseCh {
		for _, fileToParse := range filesToParse {
			mboxReader.setFilePath(fileToParse.Filename)
			rulesets := fileToParse.Rules
			for _, ruleset := range rulesets {
				mboxReader.resetFilters()
				for _, headerRule := range ruleset.HeaderRules {
					switch ruleType := headerRule.Type; ruleType {
					case "fromdate":
						mboxReader.setFromTime(time.Parse(time.RFC3339, headerRule.Value))
					case "beforedate":
						mboxReader.setBeforeTime(time.Parse(time.RFC3339, headerRule.Value))
					case "plain": 
						mboxReader.withHeader(headerRule.Key, headerRule.Value)
					case "regex":
						mboxReader.withHeader(headerRule.Key, headerRule.Value)
					}
				}
				msg, err := mboxReader.Read()
				for err == nil {
					mailSender(msg.getRawContents(), fileToParse.RealEmail)
					msg, err = mboxReader.Read()
				}
			}
		}
	}
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yml", "A path to a yaml configuration file.")
	flag.Parse()

	appConfig := ParseAppConfig(configPath)

	stopChannel := make(chan bool)
	watchedFileList, allRules := JsonFileLoader(appConfig)
	newFilesToWatchCh := LoadNewFilesFromTextFile(appConfig.NewFilesSourceTextfile, stopChannel)

	watchMboxFiles(watchedFileList, allRules, newFilesToWatchCh, sendmailSender, stopChannel)
}