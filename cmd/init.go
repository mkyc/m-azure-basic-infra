package cmd

import (
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureSharedDir()
		ensureStateFile()
		initializeConfigFile()
		initializeStateFile()
		displayCurrentConfigFile()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().IntVar(&vmsCount, "vms_count", 3, "number of virtual machines created by module (default is 3)")
	initCmd.Flags().BoolVar(&usePublicIPs, "public_ips", true, "if created machines should have public IPs attached")
	initCmd.Flags().StringVar(&name, "name", "epiphany", "prefix given to all resources created (default is `epiphany`)") //TODO rename to prefix
	initCmd.Flags().StringVar(&vmsRsaPath, "vms_rsa", "vms_rsa", "name of rsa keypair to be provided to machines (default is `vms_rsa`)")
}
