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
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sshutils "golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

var cloneCmdArgs types.CloneCmdArgs
var cloneProcessorSettings types.FileProcessorSettings

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Initializes a git project and initializes template engine",
	Long: `
There are multiple ways to use the clone command.
1. 'clonr clone <git_url> <name_of_project>'
    * Clones a remote git repository
    * Replace <git_url> with the url you would use if running git clone <url>
    * <name_of_project> is optional. This will be the name of the directory (inside your working directory) where the project will be cloned to.
    * If you don't provide a name for your project, the name will be clonr-app
2. 'clonr clone -local <local_path> <name_of_project>'
    * Clones a local directory on your filesystem.
    * Notice the '-local' flag. This indicates the local filepath. You can also use '-l' for short
    * Replace <local_path> with either an absolute or relative path to the directory you want to clone

NOTE: You can actually pass in the name using a '-name' flag, if you prefer.

This would look like this: clonr clone <git_url> -name <name_of_project>
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Initializing clonr project... Please wait")
		cloneCmdArgs.Args = args
		cloneProject(&cloneCmdArgs, &cloneProcessorSettings)
	},
}

func init() {
	RootCmd.AddCommand(cloneCmd)
	cloneProcessorSettings.StringInputReader = utils.StringInputReader
	cloneProcessorSettings.MultipleChoiceInputReader = utils.MultipleChoiceInputReader
	cloneCmd.Flags().StringVarP(&cloneCmdArgs.NameFlag, "name", "n", config.Global().DefaultProjectName, "The git URL to read from")
	cloneCmd.Flags().BoolVarP(&cloneCmdArgs.IsLocalPath, "local", "l", false, "Indicates that the path you provide is on your local machine.") //(&cloneCmdLocalFlag, "l", false)
	cloneCmd.Flags().BoolVarP(&cloneCmdArgs.IsAlias, "alias", "a", false, "Indicates that the argument you provided is an alias")              //(&cloneCmdLocalFlag, "l", false)
}

func cloneProject(cmdArgs *types.CloneCmdArgs, processorSettings *types.FileProcessorSettings) {
	pwd, fsErr := os.Getwd()
	utils.ExitIfError(fsErr)
	projectName, argErr := determineProjectName(cmdArgs)
	utils.ExitIfError(argErr)
	destination := pwd + "/" + projectName

	if cmdArgs.IsAlias {
		resolveAlias(cmdArgs)
	}

	if cmdArgs.IsLocalPath {
		// Source should be the first argument passed in through the CLI
		source := cmdArgs.Args[0]
		err := copy.Copy(source, destination)
		utils.ExitIfError(err)
	} else {
		source, err := validateAndExtractUrl(cmdArgs.Args)
		utils.ExitIfError(err)

		cloneOptions := git.CloneOptions{
			URL:      source,
			Progress: os.Stdout,
		}
		if strings.Contains(source, "git@") {
			var publicKey *ssh.PublicKeys
			sshKey, _ := ioutil.ReadFile(config.Global().SSHKeyLocation)

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

	var v *viper.Viper
	v, err := utils.ViperReadConfig(destination, config.Global().ConfigFileName, config.Global().ConfigFileType)
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// This block is here for legacy purposes
		v, err = utils.ViperReadConfig(destination, ".clonrrc", config.Global().ConfigFileType)
		utils.ExitIfError(err)
	}
	processorSettings.Viper = *v
	processorSettings.ConfigFilePath = destination

	core.ProcessFiles(processorSettings)
}

func resolveAlias(args *types.CloneCmdArgs) {
	aliases := config.Global().Aliases
	if len(args.Args) > 1 {
		log.Errorln("You've entered too many arguments.")
		os.Exit(1)
	}
	alias := cast.ToStringMapString(aliases[args.Args[0]])
	if alias == nil {
		log.Errorln("You've provided an invalid alias. Try running 'clonr alias show' to see what is available.")
		os.Exit(1)
	}

	args.Args = []string{cast.ToString(alias[config.Global().AliasesUrlKey])}
	args.IsLocalPath = cast.ToBool(alias[config.Global().AliasesLocalIndicatorKey])
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
