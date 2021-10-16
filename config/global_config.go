package config

import (
	"github.com/oledakotajoe/clonr/utils"
	"github.com/spf13/viper"
)

type globalConfig struct {
	DefaultProjectName          string
	ClonrConfigFileName         string
	ClonrConfigFileType         string
	ClonrPlaceholderRegex       string
	ClonrPrefix                 string
	ClonrSuffix                 string
	ClonrVariableRegex          string
	ClonrConfigRootKeyName      string
	TemplateFileLocationKeyName string
	VariablesArrayKeyName       string
	GlobalVariablesKeyName      string
	QuestionsKeyName            string
	DefaultAnswerKeyName        string
	DefaultChoicesKeyName       string
	LogLevel                    string
}

func Global() *globalConfig {
	v, err := utils.ViperReadConfig(utils.GetLocationOfInstalledBinary(), ".clonr-config.yml", "yaml")

	v.SetDefault("DefaultProjectName", "clonr-app")
	v.SetDefault("ClonrConfigFileName", ".clonrrc")
	v.SetDefault("ClonrConfigFileType", "yaml")
	v.SetDefault("ClonrPlaceholderRegex", "\\{{1}@{1}clonr\\{{1}[a-z0-9-_]+\\}{2}")
	v.SetDefault("ClonrPrefix", "{@clonr{")
	v.SetDefault("ClonrSuffix", "}}")
	v.SetDefault("ClonrVariableRegex", "[\\w-]+")
	v.SetDefault("ClonrConfigRootKeyName", "templates")
	v.SetDefault("TemplateFileLocationKeyName", "location")
	v.SetDefault("VariablesArrayKeyName", "variables")
	v.SetDefault("GlobalVariablesKeyName", "globals")
	v.SetDefault("QuestionsKeyName", "question")
	v.SetDefault("DefaultAnswerKeyName", "default")
	v.SetDefault("DefaultChoicesKeyName", "choices")
	v.SetDefault("LogLevel", "INFO")

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			writeErr := v.WriteConfig()
			utils.CheckForError(writeErr)
		} else {
			// Config file was found but another error was produced, exit status 1
			utils.CheckForError(err)
		}
	}

	this := globalConfig{
		DefaultProjectName:          v.GetString("DefaultProjectName"),
		ClonrConfigFileName:         v.GetString("ClonrConfigFileName"),
		ClonrConfigFileType:         v.GetString("ClonrConfigFileType"),
		ClonrPlaceholderRegex:       v.GetString("ClonrPlaceholderRegex"),
		ClonrPrefix:                 v.GetString("ClonrPrefix"),
		ClonrSuffix:                 v.GetString("ClonrSuffix"),
		ClonrVariableRegex:          v.GetString("ClonrVariableRegex"),
		ClonrConfigRootKeyName:      v.GetString("ClonrConfigRootKeyName"),
		TemplateFileLocationKeyName: v.GetString("TemplateFileLocationKeyName"),
		VariablesArrayKeyName:       v.GetString("VariablesArrayKeyName"),
		GlobalVariablesKeyName:      v.GetString("GlobalVariablesKeyName"),
		QuestionsKeyName:            v.GetString("QuestionsKeyName"),
		DefaultAnswerKeyName:        v.GetString("DefaultAnswerKeyName"),
		DefaultChoicesKeyName:       v.GetString("DefaultChoicesKeyName"),
		LogLevel:                    v.GetString("LogLevel"),
	}
	return &this
}
