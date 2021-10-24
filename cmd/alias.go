package cmd

import (
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Adds an alias for a git url to be used with the clone command's -alias flag",
	Long:  `All software has versions. This is Clonr's`,
	Run: func(cmd *cobra.Command, args []string) {
		aliasCmdArgs.Args = args
		processAlias(aliasCmdArgs)
	},
}
var aliasCmdArgs types.AliasCmdArgs

func init() {
	RootCmd.AddCommand(aliasCmd)

	// Set Flags
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.AddFlag, "add", "a", true, "Use this flag to add an alias to the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.UpdateFlag, "update", "u", true, "Use this flag to update an alias already in the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.DeleteFlag, "delete", "d", true, "Use this flag to remove an alias from the list.")
	aliasCmd.Flags().BoolVarP(&aliasCmdArgs.DeleteFlag, "local", "l", true, "Use this flag to indicate that an alias points to a local directory.")
	aliasCmd.Flags().StringVarP(&aliasCmdArgs.AliasNameFlag, "name", "n", "", "The name you want to use for your alias.")
	aliasCmd.Flags().StringVarP(&aliasCmdArgs.AliasLocationFlag, "name", "n", "", "The git url or the path to the file you want to use for your alias.")
}

func processAlias(args types.AliasCmdArgs) {
	config.ForEachConfigField()
}
