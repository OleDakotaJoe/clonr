package config

import (
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"reflect"
)

type globalConfig struct {
	Viper                   *viper.Viper
	DefaultProjectName      string
	ConfigFileName          string
	ConfigFileType          string
	PlaceholderRegex        string
	PlaceholderPrefix       string
	PlaceholderSuffix       string
	VariableNameRegex       string
	TemplateRootKeyName     string
	TemplateLocationKeyName string
	VariablesKeyName        string
	GlobalsKeyName          string
	QuestionsKeyName        string
	DefaultAnswerKeyName    string
	DefaultChoicesKeyName   string
	LogLevel                string
}

func Global() *globalConfig {
	v := getConfig()
	this := globalConfig{
		Viper:                   v,
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
	v.SetConfigName(".clonr-config.yml")
	v.SetConfigType("yaml")
	v.AddConfigPath(utils.GetLocationOfInstalledBinary())
	initDefaults(v)
	err := v.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Debug("Using default configuration.")
	} else {
		utils.ExitIfError(err)
	}
	return v
}

func initDefaults(v *viper.Viper) {
	setDefaults(v, v.SetDefault)
}

func ResetGlobalToDefault() {
	v := Global().Viper
	location := utils.GetLocationOfInstalledBinary() + Global().ConfigFileName
	setDefaults(v, v.Set)
	utils.SaveConfig(v, location)
}

func setDefaults(v *viper.Viper, viperFunc func(key string, value interface{})) {
	viperFunc("DefaultProjectName", "clonr-app")
	viperFunc("ConfigFileName", ".clonr-config.yml")
	viperFunc("ConfigFileType", "yaml")
	viperFunc("PlaceholderRegex", "\\{{1}@{1}clonr\\{{1}[a-z0-9-_]+\\}{2}")
	viperFunc("PlaceholderPrefix", "{@clonr{")
	viperFunc("PlaceholderSuffix", "}}")
	viperFunc("VariableNameRegex", "[\\w-]+")
	viperFunc("TemplateRootKeyName", "templates")
	viperFunc("TemplateLocationKeyName", "location")
	viperFunc("VariablesKeyName", "variables")
	viperFunc("GlobalsKeyName", "globals")
	viperFunc("QuestionsKeyName", "question")
	viperFunc("DefaultAnswerKeyName", "default")
	viperFunc("DefaultChoicesKeyName", "choices")
	viperFunc("LogLevel", "info")

	if reflect.TypeOf(viperFunc) == reflect.TypeOf(viper.GetViper().SetDefault) {
		var err error
		err = v.BindEnv("DefaultProjectName", "CLONR_DFLT_PROJ_NAME")
		utils.ExitIfError(err)
		err = v.BindEnv("ConfigFileName", "CLONR_CONFIG_FNAME")
		utils.ExitIfError(err)
		err = v.BindEnv("ConfigFileType", "CLONR_CONFIG_FTYPE")
		utils.ExitIfError(err)
		err = v.BindEnv("PlaceholderRegex", "CLONR_PH_REGEX")
		utils.ExitIfError(err)
		err = v.BindEnv("PlaceholderPrefix", "CLONR_PH_PREFIX")
		utils.ExitIfError(err)
		err = v.BindEnv("PlaceholderSuffix", "CLONR_PH_SUFFIX")
		utils.ExitIfError(err)
		err = v.BindEnv("VariableNameRegex", "CLONR_VARS_REGEX")
		utils.ExitIfError(err)
		err = v.BindEnv("TemplateRootKeyName", "CLONR_TEMPLATE_ROOT_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("TemplateLocationKeyName", "CLONR_TEMPLATE_LOC_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("VariablesKeyName", "CLONR_VARS_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("GlobalsKeyName", "CLONR_GLOBALS_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("QuestionsKeyName", "CLONR_QUES_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("DefaultAnswerKeyName", "CLONR_DFLT_ANS_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("DefaultChoicesKeyName", "CLONR_CHOICES_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("LogLevel", "CLONR_LOG")
		utils.ExitIfError(err)
	}

}
