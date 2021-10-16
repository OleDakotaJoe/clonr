package cmd

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Change clonr's configuration",
	Long:  `Configure your clonr setup however your like.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Global())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
