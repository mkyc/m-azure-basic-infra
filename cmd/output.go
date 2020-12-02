package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// outputCmd represents the output command
var outputCmd = &cobra.Command{
	Use:   "output",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("output called")
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		_, s, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err)
		}
		m, err := getTerraformOutput()
		if err != nil {
			logger.Fatal().Err(err)
		}

		s.AzBI.Output = produceOutput(m)
		err = saveState(stateFilePath, s)
		if err != nil {
			logger.Fatal().Err(err)
		}

		b, err := s.Marshall()
		if err != nil {
			logger.Fatal().Err(err)
		}
		logger.Info().Msg(string(b))
		fmt.Println("Updated output: \n" + string(b))
	},
}

func init() {
	rootCmd.AddCommand(outputCmd)
}
