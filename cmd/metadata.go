package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	useJson bool
)

// metadataCmd represents the metadata command
var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "outputs module metadata information",
	Long: `Outputs module metadata information.

This information is required for inter-module dependency checking. This command is not intended for human users.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("BindPFlags failed")
		}

		useJson = viper.GetBool("json")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("metadata called")
		fmt.Println(printMetadata())
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)

	metadataCmd.Flags().Bool("json", false, "get metadata in json")
}

type Metadata struct {
	Labels map[string]interface{} `yaml:"labels" json:"labels"`
}

func printMetadata() string {
	logger.Debug().Msg("printMetadata")

	labels := Metadata{Labels: map[string]interface{}{
		"version":         Version,
		"name":            "Azure Basic Infrastructure",
		"short":           moduleShortName,
		"kind":            "infrastructure",
		"provider":        "azure",
		"provides-vms":    true,
		"provides-pubips": true,
	}}
	var bytes []byte
	var err error
	if useJson {
		bytes, err = json.Marshal(labels)
	} else {
		bytes, err = yaml.Marshal(labels)
	}
	if err != nil {
		logger.Fatal().Err(err).Msg("Marshal failed")
	}
	return string(bytes)
}
