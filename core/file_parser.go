package core

import (
	"bufio"
	"clonr/config"
	"clonr/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

func RenderFile(filename string, varMap map[string]string) {
	input, err := ioutil.ReadFile(filename)
	utils.CheckForError(err)
	inputFileAsString := string(input)
	for key, value := range varMap {
		pattern := config.GlobalConfig().ClonrPrefix + key + config.GlobalConfig().ClonrSuffix
		log.Debugf("Line: %s, Pattern: %s, value: %s, key: %s", inputFileAsString, pattern, value, key)
		inputFileAsString= strings.Replace(inputFileAsString, pattern, value , 1)
		log.Debugf("Rendering Variable: %s", key)
	}
	err = ioutil.WriteFile(filename, []byte(inputFileAsString), 0644)
	utils.CheckForError(err)
}

func ViperReadConfig(configFilePath string) *viper.Viper {
	v := viper.GetViper()
	v.SetConfigName(config.GlobalConfig().ClonrConfigFileName)
	v.AddConfigPath(configFilePath)
	err := v.ReadInConfig()
	utils.CheckForError(err)
	return v
}

func ProcessFiles(configFilePath string) {
	v := ViperReadConfig(configFilePath)

	paths := v.GetStringMapString("paths")
	for path := range paths {
		key := "paths."+ path
		rawVarMap := v.GetStringMap(key)
		processedVarMap := make(map[string]string)
		for variable, question := range rawVarMap {
			processedVarMap[variable] = answerQuestion(fmt.Sprintf("%s", question))
		}
		// Renders the file below
		log.Infof("Rendering file: %s, with vars: %s", path, processedVarMap)
		RenderFile(configFilePath + path, processedVarMap)
	}
}

func answerQuestion(question string)  string {
	fmt.Println("")
	fmt.Println(question)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}