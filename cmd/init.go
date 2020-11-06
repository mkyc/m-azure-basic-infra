package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		setup()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

}

func setup() {
	sharedDirectory := os.Getenv("M_SHARED")
	if len(sharedDirectory) == 0 {
		log.Fatal(fmt.Errorf("expected M_SHARED environment variable"))
	}

	err := os.MkdirAll(filepath.Join(sharedDirectory, moduleShortName), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
