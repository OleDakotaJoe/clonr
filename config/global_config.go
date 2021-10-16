package config

import (
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type globalConfig struct {
	DefaultProjectName  string
	ConfigFileName        string
	ConfigFileType         string
	PlaceholderRegex       string
	PlaceholderPrefix  string
	PlaceholderSuffix   string
	VariableNameRegex   string
	TemplateRootKeyName string
	TemplateLocationKeyName string
	VariablesKeyName        string
	GlobalsKeyName              string
	QuestionsKeyName      string
	DefaultAnswerKeyName        string
	DefaultChoicesKeyName string
	LogLevel              string
}

func setDefaults(v *viper.Viper) {
	var err error
	v.SetDefault("DefaultProjectName", "clonr-app")
	err = v.BindEnv("DefaultProjectName", "CLONR_DFLT_PROJ_NAME"); utils.CheckForError(err)

	v.SetDefault("ConfigFileName", ".clonrrc")
	err = v.BindEnv("ConfigFileName", "CLONR_CONFIG_FNAME"); utils.CheckForError(err)

	v.SetDefault("ConfigFileType", "yaml")
	err = v.BindEnv("ConfigFileType", "CLONR_CONFIG_FTYPE"); utils.CheckForError(err)

	v.SetDefault("PlaceholderRegex", "\\{{1}@{1}clonr\\{{1}[a-z0-9-_]+\\}{2}")
	err = v.BindEnv("PlaceholderRegex", "CLONR_PH_REGEX"); utils.CheckForError(err)

	v.SetDefault("PlaceholderPrefix", "{@clonr{")
	err = v.BindEnv("PlaceholderPrefix", "CLONR_PH_PREFIX"); utils.CheckForError(err)

	v.SetDefault("PlaceholderSuffix", "}}")
	err = v.BindEnv("PlaceholderSuffix", "CLONR_PH_SUFFIX"); utils.CheckForError(err)

	v.SetDefault("VariableNameRegex", "[\\w-]+")
	err = v.BindEnv("VariableNameRegex", "CLONR_VARS_REGEX"); utils.CheckForError(err)

	v.SetDefault("TemplateRootKeyName", "templates")
	err = v.BindEnv("TemplateRootKeyName", "CLONR_TEMPLATE_ROOT_KEY"); utils.CheckForError(err)

	v.SetDefault("TemplateLocationKeyName", "location")
	err = v.BindEnv("TemplateLocationKeyName", "CLONR_TEMPLATE_LOC_KEY"); utils.CheckForError(err)

	v.SetDefault("VariablesKeyName", "variables")
	err = v.BindEnv("VariablesKeyName", "CLONR_VARS_KEY"); utils.CheckForError(err)

	v.SetDefault("GlobalsKeyName", "globals")
	err = v.BindEnv("GlobalsKeyName", "CLONR_GLOBALS_KEY"); utils.CheckForError(err)

	v.SetDefault("QuestionsKeyName", "question")
	err = v.BindEnv("QuestionsKeyName", "CLONR_QUES_KEY"); utils.CheckForError(err)

	v.SetDefault("DefaultAnswerKeyName", "default")
	err = v.BindEnv("DefaultAnswerKeyName", "CLONR_DFLT_ANS_KEY"); utils.CheckForError(err)

	v.SetDefault("DefaultChoicesKeyName", "choices")
	err = v.BindEnv("DefaultChoicesKeyName", "CLONR_CHOICES_KEY"); utils.CheckForError(err)

	v.SetDefault("LogLevel", "info")
	err = v.BindEnv("LogLevel", "CLONR_LOG"); utils.CheckForError(err)
}

func Global() *globalConfig {
	v := getConfig()
	this := globalConfig{
		DefaultProjectName:      v.GetString("DefaultProjectName"),
		ConfigFileName:          v.GetString("ConfigFileName"),
		ConfigFileType:          v.GetString("ConfigFileType"),
		PlaceholderRegex:        v.GetString("PlaceholderRegex"),
		PlaceholderPrefix:       v.GetString("PlaceholderPrefix"),
		PlaceholderSuffix:       v.GetString("PlaceholderSuffix"),
		VariableNameRegex:       v.GetString("VariableNameRegex"),
		TemplateRootKeyName:     v.GetString("TemplateRootKeyName"),
		TemplateLocationKeyName: v.GetString("TemplateLocationKeyName"),
		VariablesKeyName:        v.GetString("VariablesKeyName"),
		GlobalsKeyName:          v.GetString("GlobalsKeyName"),
		QuestionsKeyName:        v.GetString("QuestionsKeyName"),
		DefaultAnswerKeyName:    v.GetString("DefaultAnswerKeyName"),
		DefaultChoicesKeyName:   v.GetString("DefaultChoicesKeyName"),
		LogLevel:                v.GetString("LogLevel"),
	}
	return &this
}

func getConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName(utils.GetLocationOfInstalledBinary())
	v.SetConfigType("yaml")
	v.AddConfigPath(".clonr-config.yml")
	setDefaults(v)
	err := v.ReadInConfig()
	if err == err.(viper.ConfigFileNotFoundError) {
		log.Debug("Using default configuration.")
	} else {
		utils.CheckForError(err)
	}
	return v
}
