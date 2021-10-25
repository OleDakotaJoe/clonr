package cmd

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"strings"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Adds an alias for a git url to be used with the clone command's -alias flag",
	Long:  `All software has versions. This is Clonr's`,
	Run: func(cmd *cobra.Command, args []string) {
		aliasCmdArgs.Args = args
		processAlias(&aliasCmdArgs)
	},
}
var aliasCmdArgs types.AliasCmdArgs

func init() {
	RootCmd.AddCommand(aliasCmd)

	// Set Flags
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.AddFlag, "add", "a", false, "Use this flag to add an alias to the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.UpdateFlag, "update", "u", false, "Use this flag to update an alias already in the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.DeleteFlag, "delete", "d", false, "Use this flag to remove an alias from the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.IsLocalFlag, "local", "l", false, "Use this flag to indicate that an alias points to a local directory.")
	aliasCmd.Flags().StringVarP(&aliasCmdArgs.AliasNameFlag, "name", "n", "", "The name you want to use for your alias.")
	// TODO: add show command
}

func processAlias(args *types.AliasCmdArgs) {
	fmt.Printf("%+v\n", args)
	if !isValidFlags(args) {
		log.Errorln("You must provide exactly one action flag. You must choose only one, -update (-u), -add (-a), or -delete (-d). ")
		return
	}

	for args.AliasNameFlag == "" {
		args.AliasNameFlag = utils.StringInputReader("What do you want the alias name to be?")
	}

	setTemplateLocationForAlias(args)
	aliasManager(args)

}

func isValidFlags(args *types.AliasCmdArgs) bool {
	add := args.AddFlag
	update := args.UpdateFlag
	rm := args.DeleteFlag

	conditions := []bool{add, update, rm}
	trueCount := 0
	for _,condition := range conditions {
		if condition {
			trueCount++
		}
	}

	return trueCount == 1
}

func setTemplateLocationForAlias(args *types.AliasCmdArgs) {
	var templateLocation string
	switch len(args.Args) {
	case 0:
		templateLocation = utils.StringInputReader("What is the git address, or the local path to the template")
		// TODO: add question about is this local
		break
	case 1:
		templateLocation = args.Args[0]
		break
	default:
		log.Errorln("You must provide no more than one argument, the git url or local filepath (if you passed -l).")
		os.Exit(1)
	}
	// TODO: add is this local question
	if args.IsLocalFlag || strings.Contains(templateLocation, "git@") {
		args.AliasLocation = templateLocation
	} else {
		_, err := url.ParseRequestURI(templateLocation)
		utils.ExitIfError(err)
		args.AliasLocation = templateLocation
	}
}

func aliasManager(args *types.AliasCmdArgs) {
	existingAliases := config.Global().Aliases
	resultingAliases := make(map[string]interface{})

	if args.AddFlag {
		if _, ok := existingAliases[args.AliasNameFlag]; ok {
			fmt.Println("This Alias already exists, if you continue you will override it.")
			if !getConfirmation("update", args) {
				return
			}
		}
			resultingAliases =  utils.MergeStringMaps(existingAliases, makeAliasMap(args))
		if !getConfirmation("add", args) {
			return
		}
		log.Infof("Adding alias: %s, %s\n", args.AliasNameFlag, args.AliasLocation)
	} else

	if args.UpdateFlag {
		resultingAliases = existingAliases
		resultingAliases[args.AliasNameFlag] = utils.MergeStringMaps(existingAliases, makeAliasMap(args))
		if !getConfirmation("update", args) {
			return
		}
		log.Infof("Updating alias to: %s, %s\n", args.AliasNameFlag, args.AliasLocation)
	} else

	if args.DeleteFlag {
		resultingAliases = existingAliases
		delete(resultingAliases, args.AliasNameFlag)
		if !getConfirmation("delete", args) {
			return
		}
		log.Infof("Deleting Alias: %s\n", args.AliasNameFlag)
	}
	
	config.SetPropertyAndSave("Aliases", resultingAliases)
}


func makeAliasMap(args *types.AliasCmdArgs) map[string]interface{} {
	return map[string]interface{}{args.AliasNameFlag: map[string]interface{}{ "url": args.AliasLocation, "local": args.IsLocalFlag }}
}

func getConfirmation(action string, args *types.AliasCmdArgs) bool {
	ans := utils.StringInputReader(fmt.Sprintf("Are you sure you want to %s the alias: %s? (y/n)", action, args.AliasNameFlag))
	if strings.ToLower(ans) != "y" {
		log.Infoln("No changes have been made!")
		return false
	}
	return true
}