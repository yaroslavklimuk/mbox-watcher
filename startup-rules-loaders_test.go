package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestJsonFileLoader(t *testing.T) {
	type FileToRuleset struct {
		Filename  string
		RulesetId string
	}

	type JsonFileLoaderTestCase struct {
		RulesFile       string
		FilesToRulesets map[string]string
		JsonRuleset     string
	}

	var jsonFileLoaderTestTable = []JsonFileLoaderTestCase{
		JsonFileLoaderTestCase{
			RulesFile: "test-startup-rules.json",
			FilesToRulesets: map[string]string{
				"uhdfuighdf@order.penguins.ru": "76f876fgdg",
				"7f89d7g98@order.penguins.ru":  "76f876fgdg",
			},
			JsonRuleset: `[{"HeaderRules":[{"Key":"From","Value":"noreply@ofd.org","Type":"plain"},{"Key":"Subject","Value":"/^.*Your order receipt.*$/","Type":"regex"}]}]`,
		},
		JsonFileLoaderTestCase{
			RulesFile: "test-startup-rules-2.json",
			FilesToRulesets: map[string]string{
				"uhdfuighdf@order.penguins.ru": "4546576576",
				"7f89d7g98@order.penguins.ru":  "4546576576",
			},
			JsonRuleset: `[{"HeaderRules":[{"Key":"From","Value":"noreply@ofd.org","Type":"plain"},{"Key":"Subject","Value":"/^.*Your order receipt.*$/","Type":"regex"}]},{"HeaderRules":[{"Key":"Subject","Value":"/^.*Your order receipt.*$/","Type":"regex"},{"Key":"Date","Value":"2020-04-03 00:00:00","Type":"fromdate"},{"Key":"Date","Value":"2020-04-19 13:14:15","Type":"beforedate"}]}]`,
		},
	}

	conf := ParseAppConfig("config.yml")
	conf.RulesSource = "rulesfilejson"

	wd, _ := os.Getwd()
	for _, tcase := range jsonFileLoaderTestTable {
		t.Run(tcase.RulesFile, func(t *testing.T) {
			conf.RulesFileJson.Path = wd + string(os.PathSeparator) + "testcases" +
				string(os.PathSeparator) + tcase.RulesFile
			actualFileList, actualRules := JsonFileLoader(conf)

			if len(actualFileList) != len(tcase.FilesToRulesets) {
				t.Errorf("number of files error")
			}

			for _, actualFile := range actualFileList {
				expectedRulesetId := tcase.FilesToRulesets[actualFile.Filename]

				expectedRulesetJson, _ := json.Marshal(actualRules[expectedRulesetId])
				actualRulesetJson, _ := json.Marshal(actualFile.Rules)

				expectedRuleset := string(expectedRulesetJson)
				actualRuleset := string(actualRulesetJson)

				if actualRuleset != tcase.JsonRuleset {
					t.Errorf("Watched file ruleset structure is not correct. Got %v, want %v", actualRuleset, tcase.JsonRuleset)
				}

				if expectedRuleset != actualRuleset {
					t.Errorf("Watched file ruleset pointer is not correct. Got %v, want %v", actualRuleset, expectedRuleset)
				}
			}
		})
	}
}
