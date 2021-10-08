package core

import (
	"github.com/spf13/viper"
)

type ClonrVarMap map[string]string
type InputReader func(prompt string) string
type FileMap map[string]ClonrVarMap

type FileProcessorSettings struct {
	ConfigFilePath    string
	Reader            InputReader
	CloneCmdArguments []string
	Viper             viper.Viper

	mainTemplateMap FileMap
	globalVariables ClonrVarMap
	templateVarMap  ClonrVarMap
}
