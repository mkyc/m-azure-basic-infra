package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	inJson bool
)

// metadataCmd represents the metadata command
var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		//TODO move to debug
		log.Println("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}

		inJson = viper.GetBool("json")
	},
	Run: func(cmd *cobra.Command, args []string) {
		//TODO move to debug
		log.Println("metadata called")
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
	//TODO move to debug
	log.Println("printMetadata")

	l := Metadata{Labels: map[string]interface{}{
		"version":         Version,
		"name":            "Azure Basic Infrastructure",
		"short":           moduleShortName,
		"kind":            "infrastructure",
		"provider":        "azure",
		"provides-vms":    true,
		"provides-pubips": true,
	}}
	var b []byte
	var err error
	if inJson {
		b, err = json.Marshal(l)
	} else {
		b, err = yaml.Marshal(l)
	}
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
