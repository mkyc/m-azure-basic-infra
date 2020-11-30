package cmd

import (
	"fmt"
	state "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
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
		moduleDirectoryPath := filepath.Join(SharedDirectory, moduleShortName)
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		log.Println("ensure directories")
		err := ensureDirectory(moduleDirectoryPath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("load state file")
		s, err := loadState(stateFilePath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("load config file")
		c, err := loadConfig(configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		if s.AzBI.Status != state.Initialized && s.AzBI.Status != state.Destroyed {
			log.Fatal("impossibru state")
		}

		log.Println("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("backup config file")
		err = backupFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		c.Params.VmsCount = to.IntPtr(vmsCount)
		c.Params.UsePublicIP = to.BooPtr(usePublicIPs)
		c.Params.Name = to.StrPtr(name)
		c.Params.RsaPublicKeyPath = to.StrPtr(filepath.Join(SharedDirectory, fmt.Sprintf("%s.pub", vmsRsaPath)))

		s.AzBI.Status = state.Initialized

		log.Println("save config")
		err = saveConfig(configFilePath, c)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("save state")
		err = saveState(stateFilePath, s)
		if err != nil {
			log.Fatal(err)
		}

		b, err := c.Marshall()
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
