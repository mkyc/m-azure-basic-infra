package cmd

import (
	"fmt"
	"path/filepath"

	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/save"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies planned changes on Azure cloud",
	Long: `Applies planned changes on Azure cloud. 

Using plan file created with 'plan' command this command performs actual 'terraform apply' operation. This command
performs following steps: 
 - validates presence of config and module state files
 - checks that module status is either 'Initialized' or 'Destroyed'
 - performs 'terraform apply' operation using existing plan file
 - updates module state file with applied config
 - saves terraform output to module state file. 

This command should always be preceded by 'plan' command.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("BindPFlags failed")
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
			logger.Fatal().Err(err).Msg("checkAndLoad failed")
		}

		if state.GetAzBIState() == nil {
			logger.Fatal().Msg("please run init and plan first")
		}

		// TODO check if there is terraform plan file present
		err = showModulePlan(config, state)
		if err != nil {
			logger.Fatal().Err(err).Msg("showModulePlan failed")
		}
		output, err := terraformApply()
		if err != nil {
			logger.Fatal().Err(err).Msgf("registered following output: \n%s\n", output)
		}

		state.AzBI.Config = config
		state.AzBI.Status = st.Applied

		logger.Debug().Msg("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("backupFile failed")
		}

		terraformOutputMap, err := getTerraformOutputMap()
		if err != nil {
			logger.Fatal().Err(err).Msg("getTerraformOutputMap failed")
		}

		state.AzBI.Output = produceOutput(terraformOutputMap)

		logger.Debug().Msg("save state")
		err = save.State(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveState failed")
		}

		bytes, err := state.Marshal()
		if err != nil {
			logger.Fatal().Err(err).Msg("state.Marshal failed")
		}
		logger.Info().Msg(string(bytes))
		fmt.Println("State after apply: \n" + string(bytes))

		msg, err := count(output)
		if err != nil {
			logger.Fatal().Err(err).Msg("count failed")
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
