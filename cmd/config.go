package cmd

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/utils"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"text/tabwriter"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Change clonr's configuration",
	Long:  `Configure your clonr setup however your like.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Sub-commands
	configCmd.AddCommand(showCmd)
	configCmd.AddCommand(setCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display clonr's current configuration.",
	Long:  `Display clonr's current configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := *config.Global()
		v := reflect.ValueOf(conf)
		typeConf := v.Type()
		for i := 0; i < v.NumField(); i++ {
			writer := tabwriter.NewWriter(os.Stdout, 28, 8, 1, '\t', tabwriter.AlignRight)
			field := typeConf.Field(i).Name
			value := cast.ToString(v.Field(i))
			_, pErr := fmt.Fprintf(writer, "%s:\t %s \n", field, value)
			utils.ExitIfError(pErr)
			wErr := writer.Flush()
			utils.ExitIfError(wErr)
		}
	},
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Make an adjustment to clonr's configuration",
	Long:  `Display clonr's current configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := config.Global().Viper
		err := v.WriteConfig()
		if os.IsNotExist(err) {
			err := v.WriteConfigAs(utils.GetLocationOfInstalledBinary())
			fmt.Println(utils.GetLocationOfInstalledBinary())
			utils.ExitIfError(err)
		}
	},
}

