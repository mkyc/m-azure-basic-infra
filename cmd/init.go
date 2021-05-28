package cmd

import (
	"fmt"
	"path/filepath"
	"reflect"

	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/load"
	"github.com/epiphany-platform/e-structures/utils/save"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name       string
	vmsRsaPath string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes module configuration file",
	Long:  `Initializes module configuration file (in ` + filepath.Join(defaultSharedDirectory, moduleShortName, configFileName) + `/). `,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("BindPFlags failed")
		}

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
			logger.Fatal().Err(err).Msg("ensureDirectory failed")
		}
		logger.Debug().Msg("load state file")
		state, err := load.State(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("loadState failed")
		}
		logger.Debug().Msg("load config file")
		config, err := load.AzBIConfig(configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("loadConfig failed")
		}

		if state.GetAzBIState() == nil {
			state.AzBI = &st.AzBIState{}
		}

		if state.GetAzBIState() != nil && !reflect.DeepEqual(state.GetAzBIState(), &st.AzBIState{}) {
			if state.AzBI.Status != st.Initialized && state.AzBI.Status != st.Destroyed {
				logger.Fatal().Err(fmt.Errorf("unexpected state: %v", state.AzBI.Status)).Msg("incorrect state")
			}
		}

		logger.Debug().Msg("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("backupFile failed")
		}
		logger.Debug().Msg("backup config file")
		err = backupFile(configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("backupFile failed")
		}

		config.Params.Name = to.StrPtr(name)
		config.Params.RsaPublicKeyPath = to.StrPtr(filepath.Join(SharedDirectory, fmt.Sprintf("%s.pub", vmsRsaPath)))

		state.AzBI.Status = st.Initialized

		logger.Debug().Msg("save config")
		err = save.AzBIConfig(configFilePath, config)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveConfig failed")
		}
		logger.Debug().Msg("save state")
		err = save.State(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveState failed")
		}

		bytes, err := config.Marshal()
		if err != nil {
			logger.Fatal().Err(err).Msg("config.Marshal failed")
		}
		logger.Debug().Msg(string(bytes))
		fmt.Println("Initialized config: \n" + string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().String("name", "epiphany", "prefix given to all resources created") //TODO rename to prefix
	initCmd.Flags().String("vms_rsa", "vms_rsa", "name of rsa keypair to be provided to machines")
}
