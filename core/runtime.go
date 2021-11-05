package core

import (
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	"github.com/robertkrimen/otto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"path/filepath"
	"strings"
)

func RunScriptAndReturnString(script string, settings *types.FileProcessorSettings) (string, error) {
	vm := otto.New()
	addAllFunctionsToVmContext(*vm, settings)
	log.Debugf("Running Script: %s", script)
	val, err := vm.Run(script)
	log.Debugf("Returned value from script:  %s", val)
	utils.ExitIfError(err)
	result, gErr := vm.Get(config.Global().ConditionalReturnVarName)
	utils.ExitIfError(gErr)
	return result.ToString()
}

func RunScriptAndReturnBool(script string, settings *types.FileProcessorSettings) (bool, error) {
	vm := otto.New()
	addAllFunctionsToVmContext(*vm, settings)
	log.Debugf("Running Script: %s", script)
	val, err := vm.Run(script)
	log.Debugf("Returned value from script:  %s", val)
	utils.ExitIfError(err)
	result, gErr := vm.Get(config.Global().ConditionalReturnVarName)
	utils.ExitIfError(gErr)
	return result.ToBoolean()
}

func addAllFunctionsToVmContext(vm otto.Otto, settings *types.FileProcessorSettings) {
	addGetClonrVarToContext(vm, settings)
	addGetClonrBoolToContext(vm, settings)
}

func addGetClonrVarToContext(vm otto.Otto, settings *types.FileProcessorSettings) {
	log.Debugln("Adding getClonrVar() function to javascript runtime")
	err := vm.Set("getClonrVar", func(call otto.FunctionCall) otto.Value {
		log.Debugln("Called getClonrVar()")
		result, _ := vm.ToValue(getClonrVar(&types.RuntimeClonrVarDTO{
			FunctionCall:          call,
			FileProcessorSettings: *settings,
			VM:                    vm,
		}))
		return result
	})
	utils.ExitIfError(err)

}

func addGetClonrBoolToContext(vm otto.Otto, processorSettings *types.FileProcessorSettings) {
	log.Debugln("Adding getClonrBool() function to javascript runtime")
	err := vm.Set("getClonrBool", func(call otto.FunctionCall) otto.Value {
		log.Debugln("Called getClonrBool()")
		result, _ := vm.ToValue(getClonrBool(&types.RuntimeClonrVarDTO{
			FunctionCall:          call,
			FileProcessorSettings: *processorSettings,
			VM:                    vm,
		}))
		return result
	})
	utils.ExitIfError(err)
}

func getClonrVar(dto *types.RuntimeClonrVarDTO) string {
	args, err := dto.Argument(0).ToString()
	utils.ExitIfError(err)
	// TODO: add regex check and error thrown if not match
	args = strings.Replace(args, "]", "", 1)
	argsArray := strings.Split(args, "[")
	location := argsArray[0]
	variable := argsArray[1]
	clonrVarMap := dto.MainTemplateMap
	globalVarMap := dto.GlobalsVarMap
	log.Debugf("clonrVarMap: %s", clonrVarMap)
	log.Debugf("globalVarMap: %s", globalVarMap)
	log.Debugf("Looking into maps for location: %s, variable: %s", location, variable)

	if location == config.Global().GlobalsKeyName {
		result := globalVarMap[variable]
		log.Debugf("Got '%s' when trying to access value for %s", result, variable)
		return result
	} else {
		location, err = filepath.Abs(dto.FileProcessorSettings.ConfigFilePath + "/" + argsArray[0])
		utils.ExitIfError(err)
	}
	varMap := clonrVarMap[location]
	result := varMap[variable]
	log.Debugf("ClonrVar being passed into javascript runtime: %s", result)
	return result
}

func getClonrBool(dto *types.RuntimeClonrVarDTO) bool {
	return cast.ToBool(getClonrVar(dto))
}

func ExtractScriptWithTags(fileAsString string) string {
	beginningIndex := strings.Index(fileAsString, config.Global().ConditionalExprPrefix)
	endingIndex := strings.Index(fileAsString, config.Global().ConditionalExprSuffix)

	if beginningIndex == -1 || endingIndex == -1 {
		if beginningIndex != -1 {
			utils.ExitIfError(utils.ThrowError("Something is wrong. Your file contains a closing script  tag, but not an opening one."))
		}
		if endingIndex != -1 {
			utils.ExitIfError(utils.ThrowError("Something is wrong. Your file contains an opening script tag, but not a closing one."))
		}
		return ""
	}
	script := fileAsString[beginningIndex : endingIndex+len(config.Global().ConditionalExprSuffix)]
	log.Debugf("Extracted Script with tags: %s", script)
	return script
}

func RemoveTagsFromScript(script string) string {
	trimmedScript := strings.Replace(script, config.Global().ConditionalExprPrefix, "", 1)
	trimmedScript = strings.Replace(trimmedScript, config.Global().ConditionalExprSuffix, "", 1)

	log.Debugf("Trimmed Script: %s", trimmedScript)
	return trimmedScript
}
