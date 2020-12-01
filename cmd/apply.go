package cmd

import (
	state "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/spf13/viper"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
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
		log.Println("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}

		clientId = viper.GetString("client_id")
		clientSecret = viper.GetString("client_secret")
		subscriptionId = viper.GetString("subscription_id")
		tenantId = viper.GetString("tenant_id")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("apply called")

		//TODO check if not applied already

		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		c, s, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		err = showModulePlan(c, s)
		if err != nil {
			log.Fatal(err)
		}
		err = terraformApply()
		if err != nil {
			log.Fatal(err)
		}

		s.AzBI.Config = c
		s.AzBI.Status = state.Applied

		log.Println("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("save state")
		err = saveState(stateFilePath, s)
		if err != nil {
			log.Fatal(err)
		}

		m, err := getTerraformOutput()
		if err != nil {
			log.Fatal(err)
		}

		s.AzBI.Output = produceOutput(m)
		err = saveState(stateFilePath, s)
		if err != nil {
			log.Fatal(err)
		}

		b, err := s.Marshall()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(b))
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().String("client_id", "", "Azure client identifier")
	applyCmd.Flags().String("client_secret", "", "Azure client secret")
	applyCmd.Flags().String("subscription_id", "", "Azure subscription identifier")
	applyCmd.Flags().String("tenant_id", "", "Azure tenant identifier")
}
