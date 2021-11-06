package types

import (
	"github.com/spf13/viper"
	v8 "rogchap.com/v8go"
)

type ClonrVarMap map[string]string
type InputReader func(prompt string) string
type ChoiceReader func(prompt string, choices []string) string
type FileMap map[string]ClonrVarMap

type FileProcessorSettings struct {
	ConfigFilePath            string
	StringInputReader         InputReader
	MultipleChoiceInputReader ChoiceReader
	Viper                     viper.Viper

	MainTemplateMap FileMap
	GlobalsVarMap   ClonrVarMap
	TemplatesVarMap ClonrVarMap
}

type ConfigFieldMutator struct {
	Property      string
	Value         string
	ConfigMutator func(mutator *ConfigFieldMutator)
	Result        interface{}
	Callback      func(mutator *ConfigFieldMutator)
}

type RuntimeDTO struct {
	*v8.FunctionCallbackInfo
	FileProcessorSettings
	*v8.Isolate
}

type ClonrVarDTO struct {
	Args            []string
	MainTemplateMap FileMap
	GlobalsVarMap   ClonrVarMap
	ConfigFilePath  string
}
