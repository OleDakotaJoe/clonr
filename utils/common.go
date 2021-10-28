package utils

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"
)

func ExitIfError(err error) {
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
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType(configFileType)
	v.AddConfigPath(configFilePath)
	log.Debugf("Config File Location: %s", v.ConfigFileUsed())
	err := v.ReadInConfig()
	return v, err
}

func SaveConfig(v *viper.Viper, location string) {
	err := v.WriteConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err := v.WriteConfigAs(location)
		ExitIfError(err)
	} else if err != nil {
		ExitIfError(err)
	}
}

func StringInputReader(prompt string) string {
	fmt.Println()
	fmt.Println(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// MultipleChoiceInputReader This method takes in a prompt string, and a list of choices, and asks the user for input
// based on the list of choices, the user enters a number to corresponding to their selection.
// The value of the selection is returned.
func MultipleChoiceInputReader(prompt string, choices []string) string {
	fmt.Println()
	for index, choice := range choices {
		PrintTabFormattedText(fmt.Sprintf("[%d]:", index+1), choice, 4, 10, 4)
	}
	fmt.Println()
	fmt.Println(prompt)
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

func GetPassword() string {
	fmt.Println("\nPlease enter your password: ")
	passwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	ExitIfError(err)
	return string(passwd)
}

func PrintTabFormattedText(col1 string, col2 string, minWidth, tabWidth int, padding int) {
	writer := tabwriter.NewWriter(os.Stdout, minWidth, tabWidth, padding, '\t', tabwriter.AlignRight)
	_, pErr := fmt.Fprintf(writer, "%s \t %s \n", col1, col2)
	ExitIfError(pErr)
	wErr := writer.Flush()
	ExitIfError(wErr)
}

func DoesMatchPattern(regex string, input string) bool {
	match, err := regexp.Match(regex, []byte(input))
	ExitIfError(err)
	return match
}

func MergeStringMaps(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func GetConfirmationOrExit(prompt string) {
	prompt += " (y/n)"
	ans := StringInputReader(prompt)
	if strings.ToLower(ans) != "y" {
		log.Infoln("No changes have been made!")
		os.Exit(0)
	}

}

func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + "/" + os.Getenv("HOMEPATH")
	} else {
		return os.Getenv("HOME")
	}
}
