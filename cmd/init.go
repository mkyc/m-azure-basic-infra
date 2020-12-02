package cmd

import (
	"errors"
	"fmt"
	state "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"
)

var (
	vmsCount     int
	usePublicIPs bool
	name         string
	vmsRsaPath   string
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
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err)
		}

		vmsCount = viper.GetInt("vms_count")
		usePublicIPs = viper.GetBool("public_ips")
		name = viper.GetString("name")
		vmsRsaPath = viper.GetString("vms_rsa")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("init called")
		moduleDirectoryPath := filepath.Join(SharedDirectory, moduleShortName)
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		logger.Debug().Msg("ensure directories")
		err := ensureDirectory(moduleDirectoryPath)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Debug().Msg("load state file")
		s, err := loadState(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Debug().Msg("load config file")
		c, err := loadConfig(configFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}

		if !reflect.DeepEqual(s.AzBI, &state.AzBIState{}) && s.AzBI.Status != state.Initialized && s.AzBI.Status != state.Destroyed {
			logger.Fatal().Err(errors.New(string("unexpected state: " + s.AzBI.Status)))
		}

		logger.Debug().Msg("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Debug().Msg("backup config file")
		err = backupFile(configFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}

		c.Params.VmsCount = to.IntPtr(vmsCount)
		c.Params.UsePublicIP = to.BooPtr(usePublicIPs)
		c.Params.Name = to.StrPtr(name)
		c.Params.RsaPublicKeyPath = to.StrPtr(filepath.Join(SharedDirectory, fmt.Sprintf("%s.pub", vmsRsaPath)))

		s.AzBI.Status = state.Initialized

		logger.Debug().Msg("save config")
		err = saveConfig(configFilePath, c)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Debug().Msg("save state")
		err = saveState(stateFilePath, s)
		if err != nil {
			logger.Fatal().Err(err)
		}

		b, err := c.Marshall()
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Info().Msg(string(b))
		fmt.Println("Initialized config: \n" + string(b))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().Int("vms_count", 3, "number of virtual machines created by module (default is 3)")
	initCmd.Flags().Bool("public_ips", true, "if created machines should have public IPs attached")
	initCmd.Flags().String("name", "epiphany", "prefix given to all resources created (default is `epiphany`)") //TODO rename to prefix
	initCmd.Flags().String("vms_rsa", "vms_rsa", "name of rsa keypair to be provided to machines (default is `vms_rsa`)")
}
