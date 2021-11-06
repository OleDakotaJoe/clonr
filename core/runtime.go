package core

import (
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"path/filepath"
	v8 "rogchap.com/v8go"
	"strings"
)

func RunScriptAndReturnString(script string, settings *types.FileProcessorSettings) (string, error) {
	ctx := prepareVMContext(settings)
	_, rErr := ctx.RunScript(script, "script.js")
	utils.ExitIfError(rErr)
	obj := ctx.Global()
	val, err := obj.Get(config.Global().ConditionalReturnVarName)
	utils.ExitIfError(err)
	return val.String(), err
}

func RunScriptAndReturnBool(script string, settings *types.FileProcessorSettings) (bool, error) {
	ctx := prepareVMContext(settings)
	_, rErr := ctx.RunScript(script, "script.js")
	utils.ExitIfError(rErr)
	obj := ctx.Global()
	val, err := obj.Get(config.Global().ConditionalReturnVarName)
	utils.ExitIfError(err)
	return val.Boolean(), err
}

func prepareVMContext(settings *types.FileProcessorSettings) *v8.Context {
	iso, _ := v8.NewIsolate()
	global, _ := v8.NewObjectTemplate(iso)

	addGetClonrVarToContext(iso, settings, global)
	addGetClonrBoolToContext(iso, settings, global)
	ctx, _ := v8.NewContext(iso, global)

	return ctx
}

func addGetClonrVarToContext(iso *v8.Isolate, settings *types.FileProcessorSettings, global *v8.ObjectTemplate) {
	getClonrVar, _ := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		result := getClonrVar(&types.RuntimeDTO{
			FunctionCallbackInfo:  info,
			FileProcessorSettings: *settings,
			Isolate:               iso,
		})

		val, _ := v8.NewValue(iso, result)
		return val
	})

	err := global.Set("getClonrVar", getClonrVar)
	utils.ExitIfError(err)
}

func addGetClonrBoolToContext(iso *v8.Isolate, settings *types.FileProcessorSettings, global *v8.ObjectTemplate) {

	getClonrBool, _ := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		result := getClonrBool(&types.RuntimeDTO{
			FunctionCallbackInfo:  info,
			FileProcessorSettings: *settings,
			Isolate:               iso,
		})
		val, _ := v8.NewValue(iso, result)
		return val
	})

	err := global.Set("getClonrBool", getClonrBool)
	utils.ExitIfError(err)
}

func getClonrVar(dto *types.RuntimeDTO) string {
	var args []string
	counter := 0
	for range dto.Args() {
		arg := dto.Args()[counter].String()
		args = append(args, arg)
		counter++
	}

	result := resolveClonrVariable(&types.ClonrVarDTO{
		Args:            args,
		MainTemplateMap: dto.MainTemplateMap,
		GlobalsVarMap:   dto.GlobalsVarMap,
		ConfigFilePath:  dto.ConfigFilePath,
	})

	return result
}

func getClonrBool(dto *types.RuntimeDTO) bool {
	return cast.ToBool(getClonrVar(dto))
}

func resolveClonrVariable(dto *types.ClonrVarDTO) string {
	// TODO: add regex check and error thrown if not match
	arg := dto.Args[0]
	arg = strings.Replace(arg, "]", "", 1)
	argsArray := strings.Split(arg, "[")
	location := argsArray[0]
	variable := argsArray[1]
	clonrVarMap := dto.MainTemplateMap
	globalVarMap := dto.GlobalsVarMap
	log.Debugf("clonrVarMap: %s", clonrVarMap)
	log.Debugf("globalVarMap: %s", globalVarMap)
	log.Debugf("Looking into maps for location: %s, variable: %s", location, variable)

	if location == config.Global().GlobalsKeyName {
		result := cast.ToString(globalVarMap[variable])
		log.Debugf("Got '%s' when trying to access value for %s", result, variable)
		return result
	} else {
		var err error
		location, err = filepath.Abs(dto.ConfigFilePath + "/" + argsArray[0])
		utils.ExitIfError(err)
	}
	varMap := clonrVarMap[location]
	result := varMap[variable]
	log.Debugf("ClonrVar being passed into javascript runtime: %s", result)
	return result
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
