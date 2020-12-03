package cmd

import (
	"errors"
	"fmt"
	state "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"
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
		c, s, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}

		if !reflect.DeepEqual(s.AzBI, &state.AzBIState{}) && s.AzBI.Status != state.Initialized && s.AzBI.Status != state.Destroyed {
			logger.Fatal().Err(errors.New(string("unexpected state: " + s.AzBI.Status)))
		}

		err = showModulePlan(c, s)
		if err != nil {
			logger.Fatal().Err(err)
		}
		output, err := terraformApply()
		if err != nil {
			logger.Error().Msgf("registered following output: \n%s\n", output)
			logger.Fatal().Err(err)
		}

		s.AzBI.Config = c
		s.AzBI.Status = state.Applied

		logger.Debug().Msg("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Debug().Msg("save state")
		err = saveState(stateFilePath, s)
		if err != nil {
			logger.Fatal().Err(err)
		}

		m, err := getTerraformOutput()
		if err != nil {
			logger.Fatal().Err(err)
		}

		s.AzBI.Output = produceOutput(m)
		err = saveState(stateFilePath, s)
		if err != nil {
			logger.Fatal().Err(err)
		}

		b, err := s.Marshall()
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Info().Msg(string(b))
		fmt.Println("State after apply: \n" + string(b))

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
