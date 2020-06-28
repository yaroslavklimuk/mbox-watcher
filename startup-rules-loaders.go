package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type RawRulesJSON struct {
	Rules []struct {
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
}

func JsonFileLoader(c Config) (WatchedFileList, map[string]FileRuleSets) {
	jsonFile, err := os.Open(c.RulesFileJson.Path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	var tempRulesStruct RawRulesJSON

	jsonContents, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonContents, &tempRulesStruct)

	var watchedFileList WatchedFileList
	allRules := make(map[string]FileRuleSets)

	for _, rule := range tempRulesStruct.Rules {
		var tempRuleSetJar FileRuleSets
		for _, rawHeadRuleSet := range rule.HeaderRules {
			var ruleSet = RuleSet{
				HeaderRules: []HeaderRule{},
			}
			for _, rawHeadRule := range rawHeadRuleSet {
				headRule := HeaderRule(rawHeadRule)
				ruleSet.HeaderRules = append(ruleSet.HeaderRules, headRule)
			}
			tempRuleSetJar = append(tempRuleSetJar, ruleSet)
		}

		for _, rawFile := range rule.Files {
			var tempWatchedFile WatchedFile
			tempWatchedFile.Filename = rawFile.Name
			tempWatchedFile.RealEmail = rawFile.RealEmail
			tempWatchedFile.Rules = &tempRuleSetJar
			watchedFileList = append(watchedFileList, tempWatchedFile)
		}

		allRules[rule.RulesetId] = tempRuleSetJar
	}
	return watchedFileList, allRules
}

// func MysqlLoader(c Config) (WatchedFileList, map[string]FileRuleSets) {

// 	dsn := fmt.Sprintf("<%s>:<%s>@tcp(%s:%s)/<%s>",
// 		c.Database.Username, c.Database.Password, c.Database.Hostname, c.Database.Port, c.Database.Database)

// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// }
