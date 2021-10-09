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
		fmt.Println("v1.0.1")
		fmt.Println("Clonr CLI -- Project Templating Engine")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
