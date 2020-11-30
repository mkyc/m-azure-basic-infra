package cmd

import (
	"github.com/spf13/viper"
	"log"

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
		log.Println("destroy called")
		c, _ := checkStateAndConfigExistenceAndLoadThem()
		templateTfVars(c)
		terraformDestroy()
		updateStateAfterDestroy()
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.Flags().String("client_id", "", "Azure client identifier")
	destroyCmd.Flags().String("client_secret", "", "Azure client secret")
	destroyCmd.Flags().String("subscription_id", "", "Azure subscription identifier")
	destroyCmd.Flags().String("tenant_id", "", "Azure tenant identifier")
}
