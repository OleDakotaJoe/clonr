package types

import (
	"github.com/spf13/viper"
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
	GlobalVariables ClonrVarMap
	TemplateVarMap  ClonrVarMap
}
