package cmd

import (
	"errors"
	"fmt"
	st "github.com/epiphany-platform/e-structures/state/v0"
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
		state, err := loadState(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Debug().Msg("load config file")
		config, err := loadConfig(configFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}

		if !reflect.DeepEqual(state.AzBI, &st.AzBIState{}) && state.AzBI.Status != st.Initialized && state.AzBI.Status != st.Destroyed {
			logger.Fatal().Err(errors.New(string("unexpected state: " + state.AzBI.Status)))
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

		config.Params.VmsCount = to.IntPtr(vmsCount)
		config.Params.UsePublicIP = to.BooPtr(usePublicIPs)
		config.Params.Name = to.StrPtr(name)
		config.Params.RsaPublicKeyPath = to.StrPtr(filepath.Join(SharedDirectory, fmt.Sprintf("%s.pub", vmsRsaPath)))

		state.AzBI.Status = st.Initialized

		logger.Debug().Msg("save config")
		err = saveConfig(configFilePath, config)
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Debug().Msg("save state")
		err = saveState(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err)
		}

		bytes, err := config.Marshall()
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Info().Msg(string(bytes))
		fmt.Println("Initialized config: \n" + string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().Int("vms_count", 3, "number of virtual machines created by module (default is 3)")
	initCmd.Flags().Bool("public_ips", true, "if created machines should have public IPs attached")
	initCmd.Flags().String("name", "epiphany", "prefix given to all resources created (default is `epiphany`)") //TODO rename to prefix
	initCmd.Flags().String("vms_rsa", "vms_rsa", "name of rsa keypair to be provided to machines (default is `vms_rsa`)")
}
