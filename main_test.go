package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

func TestWatchMboxFiles(t *testing.T) {

	type WatchMboxFilesTestCase struct {
		RuleSet []struct {
			HeaderRules [][]struct {
				Key   string `json:"key"`
				Value string `json:"value"`
				Type  string `json:"type"`
			} `json:"headerrules"`
			RulesetId string `json:"rulesetid"`
			Files     []struct {
				Name      string `json:"name"`
				RealEmail string `json:"real-email"`
			} `json:"files"`
		} `json:"rules"`
		MessagesToSend    map[string][]string `json:"messages-to-send"`
		ForwardedMessages map[string][]string `json:"forwarded-messages-md5"`
	}

	testTable := make([]WatchMboxFilesTestCase, 0)
	data, err := ioutil.ReadFile("testcases/load-new-file-from-textfile-test-table.json")
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(data, &testTable)

	for ind, tcase := range testTable {
		t.Run(string(ind), func(t *testing.T) {

		})
	}
}
