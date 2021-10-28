package core

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ProcessFiles(settings *types.FileProcessorSettings) {
	processGlobalsVarMap(settings)
	processTemplatesVarMap(settings)
	renderAllFiles(settings)
}

func processTemplatesVarMap(settings *types.FileProcessorSettings) {

	configFilePath := settings.ConfigFilePath
	v := settings.Viper

	configRootKey := config.Global().TemplateRootKeyName
	paths := v.GetStringMap(configRootKey)
	log.Debugf("Paths: %s", paths)

	masterVariableMap := make(types.FileMap)

	for path := range paths {
		log.Debugf("Processing path: %s", path)
		pathData := cast.ToStringMap(paths[path])
		fileLocation, _ := filepath.Abs(configFilePath + cast.ToString(pathData[config.Global().TemplateLocationKeyName]))
		variableKey := configRootKey + "." + path + "." + config.Global().VariablesKeyName
		variables := v.GetStringMap(variableKey)

		log.Debugf("Raw pathData: %s", pathData)
		log.Debugf("Processing file at location: %s", fileLocation)
		log.Debugf("Variables: %s", variables)

		masterVariableMap[fileLocation] = generateVarMap(settings, variableKey)
	}
	settings.MainTemplateMap = masterVariableMap
}

func processGlobalsVarMap(processorSettings *types.FileProcessorSettings) {
	variablesMapKey := config.Global().GlobalsKeyName + "." + config.Global().VariablesKeyName
	processorSettings.GlobalVariables = generateVarMap(processorSettings, variablesMapKey)
}

func generateVarMap(processorSettings *types.FileProcessorSettings, variableKey string) types.ClonrVarMap {
	v := processorSettings.Viper
	unprocessedVarMap := v.GetStringMapString(variableKey)
	processedVarMap := make(types.ClonrVarMap)

	for key, _ := range unprocessedVarMap {
		baseKey := variableKey + "." + key + "."

		questionKey := baseKey + config.Global().QuestionsKeyName
		defaultAnswerKey := baseKey + config.Global().DefaultAnswerKeyName
		choicesKey := baseKey + config.Global().DefaultChoicesKeyName
		validationKey := baseKey + config.Global().ValidationKeyName

		question := cast.ToString(v.Get(questionKey))
		defaultAnswer := cast.ToString(v.Get(defaultAnswerKey))
		choices := v.GetStringSlice(choicesKey)
		validationRegex := v.GetString(validationKey)

		isMultipleChoice := false
		if len(choices) > 0 {
			isMultipleChoice = true
		}

		if key != config.Global().GlobalsKeyName {
			if defaultAnswer != "" {
				question += fmt.Sprintf(" (%s)", defaultAnswer)
			}
			var answer string
			if isMultipleChoice {
				answer = processorSettings.MultipleChoiceInputReader(question, choices)
			} else {
				// Asks for an answer recursively until it gets one that matches the regex (if applicable)
				isValidAnswer := false
				for !isValidAnswer {
					answer = processorSettings.StringInputReader(question)
					if validationRegex != "" && !(defaultAnswer != "" && answer == "") {
						isValidAnswer = utils.DoesMatchPattern(validationRegex, answer)
					} else {
						break
					}
				}
			}

			if answer == "" && defaultAnswer != "" {
				processedVarMap[key] = defaultAnswer
			} else {
				processedVarMap[key] = answer
			}
		} else {
			processedVarMap[key] = "" // just need a placeholder here so that the globals indicator ends up in the master variableName map
		}
	}
	return processedVarMap
}

func renderAllFiles(settings *types.FileProcessorSettings) {
	fileToVariableMap := settings.MainTemplateMap
	for filepath, variableMap := range fileToVariableMap {
		log.Infof("Rendering file: %s, with vars: %s", filepath, variableMap)
		renderFile(filepath, &variableMap, settings)
	}
}

func renderFile(filepath string, varMap *types.ClonrVarMap, settings *types.FileProcessorSettings) {
	input, err := ioutil.ReadFile(filepath)
	utils.ExitIfError(err)
	inputFileAsString := string(input)
	newFilepath := filepath
	for key, value := range *varMap {
		if key == config.Global().GlobalsKeyName {
			// if globals are provided, this is marked by a "globals" key in the .clonr-config.yml file, loop through globals map to check
			for key, value := range settings.GlobalVariables {
				// handle globals here
				globalPlaceholder := config.Global().PlaceholderPrefix + config.Global().GlobalsKeyName + "." + key + config.Global().PlaceholderSuffix
				if strings.Contains(filepath, globalPlaceholder) {
					newFilepath = strings.Replace(filepath, globalPlaceholder, value, -1)
					log.Debugf("Filepath: %s will be renamed to %s", filepath, newFilepath)
				}
				inputFileAsString = strings.Replace(inputFileAsString, globalPlaceholder, value, -1) // -1 makes it replace every occurrence in that file.
				log.Infof("Rendering Variable: %s", key)
			}
		} else {
			// handle
			clonrPlaceholder := config.Global().PlaceholderPrefix + key + config.Global().PlaceholderSuffix
			if strings.Contains(filepath, clonrPlaceholder) {
				newFilepath = strings.Replace(filepath, clonrPlaceholder, value, -1)
				log.Debugf("Filepath: %s will be renamed to %s", filepath, newFilepath)
			}
			inputFileAsString = strings.Replace(inputFileAsString, clonrPlaceholder, value, -1) // -1 makes it replace every occurrence in that file.

		}
	}
	wErr := ioutil.WriteFile(filepath, []byte(inputFileAsString), 0644)
	if filepath != newFilepath {
		rErr := os.Rename(filepath, newFilepath)
		utils.ExitIfError(rErr)
		log.Debugf("File has been renamed from '%s' to '%s'", filepath, newFilepath)
	}
	utils.ExitIfError(wErr)
}
