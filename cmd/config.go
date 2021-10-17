package cmd

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"strings"
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
	configCmd.AddCommand(resetCmd)
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
			if field != "Viper" {
				_, pErr := fmt.Fprintf(writer, "%s:\t %s \n", field, value)
				utils.ExitIfError(pErr)
			}
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
		propertyName := args[0]
		value := args[1]
		v := config.Global().Viper
		switch propertyName {
		case "DefaultProjectName":
			v.Set("DefaultProjectName", value)
			break
		case "ConfigFileName":
			v.Set("ConfigFileName", value)
			break
		case "ConfigFileType":
			v.Set("ConfigFileType", value)
			break
		case "PlaceholderRegex":
			v.Set("PlaceholderRegex", value)
			break
		case "PlaceholderPrefix":
			v.Set("PlaceholderPrefix", value)
			break
		case "PlaceholderSuffix":
			v.Set("PlaceholderSuffix", value)
			break
		case "TemplateRootKeyName":
			v.Set("TemplateRootKeyName", value)
			break
		case "TemplateLocationKeyName":
			v.Set("TemplateLocationKeyName", value)
			break
		case "VariablesKeyName":
			v.Set("VariablesKeyName", value)
			break
		case "GlobalsKeyName":
			v.Set("GlobalsKeyName", value)
			break
		case "QuestionsKeyName":
			v.Set("QuestionsKeyName", value)
			break
		case "DefaultAnswerKeyName":
			v.Set("DefaultAnswerKeyName", value)
			break
		case "DefaultChoicesKeyName":
			v.Set("DefaultChoicesKeyName", value)
			break
		case "LogLevel":
			v.Set("LogLevel", value)
			break
		default:
			fmt.Println()
			log.Errorf("%s is not a clonr property.", propertyName)
			fmt.Printf("\nRun 'clonr config show' for a list of options. No changes were made\n")
		}
		utils.SaveConfig(v, utils.GetLocationOfInstalledBinary()+config.Global().ConfigFileName)
	},
}
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Display clonr's current configuration.",
	Long:  `Display clonr's current configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		response := utils.StringInputReader("Are you sure you what to do this? (y/n)")
		if strings.ToLower(response) != "y" {
			return
		}
		config.ResetGlobalToDefault()
	},
}
