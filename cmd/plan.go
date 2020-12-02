package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

var (
	doDestroy bool
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
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
		doDestroy = viper.GetBool("destroy")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("plan called")
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		c, s, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}

		err = templateTfVars(c)
		if err != nil {
			logger.Fatal().Err(err)
		}
		if !doDestroy {
			err = showModulePlan(c, s)
			if err != nil {
				logger.Fatal().Err(err)
			}
			msg, err := count(terraformPlan())
			if err != nil {
				logger.Fatal().Err(err)
			}
			logger.Info().Msg("Will perform following changes: " + msg)
			fmt.Println("Will perform following changes: \n\t" + msg)
		} else {
			msg, err := count(terraformPlanDestroy())
			if err != nil {
				logger.Fatal().Err(err)
			}
			logger.Info().Msg("Will perform following changes: " + msg)
			fmt.Println("Will perform following changes: \n\t" + msg)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	planCmd.Flags().String("client_id", "", "Azure client identifier")
	planCmd.Flags().String("client_secret", "", "Azure client secret")
	planCmd.Flags().String("subscription_id", "", "Azure subscription identifier")
	planCmd.Flags().String("tenant_id", "", "Azure tenant identifier")
	planCmd.Flags().Bool("destroy", false, "make plan for destroy")
}
