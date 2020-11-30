package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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
	PreRun: func(cmd *cobra.Command, args []string) {
		log.Println("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}

		vmsCount = viper.GetInt("vms_count")
		usePublicIPs = viper.GetBool("public_ips")
		name = viper.GetString("name")
		vmsRsaPath = viper.GetString("vms_rsa")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("init called")
		c := backupOrAndInitializeFiles(vmsCount, usePublicIPs, name, vmsRsaPath)
		b, err := c.Save()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().Int("vms_count", 3, "number of virtual machines created by module (default is 3)")
	initCmd.Flags().Bool("public_ips", true, "if created machines should have public IPs attached")
	initCmd.Flags().String("name", "epiphany", "prefix given to all resources created (default is `epiphany`)") //TODO rename to prefix
	initCmd.Flags().String("vms_rsa", "vms_rsa", "name of rsa keypair to be provided to machines (default is `vms_rsa`)")
}
