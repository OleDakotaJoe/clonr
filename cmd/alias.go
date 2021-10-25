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
	Short: "Adds an alias for a git url to be used with the clone command's '--alias' flag",
	Long: `
There are many ways to use the alias command.
1. 'clonr alias show': displays a list of available aliases
2. 'clonr alias (-a|-u) --name=<alias_name> --url=<template_url> (-l)': sets the alias with the values you specified
3. 'clonr alias (-a|-u) <alias_name> <template_url> (-l)'
4. 'clonr alias (-a|-u) --name=<alias_name> <template_url> (-l)'
5. 'clonr alias (-a|-u) <alias_name> --url=<template_url> (-l)'
6. 'clonr alias (-a|-u)': walks you through a wizard to set the alias
7. 'clonr alias -d ': walks you through a wizard for deleting the alias
8. 'clonr alias -d <alias_name>' deletes the specified alias
9. 'clonr alias -d --name=<alias-name>' deletes the specified alias
`,
	Run: func(cmd *cobra.Command, args []string) {
		aliasCmdArgs.Args = args
		processAlias(&aliasCmdArgs)
	},
}

var aliasShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays all currently saved aliases",
	Long:  `Run 'clonr alias show' to see a list of all available aliases`,
	Run: func(cmd *cobra.Command, args []string) {
		displayAliases()
	},
}
var aliasCmdArgs types.AliasCmdArgs

func init() {
	RootCmd.AddCommand(aliasCmd)
	aliasCmd.AddCommand(aliasShowCmd)
	aliasCmdArgs.StringInputReader = utils.StringInputReader
	aliasCmdArgs.ConfirmFunction = utils.GetConfirmationOrExit
	// Set Flags
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.AddFlag, "add", "a", false, "Use this flag to add an alias to the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.UpdateFlag, "update", "u", false, "Use this flag to update an alias already in the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.DeleteFlag, "delete", "d", false, "Use this flag to remove an alias from the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.IsLocalFlag, "local", "l", false, "Use this flag to indicate that an alias points to a local directory.")
	aliasCmd.Flags().StringVarP(&aliasCmdArgs.AliasNameFlag, "name", "n", "", "The name you want to use for your alias.")
	aliasCmd.Flags().StringVar(&aliasCmdArgs.AliasLocationFlag, "url", "", "The url or local filepath that the alias represents.")
}

func processAlias(args *types.AliasCmdArgs) {
	if !isValidFlags(args) {
		log.Errorln("You must provide exactly one action flag. You must choose only one, -update (-u), -add (-a), or -delete (-d). ")
		os.Exit(1)
	}

	if len(args.Args) > 2 {
		log.Errorln("Too many arguments.")
		os.Exit(1)
	}

	setNameForAlias(args)
	setTemplateLocationForAlias(args)
	aliasManager(args)

}

func setNameForAlias(args *types.AliasCmdArgs) {
	args.ActualAliasName = args.AliasNameFlag
	for args.ActualAliasName == "" {
		var prompt string
		if args.DeleteFlag {
			displayAliases()
			prompt = "Which alias do you want to delete?"
		} else if len(args.Args) == 0 {
			if args.AddFlag {
				prompt = "What do you want the alias name to be?"
			} else if args.UpdateFlag {
				displayAliases()
				prompt = "Which alias do you want to update?"
			}
		} else if len(args.Args) > 0 {
			args.ActualAliasName = args.Args[0]
			break
		}
		args.ActualAliasName = args.StringInputReader(prompt)
		existingAliases := config.Global().Aliases
		if !args.AddFlag {
			if _, ok := existingAliases[args.ActualAliasName]; !ok {
				fmt.Println("That alias does not exist!!!")
				args.ActualAliasName = ""
			}
		}
	}
	log.Infof("Alias: %s", args.ActualAliasName)
}

func displayAliases() {
	aliases := config.Global().Aliases
	for alias, props := range aliases {
		propsMap := cast.ToStringMapString(props)
		i := 0
		for k, v := range propsMap {
			propValuePairString := fmt.Sprintf("%s: %s", k, v)
			if i == 0 {
				fmt.Printf("%s:\n", alias)
				i++
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
	if args.DeleteFlag || args.AliasLocationFlag != "" {
		return
	}

	var templateLocation string
	if args.AliasNameFlag == "" {
		if len(args.Args) < 2 {
			templateLocation = args.StringInputReader("What is the git address, or the local path to the template")
		} else {
			templateLocation = args.Args[1]
		}
	} else {
		if len(args.Args) > 1 {
			log.Errorln("Too many Arguments.")
			os.Exit(1)
		} else if len(args.Args) == 1 {
			templateLocation = args.Args[0]
		} else {
			templateLocation = args.StringInputReader("What is the git address, or the local path to the template")
		}
	}

	if !args.IsLocalFlag {
		ans := args.StringInputReader("Is this a directory on your local machine? (y/n)")
		if strings.ToLower(ans) == "y" {
			args.IsLocalFlag = true
		}
	}

	if args.IsLocalFlag || strings.Contains(templateLocation, "git@") {
		args.ActualAliasLocation = templateLocation
	} else {
		_, err := url.ParseRequestURI(templateLocation)
		utils.ExitIfError(err)
		args.ActualAliasLocation = templateLocation
	}
}

func aliasManager(args *types.AliasCmdArgs) {
	existingAliases := config.Global().Aliases
	resultingAliases := make(map[string]interface{})

	if args.AddFlag {
		if _, ok := existingAliases[args.ActualAliasName]; ok {
			fmt.Println("This Alias already exists, if you continue you will override it.")
			args.ConfirmFunction(fmt.Sprintf("Are you sure you want to update the alias: %s?", args.ActualAliasName))
		} else {
			args.ConfirmFunction(fmt.Sprintf("Are you sure you want to add the alias: %s?", args.ActualAliasName))
		}
		resultingAliases = utils.MergeStringMaps(existingAliases, makeAliasMap(args))
		log.Infof("Adding alias: %s, %s\n", args.ActualAliasName, args.ActualAliasLocation)
	} else if args.UpdateFlag {
		resultingAliases = existingAliases
		resultingAliases[args.ActualAliasName] = utils.MergeStringMaps(existingAliases, makeAliasMap(args))
		args.ConfirmFunction(fmt.Sprintf("Are you sure you want to update the alias: %s?", args.ActualAliasName))
		log.Infof("Updating alias to: %s, %s\n", args.ActualAliasName, args.ActualAliasLocation)
	} else if args.DeleteFlag {
		resultingAliases = existingAliases
		delete(resultingAliases, args.ActualAliasName)
		args.ConfirmFunction(fmt.Sprintf("Are you sure you want to delete the alias: %s?", args.ActualAliasName))
		// Here we are setting Aliases to an empty string to trick viper into removing all aliases.
		// After this, we add all aliases and save config, without the deleted alias
		config.SetPropertyAndSave("Aliases", "")
		log.Infof("Deleting Alias: %s\n", args.ActualAliasName)
	}

	config.SetPropertyAndSave("Aliases", resultingAliases)
}

func makeAliasMap(args *types.AliasCmdArgs) map[string]interface{} {
	var location string
	if args.IsLocalFlag {
		location, _ = filepath.Abs(args.ActualAliasLocation)
	} else {
		location = args.ActualAliasLocation
	}

	return map[string]interface{}{args.ActualAliasName: map[string]interface{}{
		config.Global().AliasesUrlKey:            location,
		config.Global().AliasesLocalIndicatorKey: args.IsLocalFlag,
	}}
}
