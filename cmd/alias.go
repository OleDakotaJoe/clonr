package cmd

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Adds an alias for a git url to be used with the clone command's -alias flag",
	// TODO: add better long descriptions
	Long: `
There are multiple ways to use the alias command.
1. 'clonr alias show': displays a list of available aliases
2. 'clonr alias -a -name=<alias_name> <alias_url>': sets the property to the value you specify. Beware, some of these can be destructive
      - use 'clonr alias -a' to walk through a short alias wizard
      - use 'clonr config set <property>' and you will be prompted for the value
      - use 'clonr config set <property> <value>' and if the property you chose exists, it will be set to the value you specified.
3. 'clonr config reset': resets the configuration back to default settings

NOTE: You can actually pass in the name using a '-name' or '-n' flag, if you prefer.

This would look like this: clonr clone <git_url> -name <name_of_project>
`,
	Run: func(cmd *cobra.Command, args []string) {
		aliasCmdArgs.Args = args
		processAlias(&aliasCmdArgs)
	},
}

var aliasShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays all currently saved aliases",
	// TODO: add better long descriptions
	Long: `some description`,
	Run: func(cmd *cobra.Command, args []string) {
		displayAliases()
	},
}
var aliasCmdArgs types.AliasCmdArgs

func init() {
	RootCmd.AddCommand(aliasCmd)
	aliasCmd.AddCommand(aliasShowCmd)

	// Set Flags
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.AddFlag, "add", "a", false, "Use this flag to add an alias to the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.UpdateFlag, "update", "u", false, "Use this flag to update an alias already in the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.DeleteFlag, "delete", "d", false, "Use this flag to remove an alias from the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.IsLocalFlag, "local", "l", false, "Use this flag to indicate that an alias points to a local directory.")
	aliasCmd.Flags().StringVarP(&aliasCmdArgs.AliasNameFlag, "name", "n", "", "The name you want to use for your alias.")
}

func processAlias(args *types.AliasCmdArgs) {
	if !isValidFlags(args) {
		log.Errorln("You must provide exactly one action flag. You must choose only one, -update (-u), -add (-a), or -delete (-d). ")
		return
	}

	for args.AliasNameFlag == "" {
		var prompt string
		if args.DeleteFlag {
			prompt = "Which alias do you want to delete?"
			displayAliases()
		} else {
			prompt = "What do you want the alias name to be?"
		}
		args.AliasNameFlag = utils.StringInputReader(prompt)
	}

	setTemplateLocationForAlias(args)
	aliasManager(args)

}

func displayAliases() {
	aliases := config.Global().Aliases
	for alias, props := range aliases {
		propsMap := cast.ToStringMapString(props)
		propsKeys := utils.GetKeysFromMap(propsMap)
		for i, prop := range propsKeys {
			value := propsMap[prop]
			propValuePairString := fmt.Sprintf("%s: %s", prop, value)
			if i == 0 {
				fmt.Printf("%s:\n", alias)
			}
			fmt.Printf("\t%s\n", propValuePairString)
		}
	}
}

func isValidFlags(args *types.AliasCmdArgs) bool {
	add := args.AddFlag
	update := args.UpdateFlag
	rm := args.DeleteFlag

	conditions := []bool{add, update, rm}
	trueCount := 0
	for _, condition := range conditions {
		if condition {
			trueCount++
		}
	}

	return trueCount == 1
}

func setTemplateLocationForAlias(args *types.AliasCmdArgs) {
	if args.DeleteFlag {
		return
	}

	var templateLocation string
	switch len(args.Args) {
	case 0:
		templateLocation = utils.StringInputReader("What is the git address, or the local path to the template")
		break
	case 1:
		templateLocation = args.Args[0]
		break
	default:
		log.Errorln("You must provide no more than one argument, the git url or local filepath (if you passed -l).")
		os.Exit(1)
	}
	if !args.IsLocalFlag {
		ans := utils.StringInputReader("Is this a directory on your local machine? (y/n)")
		if strings.ToLower(ans) == "y" {
			args.IsLocalFlag = true
		}
	}

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
			utils.GetConfirmationOrExit(fmt.Sprintf("Are you sure you want to update the alias: %s?", args.AliasNameFlag))
		} else {
			utils.GetConfirmationOrExit(fmt.Sprintf("Are you sure you want to add the alias: %s?", args.AliasNameFlag))
		}
		resultingAliases = utils.MergeStringMaps(existingAliases, makeAliasMap(args))
		log.Infof("Adding alias: %s, %s\n", args.AliasNameFlag, args.AliasLocation)
	} else if args.UpdateFlag {
		resultingAliases = existingAliases
		resultingAliases[args.AliasNameFlag] = utils.MergeStringMaps(existingAliases, makeAliasMap(args))
		utils.GetConfirmationOrExit(fmt.Sprintf("Are you sure you want to update the alias: %s?", args.AliasNameFlag))
		log.Infof("Updating alias to: %s, %s\n", args.AliasNameFlag, args.AliasLocation)
	} else if args.DeleteFlag {
		resultingAliases = existingAliases
		delete(resultingAliases, args.AliasNameFlag)
		utils.GetConfirmationOrExit(fmt.Sprintf("Are you sure you want to delete the alias: %s?", args.AliasNameFlag))
		// Here we are setting Aliases to an empty string to trick viper into removing all aliases.
		// After this, we add all aliases and save config, without the deleted alias
		config.SetPropertyAndSave("Aliases", "")
		log.Infof("Deleting Alias: %s\n", args.AliasNameFlag)
	}

	config.SetPropertyAndSave("Aliases", resultingAliases)
}

func makeAliasMap(args *types.AliasCmdArgs) map[string]interface{} {
	var location string
	if args.IsLocalFlag {
		location, _ = filepath.Abs(args.AliasLocation)
	} else {
		location = args.AliasLocation
	}

	return map[string]interface{}{args.AliasNameFlag: map[string]interface{}{
		config.Global().AliasesUrlKey:            location,
		config.Global().AliasesLocalIndicatorKey: args.IsLocalFlag,
	}}
}
