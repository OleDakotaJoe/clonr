package core

import (
	"bufio"
	"clonr/config"
	"clonr/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io/ioutil"
	"os"
	"strings"
)

func ProcessFiles(configFilePath string) {
	v := utils.ViperReadConfig(configFilePath, config.GlobalConfig().ClonrConfigFileName, config.GlobalConfig().ClonrConfigFileType)
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

		// Renders the file below
		log.Infof("Rendering file: %s, with vars: %s", path, processedVarMap)
		renderFile(fileLocation, processedVarMap)
	}
}
func renderFile(filename string, varMap map[string]string) {
	input, err := ioutil.ReadFile(filename)
	utils.CheckForError(err)
	inputFileAsString := string(input)
	for key, value := range varMap {
		pattern := config.GlobalConfig().ClonrPrefix + key + config.GlobalConfig().ClonrSuffix
		log.Debugf("InputFileAsString: %s, Pattern: %s, value: %s, key: %s", inputFileAsString, pattern, value, key)
		inputFileAsString = strings.Replace(inputFileAsString, pattern, value, -1) // -1 makes it replace every occurrence in that file.
		log.Debugf("Rendering Variable: %s", key)
	}
	err = ioutil.WriteFile(filename, []byte(inputFileAsString), 0644)
	utils.CheckForError(err)
}

func answerQuestion(question string) string {
	fmt.Println("")
	fmt.Println(question)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
