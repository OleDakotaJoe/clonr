package core

import (
	"bufio"
	"clonr/config"
	"clonr/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
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
		pathData := cast.ToStringMap(paths[path])
		fileLocation := configFilePath + cast.ToString(pathData[config.GlobalConfig().TemplateFileLocationKeyName])
		variableKey := configRootKey + "." + path + "." + config.GlobalConfig().VariablesArrayKeyName
		variables := v.GetStringMap(variableKey)
		processedVarMap := make(map[string]string)

		log.Debugf("Raw pathData: %s", pathData)
		log.Debugf("Processing file at location: %s", fileLocation)
		log.Debugf("Variables: %s", variables)


		for variable, _ := range variables {
			questionKey := variableKey + "." + variable
			question := v.GetStringMapString(questionKey)["question"]
			processedVarMap[variable] = answerQuestion(question)
			log.Debugf("variable: %s, question: %s", variable, question)
		}

		//// Renders the file below
		log.Infof("Rendering file: %s, with vars: %s", path, processedVarMap)
		RenderFile(fileLocation, processedVarMap)
	}
}
func RenderFile(filename string, varMap map[string]string) {
	input, err := ioutil.ReadFile(filename)
	utils.CheckForError(err)
	inputFileAsString := string(input)
	for key, value := range varMap {
		pattern := config.GlobalConfig().ClonrPrefix + key + config.GlobalConfig().ClonrSuffix
		log.Debugf("InputFileAsString: %s, Pattern: %s, value: %s, key: %s", inputFileAsString, pattern, value, key)
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
	utils.CheckForError(err)
	return v
}

func answerQuestion(question string)  string {
	fmt.Println("")
	fmt.Println(question)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}