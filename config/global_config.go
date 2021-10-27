package config

import (
	"fmt"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"reflect"
)

type globalConfig struct {
	Viper                    *viper.Viper
	DefaultProjectName       string
	ConfigFileName           string
	ConfigFileType           string
	PlaceholderRegex         string
	PlaceholderPrefix        string
	PlaceholderSuffix        string
	VariableNameRegex        string
	TemplateRootKeyName      string
	TemplateLocationKeyName  string
	VariablesKeyName         string
	GlobalsKeyName           string
	QuestionsKeyName         string
	DefaultAnswerKeyName     string
	DefaultChoicesKeyName    string
	ValidationKeyName        string
	LogLevel                 string
	SSHKeyLocation           string
	Aliases                  map[string]interface{}
	AliasesKeyName           string
	AliasesUrlKey            string
	AliasesLocalIndicatorKey string
}

func Global() *globalConfig {
	v := loadConfig()
	this := globalConfig{
		Viper:                    v,
		DefaultProjectName:       v.GetString("DefaultProjectName"),
		ConfigFileName:           v.GetString("ConfigFileName"),
		ConfigFileType:           v.GetString("ConfigFileType"),
		PlaceholderRegex:         v.GetString("PlaceholderRegex"),
		PlaceholderPrefix:        v.GetString("PlaceholderPrefix"),
		PlaceholderSuffix:        v.GetString("PlaceholderSuffix"),
		VariableNameRegex:        v.GetString("VariableNameRegex"),
		TemplateRootKeyName:      v.GetString("TemplateRootKeyName"),
		TemplateLocationKeyName:  v.GetString("TemplateLocationKeyName"),
		VariablesKeyName:         v.GetString("VariablesKeyName"),
		GlobalsKeyName:           v.GetString("GlobalsKeyName"),
		QuestionsKeyName:         v.GetString("QuestionsKeyName"),
		DefaultAnswerKeyName:     v.GetString("DefaultAnswerKeyName"),
		DefaultChoicesKeyName:    v.GetString("DefaultChoicesKeyName"),
		ValidationKeyName:        v.GetString("ValidationKeyName"),
		LogLevel:                 v.GetString("LogLevel"),
		SSHKeyLocation:           v.GetString("SSHKeyLocation"),
		Aliases:                  v.GetStringMap("Aliases"),
		AliasesKeyName:           v.GetString("AliasesKeyName"),
		AliasesUrlKey:            v.GetString("AliasesUrlKey"),
		AliasesLocalIndicatorKey: v.GetString("AliasesLocalIndicatorKey"),
	}
	return &this
}

func loadConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName(".clonr-config.yml")
	v.SetConfigType("yaml")
	v.AddConfigPath(utils.GetHomeDir())
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
	location := utils.GetHomeDir() + Global().ConfigFileName
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
	viperFunc("ValidationKeyName", "validation")
	viperFunc("LogLevel", "info")
	viperFunc("SSHKeyLocation", getSshLocation())
	viperFunc("AliasesKeyName", "aliases")
	viperFunc("AliasesUrlKey", "url")
	viperFunc("AliasesLocalIndicatorKey", "local")

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
		err = v.BindEnv("ValidationKeyName", "CLONR_VALIDATION_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("LogLevel", "CLONR_LOG")
		utils.ExitIfError(err)
		err = v.BindEnv("SSHKeyLocation", "CLONR_SSH_PATH")
		utils.ExitIfError(err)
		err = v.BindEnv("AliasesKeyName", "CLONR_ALIAS_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("AliasesUrlKey", "CLONR_ALIAS_URL_KEY")
		utils.ExitIfError(err)
		err = v.BindEnv("AliasesLocalIndicatorKey", "CLONR_ALIAS_LOCAL_KEY")
		utils.ExitIfError(err)
	}

}

func SetPropertyAndSave(propertyName string, value interface{}) {
	v := Global().Viper
	log.Debugf("Setting property: %s...", propertyName)
	switch propertyName {
	case "DefaultProjectName":
		v.Set("DefaultProjectName", value)
		break
	case "PlaceholderRegex":
		v.Set("PlaceholderRegex", value)
		break
	case "PlaceholderPrefix":
		v.Set("PlaceholderPrefix", value)
		break
	case "PlaceholderSuffix":
		v.Set("PlaceholderSuffix", value)
		break
	case "TemplateRootKeyName":
		v.Set("TemplateRootKeyName", value)
		break
	case "TemplateLocationKeyName":
		v.Set("TemplateLocationKeyName", value)
		break
	case "VariablesKeyName":
		v.Set("VariablesKeyName", value)
		break
	case "GlobalsKeyName":
		v.Set("GlobalsKeyName", value)
		break
	case "QuestionsKeyName":
		v.Set("QuestionsKeyName", value)
		break
	case "DefaultAnswerKeyName":
		v.Set("DefaultAnswerKeyName", value)
		break
	case "DefaultChoicesKeyName":
		v.Set("DefaultChoicesKeyName", value)
		break
	case "ValidationKeyName":
		v.Set("ValidationKeyName", value)
		break
	case "LogLevel":
		v.Set("LogLevel", value)
		break
	case "SSHKeyLocation":
		v.Set("SSHKeyLocation", value)
		break
	case "Aliases":
		v.Set("Aliases", value)
		break
	case "AliasesKeyName":
		v.Set("AliasesKeyName", value)
		break
	case "AliasesUrlKey":
		v.Set("AliasesUrlKey", value)
		break
	case "AliasesLocalIndicatorKey":
		v.Set("AliasesLocalIndicatorKey", value)
		break
	default:
		fmt.Println()
		log.Errorf("%s is not a clonr property or cannot be configured.", propertyName)
		fmt.Printf("\nRun 'clonr config show' for a list of options. No changes were made\n")
	}
	utils.SaveConfig(v, fmt.Sprintf("%s/%s", utils.GetHomeDir(), Global().ConfigFileName))
}

func ForEachConfigField(mutator *types.ConfigFieldMutator) {
	conf := *Global()
	value := reflect.ValueOf(conf)
	typeConf := value.Type()
	for i := 0; i < value.NumField(); i++ {
		mutator.Property = typeConf.Field(i).Name
		mutator.Value = cast.ToString(value.Field(i))
		// Add any properties you don't want available for bulk manipulation to this conditional.
		if mutator.Property != "Viper" &&
			mutator.Property != "ConfigFileName" &&
			mutator.Property != "ConfigFileType" &&
			mutator.Property != "Aliases" {
			mutator.ConfigMutator(mutator)
		}
	}
	mutator.Callback(mutator)
}

func getSshLocation() string {
	return utils.GetHomeDir() + "/.ssh/id_rsa"
}
