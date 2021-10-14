package core

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io/ioutil"
	"strings"
)

func ProcessFiles(settings *FileProcessorSettings) {
	setGlobalVarMap(settings)
	getFileMapFromConfigFile(settings)
	renderAllFiles(settings)
}

var globalConfig = config.GlobalConfig()

func getFileMapFromConfigFile(settings *FileProcessorSettings) {

	configFilePath := settings.ConfigFilePath
	inputReader := settings.Reader
	v := settings.Viper

	configRootKey := globalConfig.ClonrConfigRootKeyName
	paths := v.GetStringMap(configRootKey)
	log.Debugf("Paths: %s", paths)

	masterVariableMap := make(FileMap)

	for path := range paths {
		log.Infof("Processing path: %s", path)
		pathData := cast.ToStringMap(paths[path])
		fileLocation := configFilePath + cast.ToString(pathData[globalConfig.TemplateFileLocationKeyName])
		variableKey := configRootKey + "." + path + "." + globalConfig.VariablesArrayKeyName
		variables := v.GetStringMap(variableKey)

		log.Debugf("Raw pathData: %s", pathData)
		log.Debugf("Processing file at location: %s", fileLocation)
		log.Debugf("Variables: %s", variables)

		processedVarMap := make(ClonrVarMap)
		for variableName, _ := range variables {
			attributesKey := variableKey + "." + variableName
			question := v.GetStringMapString(attributesKey)["question"]
			defaultAnswer := v.GetStringMapString(attributesKey)["default"]
			if variableName != globalConfig.GlobalVariablesKeyName {
				if defaultAnswer != "" {
					question += fmt.Sprintf(" (%s)", defaultAnswer)
				}
				answer := inputReader(question)
				if answer == "" && defaultAnswer != "" {
					processedVarMap[variableName] = defaultAnswer
				} else {
					processedVarMap[variableName] = answer
				}

			} else {
				processedVarMap[variableName] = "" // just need a placeholder here so that the globals indicator ends up in the master variableName map
			}
			log.Debugf("variableName: %s, question: %s", variableName, question)
		}

		masterVariableMap[fileLocation] = processedVarMap
	}
	settings.mainTemplateMap = masterVariableMap
}

func setGlobalVarMap(settings *FileProcessorSettings) {
	v := settings.Viper
	variablesMapKey := globalConfig.GlobalVariablesKeyName + "." + globalConfig.VariablesArrayKeyName
	unprocessedVarMap := v.GetStringMapString(variablesMapKey)
	globalVarMap := make(ClonrVarMap)

	for key, _ := range unprocessedVarMap {
		questionKey := variablesMapKey + "." + key + "." + globalConfig.QuestionsKeyName
		defaultAnswerKey := variablesMapKey + "." + key + "." + globalConfig.DefaultAnswerKeyName
		question := cast.ToString(v.Get(questionKey))
		defaultAnswer := cast.ToString(v.Get(defaultAnswerKey))
		if defaultAnswer != "" {
			question += fmt.Sprintf(" (%s)", defaultAnswer)
		}

		answer := settings.Reader(question)

		if answer == "" && defaultAnswer != "" {
			globalVarMap[key] = defaultAnswer
		} else {
			globalVarMap[key] = answer
		}
	}

	settings.globalVariables = globalVarMap
}

func renderAllFiles(settings *FileProcessorSettings) {
	fileToVariableMap := settings.mainTemplateMap
	for filepath, variableMap := range fileToVariableMap {
		log.Infof("Rendering file: %s, with vars: %s", filepath, variableMap)
		renderFile(filepath, &variableMap, settings)
	}
}

func renderFile(filepath string, varMap *ClonrVarMap, settings *FileProcessorSettings) {
	input, err := ioutil.ReadFile(filepath)
	utils.CheckForError(err)
	inputFileAsString := string(input)
	for key, value := range *varMap {
		if key == globalConfig.GlobalVariablesKeyName {
			// if globals are provided, this is marked by a "globals" key in the .clonrrc file, loop through globals map to check
			for key, value := range settings.globalVariables {
				globalPattern := globalConfig.ClonrPrefix + globalConfig.GlobalVariablesKeyName + "." + key + globalConfig.ClonrSuffix
				inputFileAsString = strings.Replace(inputFileAsString, globalPattern, value, -1) // -1 makes it replace every occurrence in that file.
				log.Infof("Rendering Variable: %s", key)
			}
		} else {
			clonrPattern := globalConfig.ClonrPrefix + key + globalConfig.ClonrSuffix
			inputFileAsString = strings.Replace(inputFileAsString, clonrPattern, value, -1) // -1 makes it replace every occurrence in that file.

		}
	}
	err = ioutil.WriteFile(filepath, []byte(inputFileAsString), 0644)
	utils.CheckForError(err)
}
