package cmd

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/core"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	sshutils "golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
	"strings"
)

var cloneCmdArgs types.CloneCmdArgs
var cloneProcessorSettings types.FileProcessorSettings

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Initializes a git project and initializes template engine",
	Long:  `This is clonr's primary command. This command will clone a project from a git repository and will `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Initializing clonr project... Please wait")
		cloneCmdArgs.Args = args
		cloneProject(&cloneCmdArgs, &cloneProcessorSettings)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneProcessorSettings.StringInputReader = utils.StringInputReader
	cloneProcessorSettings.MultipleChoiceInputReader = utils.MultipleChoiceInputReader
	cloneCmd.Flags().StringVarP(&cloneCmdArgs.NameFlag, "name", "n", config.Global().DefaultProjectName, "The git URL to read from")
	cloneCmd.Flags().BoolVarP(&cloneCmdArgs.IsLocalPath, "local", "l", false, "Indicates that the path you provide is on your local machine.") //(&cloneCmdLocalFlag, "l", false)
}

func cloneProject(cmdArgs *types.CloneCmdArgs, processorSettings *types.FileProcessorSettings) {
	args := cmdArgs.Args
	pwd, fsErr := os.Getwd()
	utils.ExitIfError(fsErr)
	projectName, argErr := determineProjectName(cmdArgs)
	utils.ExitIfError(argErr)
	destination := pwd + "/" + projectName

	if cmdArgs.IsLocalPath {
		// Source should be the first argument passed in through the CLI
		source := args[0]
		err := copy.Copy(source, destination)
		utils.ExitIfError(err)
	} else {
		source, err := validateAndExtractUrl(args)
		utils.ExitIfError(err)

		cloneOptions := git.CloneOptions{
			URL:      source,
			Progress: os.Stdout,
		}
		if strings.Contains(source, "git@") {
			var publicKey *ssh.PublicKeys
			var sshPath string

			// TODO: add functionality for user to customize location of RSA
			if runtime.GOOS == "windows" {
				sshPath = os.Getenv("HOMEDRIVE") + "/" + os.Getenv("HOMEPATH") + "/.ssh/id_rsa"
			} else {
				sshPath = os.Getenv("HOME") + "/.ssh/id_rsa"
			}
			sshKey, _ := ioutil.ReadFile(sshPath)

			_, sshErr := sshutils.ParseRawPrivateKey(sshKey)
			sshPass := ""

			if sshErr != nil {
				sshPass = utils.GetPassword()
			}

			publicKey, keyError := ssh.NewPublicKeys("git", sshKey, sshPass)
			utils.ExitIfError(keyError)
			cloneOptions.Auth = publicKey
		}

		log.Info("Clonr is cloning...")
		_, cloneErr := git.PlainClone(projectName, false, &cloneOptions)
		utils.ExitIfError(cloneErr)

		log.Debugf("Project root: %s", destination)
	}
	v, err := utils.ViperReadConfig(destination, config.Global().ConfigFileName, config.Global().ConfigFileType)
	utils.ExitIfError(err)
	processorSettings.Viper = *v
	processorSettings.ConfigFilePath = destination

	core.ProcessFiles(processorSettings)
}

func validateAndExtractUrl(args []string) (string, error) {
	log.Info("Validating source URL")
	var err error
	if len(args) == 0 {
		err = utils.ThrowError("SyntaxError: Must provide git URL")
		return "", err
	}

	sourceUrl := args[0]

	if strings.Contains(sourceUrl, "git@") {
		return sourceUrl, err
	} else {
		_, err = url.ParseRequestURI(sourceUrl)
	}
	return sourceUrl, err
}

func determineProjectName(cmdArguments *types.CloneCmdArgs) (string, error) {
	providedProjectName := cmdArguments.NameFlag
	args := cmdArguments.Args

	var result string
	var err error
	defaultProjectName := config.Global().DefaultProjectName

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
