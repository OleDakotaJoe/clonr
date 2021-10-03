package core

import (
	"bufio"
	"clonr/config"
	"clonr/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

func ProcessFiles(configFilePath string) {
	v := ViperReadConfig(configFilePath)
	configRootKey := config.GlobalConfig().ClonrConfigRootKeyName
	paths := v.GetStringMap(configRootKey)
	log.Debugf("Paths: %s", paths)
	for path := range paths {
		log.Infof("Processing path: %s", path)
		filePathKey := configRootKey +"."+ path
		pathData := v.GetStringMapString(filePathKey)
		log.Debugf("Raw pathData: %s", pathData)
		fileLocation := pathData[config.GlobalConfig().TemplateFileLocationKeyName]
		variableArrayKey := filePathKey + "." + config.GlobalConfig().VariablesArrayKeyName
		variablesMap := v.Get(variableArrayKey)
		//processedVarMap := make(map[string]string)
		log.Debug(fileLocation, variablesMap)
		//for question, variable := range variablesMap {
		//	processedVarMap[variable] = answerQuestion(fmt.Sprintf("%s", question))
		//}
		//log.Debugf("Processed varMap: %s", processedVarMap)
		//// Renders the file below
		//log.Infof("Rendering file: %s, with vars: %s", path, variablesMap)
		//RenderFile(fileLocation, variablesMap)
	}
}

func RenderFile(filename string, varMap map[string]string) {
	input, err := ioutil.ReadFile(filename)
	utils.CheckForError(err)
	inputFileAsString := string(input)
	for key, value := range varMap {
		pattern := config.GlobalConfig().ClonrPrefix + key + config.GlobalConfig().ClonrSuffix
		log.Debugf("Line: %s, Pattern: %s, value: %s, key: %s", inputFileAsString, pattern, value, key)
		inputFileAsString= strings.Replace(inputFileAsString, pattern, value , 1)
		log.Debugf("Rendering Variable: %s", key)
	}
	err = ioutil.WriteFile(filename, []byte(inputFileAsString), 0644)
	utils.CheckForError(err)
}

func ViperReadConfig(configFilePath string) *viper.Viper {
	v := viper.GetViper()
	v.SetConfigName(config.GlobalConfig().ClonrConfigFileName)
	v.AddConfigPath(configFilePath)
	log.Debugf("Config File Location: %s", v.ConfigFileUsed())
	err := v.ReadInConfig()
	log.Debug("Running")
	utils.CheckForError(err)
	log.Debug("Running 2")
	return v
}

func answerQuestion(question string)  string {
	fmt.Println("")
	fmt.Println(question)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}