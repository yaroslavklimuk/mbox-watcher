package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RulesSource string `yaml:"rules-source"`
	Database    struct {
		Hostname string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
		Database string `yaml:"db"`
	} `yaml:"database"`
	RulesFileJson struct {
		Path string `yaml:"path"`
	} `yaml: "rulesfilejson"`
	NewFilesSource         string `yaml:"new-files-source"`
	NewFilesSourceTextfile string `yaml:"new-files-source-textfile"`
	WorkingDir             string `yaml:"working-dir"`
}

func ParseAppConfig(configFile string) Config {
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
