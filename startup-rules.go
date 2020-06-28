package main

type HeaderRule struct {
	Key   string
	Value string
	Type  string
}

type RuleSet struct {
	HeaderRules []HeaderRule
}

type FileRuleSets []RuleSet

type WatchedFile struct {
	Filename  string
	RealEmail string
	Rules     *FileRuleSets
}

type WatchedFileList []WatchedFile
