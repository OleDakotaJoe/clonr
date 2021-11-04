package core

import (
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	"github.com/robertkrimen/otto"
	"github.com/spf13/cast"
	"strings"
)

func RunScriptAndReturnValue(script string, settings *types.FileProcessorSettings) string {
	vm := otto.New()
	addAllFunctionsToVmContext(*vm, settings)
	val, err := vm.Run(script)
	utils.ExitIfError(err)
	return cast.ToString(val)
}

func addAllFunctionsToVmContext(vm otto.Otto, settings *types.FileProcessorSettings) {
	addGetClonrVarToContext(vm, settings)
	addGetClonrBoolToContext(vm, settings)
}

func addGetClonrVarToContext(vm otto.Otto, settings *types.FileProcessorSettings) {
	err := vm.Set("getClonrBool", func(call otto.FunctionCall) otto.Value {
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
	err := vm.Set("getClonrBool", func(call otto.FunctionCall) otto.Value {
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

	if location == config.Global().GlobalsKeyName {
		for _, v := range globalVarMap {
			if v == variable {
				return v
			}
		}
		utils.ExitIfError(utils.ThrowError("Something went wrong. You must've passed an invalid argument to getClonrBool or getClonrVar funtion"))
	}
	return cast.ToString(cast.ToStringMap(clonrVarMap[location])[variable])

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

	return fileAsString[beginningIndex : endingIndex+len(config.Global().ConditionalExprSuffix)]
}

func RemoveTagsFromScript(script string) string {
	trimmedScript := strings.Replace(script, config.Global().ConditionalExprPrefix, "", 1)
	trimmedScript = strings.Replace(script, config.Global().ConditionalExprSuffix, "", 1)
	return trimmedScript
}
