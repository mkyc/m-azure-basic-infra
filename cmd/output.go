package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// outputCmd represents the output command
var outputCmd = &cobra.Command{
	Use:   "output",
	Short: "updates module state output subtree",
	Long: `Updates module state output subtree.

This command performs 'terraform output' operation on existing terraform state file and then saves known values to 
module state file output subtree of specific module tree.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("output called")
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		_, state, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("checkAndLoad failed")
		}
		terraformOutputMap, err := getTerraformOutputMap()
		if err != nil {
			logger.Fatal().Err(err).Msg("getTerraformOutputMap failed")
		}

		state.AzBI.Output = produceOutput(terraformOutputMap)
		err = saveState(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveState failed")
		}

		bytes, err := state.Marshal()
		if err != nil {
			logger.Fatal().Err(err).Msg("state.Marshal failed")
		}
		fmt.Println(string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(outputCmd)
}
