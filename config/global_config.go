package config

type config struct {
	DefaultProjectName          string
	ClonrConfigFileName         string
	ClonrConfigFileType         string
	ClonrRegex                  string
	ClonrPrefix                 string
	ClonrSuffix                 string
	VariableRegex               string
	ClonrConfigRootKeyName      string
	TemplateFileLocationKeyName string
	VariablesArrayKeyName       string
}

func GlobalConfig() *config {
	this := config{
		DefaultProjectName:          "clonr-app",
		ClonrConfigFileName:         ".clonrrc",
		ClonrConfigFileType:         "yaml",
		ClonrRegex:                  "\\{{1}@{1}clonr\\{{1}[a-z0-9-_]+\\}{2}",
		ClonrPrefix:                 "{@clonr{",
		ClonrSuffix:                 "}}",
		VariableRegex:               "[\\w-]+",
		ClonrConfigRootKeyName:      "paths",
		TemplateFileLocationKeyName: "location",
		VariablesArrayKeyName:       "variables",
	}
	return &this
}
