package cmd

import (
	"errors"
	"fmt"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"

	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
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
		logger.Debug().Msg("destroy called")
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		config, state, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}

		if !reflect.DeepEqual(state.AzBI, &st.AzBIState{}) && state.AzBI.Status != st.Applied {
			logger.Fatal().Err(errors.New(string("unexpected state: " + state.AzBI.Status)))
		}

		err = templateTfVars(config)
		if err != nil {
			logger.Fatal().Err(err)
		}
		output, err := terraformDestroy()
		if err != nil {
			logger.Error().Msgf("registered following output: \n%s\n", output)
			logger.Fatal().Err(err)
		}
		msg, err := count(output)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Info().Msg("Performed following changes: " + msg)
		fmt.Println("Performed following changes: \n\t" + msg)
		state = updateStateAfterDestroy(state)
		err = saveState(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.Flags().String("client_id", "", "Azure client identifier")
	destroyCmd.Flags().String("client_secret", "", "Azure client secret")
	destroyCmd.Flags().String("subscription_id", "", "Azure subscription identifier")
	destroyCmd.Flags().String("tenant_id", "", "Azure tenant identifier")
}
