package cmd

import (
	"clonr/config"
	"clonr/core"
	"clonr/utils"
	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

/* COMMAND ARGS */

type CloneCmdArguments struct {
	nameFlag    string
	isLocalPath bool
	inputMethod func(string) string
	args        []string
}

var cmdArguments = CloneCmdArguments{
	inputMethod: utils.InputPrompt,
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Initializes a git project and initializes template engine",
	Long:  `This is clonr's primary command. This command will clone a project from a git repository and will `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Initializing clonr project... Please wait")
		cmdArguments.args = args
		cloneProject(&cmdArguments)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVarP(&cmdArguments.nameFlag, "name", "n", config.GlobalConfig().DefaultProjectName, "The git URL to read from")
	cloneCmd.Flags().BoolVarP(&cmdArguments.isLocalPath, "local", "l", false, "Indicates that the path you provide is on your local machine.") //(&cloneCmdLocalFlag, "l", false)
}

func cloneProject(cmdArguments *CloneCmdArguments) {
	args := cmdArguments.args
	pwd, fsErr := os.Getwd()
	utils.CheckForError(fsErr)
	projectName, argErr := determineProjectName(cmdArguments)
	utils.CheckForError(argErr)
	destination := pwd + "/" + projectName

	if cmdArguments.isLocalPath {
		// Source should be the first argument passed in through the CLI
		source := args[0]
		err := copy.Copy(source, destination)
		utils.CheckForError(err)
	} else {
		source, err := validateAndExtractUrl(args)
		utils.CheckForError(err)

		log.Info("Clonr is cloning...")
		_, cloneErr := git.PlainClone(projectName, false, &git.CloneOptions{
			URL:      source,
			Progress: os.Stdout,
		})
		utils.CheckForError(cloneErr)

		log.Debugf("Project root: %s", destination)
	}

	core.ProcessFiles(
		&core.FileProcessorSettings{
			ConfigFilePath:    destination,
			Reader:            cmdArguments.inputMethod,
			CloneCmdArguments: cmdArguments.args,
			Viper:             *utils.ViperReadConfig(destination, config.GlobalConfig().ClonrConfigFileName, config.GlobalConfig().ClonrConfigFileType),
		})
}

func validateAndExtractUrl(args []string) (string, error) {
	log.Info("Validating source URL")
	var err error
	if len(args) == 0 {
		err = utils.ThrowError("SyntaxError: Must provide git URL")
		return "", err
	}

	_, err = url.ParseRequestURI(args[0])
	return args[0], err
}

func determineProjectName(cmdArguments *CloneCmdArguments) (string, error) {
	providedProjectName := cmdArguments.nameFlag
	args := cmdArguments.args

	var result string
	var err error
	defaultProjectName := config.GlobalConfig().DefaultProjectName

	if len(args) > 2 {
		err = utils.ThrowError("SyntaxError: Too many arguments.")
	}

	if len(args) == 1 {
		if providedProjectName != defaultProjectName {
			result = providedProjectName
		} else {
			result = defaultProjectName
		}
	}

	if len(args) == 2 {
		if providedProjectName != defaultProjectName {
			err = utils.ThrowError("SyntaxError: Too many arguments. You provided a flag and an inline argument")
		} else {
			result = args[1]
		}
	}

	log.Infof("Name of project will be: %s", result)

	return result, err
}
