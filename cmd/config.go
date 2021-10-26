package cmd

import (
	"fmt"
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	"github.com/oledakotajoe/clonr/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"os"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Change clonr's configuration",
	Long: `
The clonr config <sub_command> <...args> command has multiple sub_commands.

1. 'clonr config show': displays the current values for clonr's configuration
2. 'clonr config set <property> <value>': sets the property to the value you specify. Beware, some of these can be destructive
      - use 'clonr config set' to get a multiple choice list
      - use 'clonr config set <property>' and you will be prompted for the value
      - use 'clonr config set <property> <value>' and if the property you chose exists, it will be set to the value you specified.
3. 'clonr config reset': resets the configuration back to default settings
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(configCmd)

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
		config.ForEachConfigField(&types.ConfigFieldMutator{ConfigMutator: showProperties, Callback: func(mutator *types.ConfigFieldMutator) { /* do nothing */ }})
	},
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Make an adjustment to clonr's configuration",
	Long: `
	- Use 'clonr config set' to get a multiple choice list
	- Use 'clonr config set <property>' and you will be prompted for the value
	- Use 'clonr config set <property> <value>' and if the property you chose exists, it will be set to the value you specified.
`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			config.ForEachConfigField(&types.ConfigFieldMutator{ConfigMutator: generateConfigList, Callback: useMultipleChoiceToSetProp})
			break
		case 1:
			property := args[0]
			setValueForPropertyIfExists(property)
			break
		case 2:
			property := args[0]
			value := args[1]
			getConfirmationAndSaveProperty(property, value)
			break
		default:

		}
	},
}
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets the configuration back to default settings.",
	Long:  `Resets the configuration back to default settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.GetConfirmationOrExit("Are you sure you what to reset to default settings? This cannot be undone.")
		config.ResetGlobalToDefault()
	},
}

func showProperties(mutator *types.ConfigFieldMutator) {
	property := mutator.Property
	value := mutator.Value
	utils.PrintTabFormattedText(property, value, 28, 8, 4)
}

func generateConfigList(mutator *types.ConfigFieldMutator) {
	property := mutator.Property
	mutator.Result = append(cast.ToSlice(mutator.Result), property)
}

func useMultipleChoiceToSetProp(mutator *types.ConfigFieldMutator) {
	configArr := cast.ToStringSlice(mutator.Result)
	prompt := "Which property do you want to configure?"
	property := utils.MultipleChoiceInputReader(prompt, configArr)

	setValueForProperty(property)
}

func setValueForPropertyIfExists(property string) {
	mutator := types.ConfigFieldMutator{ConfigMutator: generateConfigList, Callback: func(mutator *types.ConfigFieldMutator) { /* do nothing */ }}
	config.ForEachConfigField(&mutator)
	result := cast.ToStringSlice(mutator.Result)

	for _, prop := range result {
		if prop == property {
			setValueForProperty(property)
			break
		}
	}
}

func setValueForProperty(property string) {
	newValue := utils.StringInputReader(fmt.Sprintf("What do you want the value of '%s' to be?", property))
	getConfirmationAndSaveProperty(property, newValue)
}

func getConfirmationAndSaveProperty(property string, value string) {
	utils.GetConfirmationOrExit(fmt.Sprintf("Are you sure you want to set %s to %s? (y/n)", property, value))
	log.Infof("Saving Property: %s as %s", property, value)
	config.SetPropertyAndSave(property, value)
}
