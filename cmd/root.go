package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	moduleShortName   = "azbi"
	configFileName    = "azbi-config.json"
	stateFileName     = "state.json"
	terraformDir      = "terraform"
	tfVarsFile        = "vars.tfvars.json"
	tfStateFile       = "terraform.tfstate"
	applyTfPlanFile   = "terraform-apply.tfplan"
	destroyTfPlanFile = "terraform-destroy.tfplan"

	defaultSharedDirectory    = "/shared"
	defaultResourcesDirectory = "/resources"
)

var (
	enableDebug bool

	Version string

	SharedDirectory    string
	ResourcesDirectory string

	clientId       string
	clientSecret   string
	subscriptionId string
	tenantId       string

	logger zerolog.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "azbi",
	Long: `AzBI module is responsible for providing basic Azure cloud resources: eg. resource group, virtual network, 
subnets, virtual machines, etc.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PersistentPreRun")

		err := viper.BindPFlags(cmd.PersistentFlags())
		if err != nil {
			logger.Fatal().Err(err)
		}

		SharedDirectory = viper.GetString("shared")
		ResourcesDirectory = viper.GetString("resources")
	},
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal().Err(err)
	}
}

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger = zerolog.New(output).With().Caller().Timestamp().Logger()

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&enableDebug, "debug", "d", false, "print debug information")

	rootCmd.PersistentFlags().String("shared", defaultSharedDirectory, "shared directory location")
	rootCmd.PersistentFlags().String("resources", defaultResourcesDirectory, "resources directory location")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if enableDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
	viper.AutomaticEnv() // read in environment variables that match
}
