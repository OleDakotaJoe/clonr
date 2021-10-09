package cmd

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/core"
	"github.com/oledakotajoe/clonr/utils"
	"github.com/otiai10/copy"
	sshutils "golang.org/x/crypto/ssh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
	"strings"
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

		cloneOptions := git.CloneOptions{
			URL:      source,
			Progress: os.Stdout,
		}
		if strings.Contains(source, "git@") {
			var publicKey *ssh.PublicKeys
			var sshPath string
			if runtime.GOOS == "windows" {
				sshPath =os.Getenv("HOMEDRIVE") + "/" +  os.Getenv("HOMEPATH") + "/.ssh/id_rsa.pub"
			} else {
				sshPath = os.Getenv("HOME") + "/.ssh/id_rsa"
			}
			sshKey, _ := ioutil.ReadFile(sshPath)

			_, sshErr := sshutils.ParseRawPrivateKey(sshKey)
			sshPass := ""
			if sshErr != nil {
				sshPass = utils.InputPrompt("What is your ssh-key password?")
			}

			publicKey, keyError := ssh.NewPublicKeys("git", sshKey, sshPass)
			utils.CheckForError(keyError)
			cloneOptions.Auth = publicKey
		}

		log.Info("Clonr is cloning...")
		_, cloneErr := git.PlainClone(projectName, false, &cloneOptions)
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

	sourceUrl := args[0]

	if strings.Contains(sourceUrl, "git@") {
		return sourceUrl, err
	} else {
		_, err = url.ParseRequestURI(sourceUrl)
	}
	return sourceUrl, err
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
