package utils

import (
	"clonr/config"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"regexp"
)

func CheckForError(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func ThrowError(message string) error {
	err := errors.New(message)
	return err
}

func ViperReadConfig(configFilePath string, configFileName string) *viper.Viper {
	v := viper.GetViper()
	v.SetConfigName(configFileName)
	v.AddConfigPath(configFilePath)
	log.Debugf("Config File Location: %s", v.ConfigFileUsed())
	err := v.ReadInConfig()
	CheckForError(err)
	return v
}

func RemoveElementFromSlice(list []string, index int) []string {
	return append(list[:index], list[index+1:]...)
}

func IsVariableValid(variable string) (bool, error) {
	return regexp.Match(config.GlobalConfig().VariableRegex, []byte(variable))
}

func GetKeysFromMap(someMap map[string]string) []string {
	keys := make([]string, len(someMap))
	i := 0
	for k := range someMap {
		keys[i] = k
		i++
	}
	return keys
}
