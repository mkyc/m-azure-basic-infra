package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	moduleShortName   = "azbi"
	configFileName    = "azbi-config.yml"
	stateFileName     = "state.yml"
	terraformDir      = "terraform"
	tfVarsFile        = "vars.tfvars.json"
	tfStateFile       = "terraform.tfstate"
	applyTfPlanFile   = "terraform-apply.tfplan"
	destroyTfPlanFile = "terraform-destroy.tfplan"
)

var (
	cfgFile string

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
	Use:   "m-azure-basic-infrastructure",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PersistentPreRun")

		err := viper.BindPFlags(cmd.PersistentFlags())
		if err != nil {
			logger.Fatal().Err(err)
		}

		SharedDirectory = viper.GetString("shared")
		ResourcesDirectory = viper.GetString("resources")
	},
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.m-azure-basic-infrastructure.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&enableDebug, "debug", "d", false, "enable debug loglevel")

	rootCmd.PersistentFlags().String("shared", "/shared", "Shared directory location (default is `/shared`")
	rootCmd.PersistentFlags().String("resources", "/resources", "Resources directory location (default is `/resources`")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if enableDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
	logger.Debug().Msg("initializing root config")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".m-azure-basic-infrastructure" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".m-azure-basic-infrastructure")
	}

	logger.Debug().Msg("read config variables")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
