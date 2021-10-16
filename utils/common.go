package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func CheckForError(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func ThrowError(message string) error {
	err := errors.New(message)
	return err
}

func ViperReadConfig(configFilePath string, configFileName string, configFileType string) (*viper.Viper, error) {
	v := viper.GetViper()
	v.SetConfigName(configFileName)
	v.SetConfigType(configFileType)
	v.AddConfigPath(configFilePath)
	log.Debugf("Config File Location: %s", v.ConfigFileUsed())
	err := v.ReadInConfig()
	return v, err
}

func StringInputReader(prompt string) string {
	fmt.Println()
	fmt.Println(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func MultipleChoiceInputReader(prompt string, choices []string) string {
	fmt.Println(prompt)
	fmt.Println()
	for index, choice := range choices {
		fmt.Printf("[%d] : %s\n", index+1, choice)
	}
	fmt.Print("Enter the number of your selection: ")

	var resultIndex int
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	temp := scanner.Text()
	parsedInt, err := strconv.ParseInt(temp, 10, 64)
	if err != nil {
		fmt.Printf("You must provide an Integer. You provided %s", temp)
		fmt.Println()
		return MultipleChoiceInputReader(prompt, choices)
	} else {
		inputValue := int(parsedInt)
		if inputValue > len(choices) || inputValue <= 0 {
			fmt.Println("Your choice was out of range.")
			fmt.Println()
			return MultipleChoiceInputReader(prompt, choices)
		} else {
			resultIndex = int(parsedInt) - 1
		}
	}

	return choices[resultIndex]
}

func GetLocationOfInstalledBinary() string {
	ex, err := os.Executable()
	CheckForError(err)
	return filepath.Dir(ex)
}

func RemoveElementFromSlice(list []string, index int) []string {
	return append(list[:index], list[index+1:]...)
}

func IsVariableValid(variable string) (bool, error) {
	return regexp.Match(config.Global().ClonrVariableRegex, []byte(variable))
}

func GetKeysFromMap(someMap map[string]string) []string {
	keys := make([]string, len(someMap))
	i := 0
	for k := range someMap {
		keys[i] = k
		i++
	}
	return keys
}

func GetPassword() string {
	fmt.Println("\nPlease enter your password: ")
	passwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	CheckForError(err)
	return string(passwd)
}
