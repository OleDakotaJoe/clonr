package cmd

import (
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var docsCmdArgs types.DocsCmdArgs

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Print the version number of Clonr",
	Long:  `All software has versions. This is Clonr's`,
	Run: func(cmd *cobra.Command, args []string) {
		dirErr := os.MkdirAll(docsCmdArgs.OutputDir, os.ModePerm)
		utils.ExitIfError(dirErr)
		err := doc.GenMarkdownTree(RootCmd, docsCmdArgs.OutputDir)
		utils.ExitIfError(err)
	},
}

func init() {
	RootCmd.AddCommand(docsCmd)
	docsCmd.Flags().StringVarP(&docsCmdArgs.OutputDir,"out","o", "./clonr-docs", "The filepath of the directory you would like to output the files")
}
