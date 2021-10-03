package config

type config struct {
	DefaultProjectName string
	ClonrConfigFileName string
	ClonrConfigFileExt string
	ClonrRegex string
	ClonrPrefix string
	ClonrSuffix string
	VariableRegex string
}


func GlobalConfig() *config {
	this := config{
		DefaultProjectName: "clonr-app",
		ClonrConfigFileName: "clonr",
		ClonrConfigFileExt: "yaml",
		ClonrRegex:  "\\{{1}@{1}clonr\\{{1}[\\w-]+\\}{2}",
		ClonrPrefix: "{@clonr{",
		ClonrSuffix: "}}",
		VariableRegex: "[\\w-]+",
	}
	return &this
}