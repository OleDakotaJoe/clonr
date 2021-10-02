package cmd

import (
	"clonr/utils"
	"fmt"
	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"strings"
)

var nameFlag string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a git project and initializes template engine",
	Long:  `All software has versions. This is Clonr'`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Initializing Clonr Project Templating CLI --- v0.0.1")
		var source = validateAndExtractUrl(args)
		destination :=  determineOutputDir(nameFlag, args)
		cloneProject(source, destination)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&nameFlag, "name", "", "The git URL to read from")
}


func validateAndExtractUrl(args []string) string {
	log.Info("Validating source URL")
	if len(args) == 0 {
		utils.ThrowError("SyntaxError: Must provide git URL" + strings.Join(args[2:], " "), 1)
	}

	_ ,err := url.ParseRequestURI(args[0])
	utils.CheckForError(err)
	log.Info("The source URL you provided is valid.")

	return args[0]
}

func cloneProject(source string, outputDir string) {
	fmt.Println("Cloning git repo... Please Wait")

	_, err := git.PlainClone(outputDir, false, &git.CloneOptions{
		URL: source,
		Progress: os.Stdout,
	} )
	utils.CheckForError(err)
}

func determineOutputDir(nameFlag string, args []string) string {
	var nameArg string
	var result string


	if len(args) == 1 {
		nameArg = args[1]
	} else if len(args) > 1 {
		utils.ThrowError("SyntaxError: Unexpected arguments" + strings.Join(args[2:], " "), 1)
	}


	if nameFlag != "" {
		result =  nameFlag
	} else {
		if &nameArg != nil {
			result = nameArg
		} else {
			result = "my-clonr-app"
		}
	}

	log.Infof("Name of project: %s", result )
	return result
}