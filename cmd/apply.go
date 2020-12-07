package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"

	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err)
		}

		clientId = viper.GetString("client_id")
		clientSecret = viper.GetString("client_secret")
		subscriptionId = viper.GetString("subscription_id")
		tenantId = viper.GetString("tenant_id")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("apply called")

		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		config, state, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}

		if !reflect.DeepEqual(state.AzBI, &st.AzBIState{}) && state.AzBI.Status != st.Initialized && state.AzBI.Status != st.Destroyed {
			logger.Fatal().Err(errors.New(string("unexpected state: " + state.AzBI.Status)))
		}

		err = showModulePlan(config, state)
		if err != nil {
			logger.Fatal().Err(err)
		}
		output, err := terraformApply()
		if err != nil {
			logger.Error().Msgf("registered following output: \n%s\n", output)
			logger.Fatal().Err(err)
		}

		state.AzBI.Config = config
		state.AzBI.Status = st.Applied

		logger.Debug().Msg("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Debug().Msg("save state")
		err = saveState(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err)
		}

		terraformOutputMap, err := getTerraformOutputMap()
		if err != nil {
			logger.Fatal().Err(err)
		}

		state.AzBI.Output = produceOutput(terraformOutputMap)
		err = saveState(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err)
		}

		bytes, err := state.Marshall()
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Info().Msg(string(bytes))
		fmt.Println("State after apply: \n" + string(bytes))

		msg, err := count(output)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Info().Msg("Performed following changes: " + msg)
		fmt.Println("Performed following changes: \n\t" + msg)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().String("client_id", "", "Azure client identifier")
	applyCmd.Flags().String("client_secret", "", "Azure client secret")
	applyCmd.Flags().String("subscription_id", "", "Azure subscription identifier")
	applyCmd.Flags().String("tenant_id", "", "Azure tenant identifier")
}
