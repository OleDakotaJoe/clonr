package cmd

import (
	"github.com/oledakotajoe/clonr/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clonr",
	Short: "Clonr is a project templating CLI.",
	Long:  `A Fast and Flexible CLI and templating engine for setting up template projects.`,
}

func Execute() {
	err := rootCmd.Execute()
	utils.CheckForError(err)
}
