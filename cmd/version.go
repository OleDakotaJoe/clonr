package cmd

import (
"fmt"
"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Clonr",
	Long:  `All software has versions. This is Clonr's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Clonr Project Templating CLI --- v0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}