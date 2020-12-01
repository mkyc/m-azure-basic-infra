package cmd

import (
	"log"
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
		log.Println("output called")
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		_, s, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		m, err := getTerraformOutput()
		if err != nil {
			log.Fatal(err)
		}

		s.AzBI.Output = produceOutput(m)
		err = saveState(stateFilePath, s)
		if err != nil {
			log.Fatal(err)
		}

		b, err := s.Marshall()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(b))
	},
}

func init() {
	rootCmd.AddCommand(outputCmd)
}
