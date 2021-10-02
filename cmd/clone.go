package cmd

import (
	"clonr/config"
	"clonr/utils"
	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"strings"
)

/* FLAG VARIABLES */
var cloneCmdNameFlag string

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Initializes a git project and initializes template engine",
	Long:  `This is clonr's primary command. This command will clone a project from a git repository and will `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Initializing clonr project... Please wait")
		cloneProject(config.DefaultConfig().DefaultProjectName, cloneCmdNameFlag, args)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVar(&cloneCmdNameFlag, "name", "", "The git URL to read from")
}


func cloneProject(defaultNameFlag string, nameFlag string, args []string) {
	var source = validateAndExtractUrl(args)
	destination :=  determineOutputDir( defaultNameFlag , nameFlag, args)
	log.Info("Cloning git repo... Please Wait")

	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL: source,
		Progress: os.Stdout,
	} )
	utils.CheckForError(err)
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

func determineOutputDir(defaultOutputDir string, outputDirFlag string, args []string) string {
	var result string

	if len(args) > 2 {
		utils.ThrowError("SyntaxError: Too many arguments. You provided: " + strings.Join(args[1:], " "), 1)
	}

	if len(args) == 1 {
		if outputDirFlag != "" {
			result = outputDirFlag
		} else {
			result = defaultOutputDir
		}
	}

	if len(args) == 2 {
		if outputDirFlag != defaultOutputDir {
			utils.ThrowError("SyntaxError: Too many arguments. You provided two name arguments, " + args[1] +
				" and " +outputDirFlag+ ". You must only provide one."  , 1)
		}

		result = args[1]
	}

	log.Infof("Name of project: %s", result )
	return result
}