package cmd

import (
	"fmt"

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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apply called")
		ensureSharedDir()
		showModulePlan()
		terraformApply()
		updateStateAfterApply()
		terraformOutput()
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringVar(&clientId, "client_id", "", "Azure client identifier")
	applyCmd.Flags().StringVar(&clientSecret, "client_secret", "", "Azure client secret")
	applyCmd.Flags().StringVar(&subscriptionId, "subscription_id", "", "Azure subscription identifier")
	applyCmd.Flags().StringVar(&tenantId, "tenant_id", "", "Azure tenant identifier")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
