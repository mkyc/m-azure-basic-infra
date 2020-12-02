package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
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
		log.Println("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}

		clientId = viper.GetString("client_id")
		clientSecret = viper.GetString("client_secret")
		subscriptionId = viper.GetString("subscription_id")
		tenantId = viper.GetString("tenant_id")
		destroy = viper.GetBool("destroy")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("plan called")
		//TODO ensure clientId, clientSecret, subscriptionId, tenantId
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		c, s, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		err = templateTfVars(c)
		if err != nil {
			log.Fatal(err)
		}
		if !destroy {
			err = showModulePlan(c, s)
			if err != nil {
				log.Fatal(err)
			}
			terraformPlan()
		} else {
			terraformPlanDestroy()
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
