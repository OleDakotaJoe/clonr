package cmd

import (
	"clonr/utils"
	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"strings"
)

var nameFlag string

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Initializes a git project and initializes template engine",
	Long:  `This is clonr's primary command. This command will clone a project from a git repository and will `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Initializing clonr project... Please wait")
		var source = validateAndExtractUrl(args)
		destination :=  determineOutputDir(nameFlag, args)
		cloneProject(source, destination)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVar(&nameFlag, "name", "", "The git URL to read from")
}


func validateAndExtractUrl(args []string) string {
	log.Info("Validating source URL")
	if len(args) == 0 {
		utils.ThrowError("SyntaxError: Must provide git URL", 1)
	}

	_ ,err := url.ParseRequestURI(args[0])
	utils.CheckForError(err)
	log.Info("The source URL you provided is valid.")

	return args[0]
}

func cloneProject(source string, outputDir string) {
	log.Info("Cloning git repo... Please Wait")

	_, err := git.PlainClone(outputDir, false, &git.CloneOptions{
		URL: source,
		Progress: os.Stdout,
	} )
	utils.CheckForError(err)
}

func determineOutputDir(nameFlag string, args []string) string {
	var nameArg string
	var result string


	if len(args) == 2 {
		nameArg = args[1]
	} else if len(args) > 2 {
		utils.ThrowError("SyntaxError: Too many arguments. You provided: " + strings.Join(args[1:], " "), 1)
	} else if len(args) < 1 {
		utils.ThrowError("SyntaxError: Too few arguments. ", 1)
	}


	if nameFlag != "" {
		if &nameArg != nil {
			utils.ThrowError("SyntaxError: Unexpected arguments. You provided:  " + strings.Join(args[1:], " "), 1)
		} else {
			result = nameFlag
		}
	} else {
		if nameArg != "" {
			result = nameArg
		} else {
			result = "clonr-app"
		}
	}

	log.Infof("Name of project: %s", result )
	return result
}