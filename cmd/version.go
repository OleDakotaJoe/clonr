package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/version-go/ldflags"
)

var Version = ldflags.Version()

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Clonr",
	Long:  `All software has versions. This is Clonr's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Clonr %s\n", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
