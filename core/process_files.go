package core

import (
	"clonr/config"
	"clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io/ioutil"
	"strings"
)

func ProcessFiles(configFilePath string, inputReader InputReader) {
	masterVarFileMap := getFileMapFromConfigFile(configFilePath, inputReader)
	renderAllFiles(masterVarFileMap)
}

func getFileMapFromConfigFile(configFilePath string, inputReader InputReader) FileMap {
	v := utils.ViperReadConfig(configFilePath, config.GlobalConfig().ClonrConfigFileName, config.GlobalConfig().ClonrConfigFileType)
	configRootKey := config.GlobalConfig().ClonrConfigRootKeyName
	paths := v.GetStringMap(configRootKey)
	log.Debugf("Paths: %s", paths)

	masterVariableMap := make(FileMap)

	for path := range paths {
		log.Infof("Processing path: %s", path)
		pathData := cast.ToStringMap(paths[path])
		fileLocation := configFilePath + cast.ToString(pathData[config.GlobalConfig().TemplateFileLocationKeyName])
		variableKey := configRootKey + "." + path + "." + config.GlobalConfig().VariablesArrayKeyName
		variables := v.GetStringMap(variableKey)

		log.Debugf("Raw pathData: %s", pathData)
		log.Debugf("Processing file at location: %s", fileLocation)
		log.Debugf("Variables: %s", variables)

		processedVarMap := make(ClonrVarMap)
		for variable, _ := range variables {
			questionKey := variableKey + "." + variable
			question := v.GetStringMapString(questionKey)["question"]
			processedVarMap[variable] = inputReader(question)
			log.Debugf("variable: %s, question: %s", variable, question)
		}

		masterVariableMap[fileLocation] = processedVarMap
	}

	return masterVariableMap
}

func renderAllFiles(fileToVariableMap FileMap) {
	for fileLocation, variableMap := range fileToVariableMap {
		log.Infof("Rendering file: %s, with vars: %s", fileLocation, variableMap)
		renderFile(fileLocation, variableMap)
	}
}

func renderFile(filename string, varMap ClonrVarMap) {
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
