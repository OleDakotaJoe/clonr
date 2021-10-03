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
		cloneProject(cloneCmdNameFlag, args)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVar(&cloneCmdNameFlag, "name", config.DefaultConfig().DefaultProjectName, "The git URL to read from")
}


func cloneProject(nameFlag string, args []string) {
	var source = validateAndExtractUrl(args)
	destination, err :=  determineOutputDir(nameFlag, args)
	utils.CheckForError(err)
	log.Info("Cloning git repo... Please Wait")

	_, cloneErr := git.PlainClone(destination, false, &git.CloneOptions{
		URL: source,
		Progress: os.Stdout,
	} )
	utils.CheckForError(cloneErr)
}

func validateAndExtractUrl(args []string) string {
	log.Info("Validating source URL")
	if len(args) == 0 {
		_, err  := utils.ThrowError("SyntaxError: Must provide git URL", 1)
		utils.CheckForError(err)
	}

	_ ,err := url.ParseRequestURI(args[0])
	utils.CheckForError(err)
	log.Info("The source URL you provided is valid.")

	return args[0]
}

func determineOutputDir(outputDirFlag string, args []string) (string, error) {
	var result string
	var err error
	defaultOutputDir := config.DefaultConfig().DefaultProjectName

	if len(args) > 2 {
		_, err  = utils.ThrowError("SyntaxError: Too many arguments. You provided: " + strings.Join(args[1:], " "))
	}

	if len(args) == 1 {
		if outputDirFlag != defaultOutputDir {
			result = outputDirFlag
		} else {
			result = defaultOutputDir
		}
	}

	if len(args) == 2 {
		if outputDirFlag != defaultOutputDir {
			_, err  = utils.ThrowError("SyntaxError: Too many arguments. You provided two name arguments, " + args[1] +
				" and " + outputDirFlag + ". You must only provide one.")
		}

		result = args[1]
	}

	log.Infof("Name of project: %s", result )
	return result, err
}