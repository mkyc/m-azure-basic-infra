package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type AzBIParams struct {
	VmsCount         int      `yaml:"size"`
	UsePublicIP      bool     `yaml:"use_public_ip"`
	Location         string   `yaml:"location"`
	Name             string   `yaml:"name"`
	AddressSpace     []string `yaml:"address_space,flow"`
	AddressPrefixes  []string `yaml:"address_prefixes,flow"`
	RsaPublicKeyPath string   `yaml:"rsa_pub_path"`
}

type AzBIConfig struct {
	Kind   string     `yaml:"kind"`
	Params AzBIParams `yaml:"azbi"`
}

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
	Run: func(cmd *cobra.Command, args []string) {
		ensureSharedDir()
		ensureStateFile()
		initializeConfigFile()
		initializeStateFile()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().IntVar(&vmsCount, "vms_count", 3, "number of virtual machines created by module (default is 3)")
	initCmd.Flags().BoolVar(&usePublicIPs, "public_ips", true, "if created machines should have public IPs attached")
	initCmd.Flags().StringVar(&name, "name", "epiphany", "prefix given to all resources created (default is `epiphany`)") //TODO rename to prefix
	initCmd.Flags().StringVar(&vmsRsaPath, "vms_rsa", "vms_rsa", "name of rsa keypair to be provided to machines (default is `vms_rsa`)")
}

func ensureSharedDir() {
	err := os.MkdirAll(filepath.Join(SharedDirectory, moduleShortName), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func ensureStateFile() {
	file, err := os.OpenFile(filepath.Join(SharedDirectory, stateFileName), os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func initializeConfigFile() {
	configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
	backupConfigFilePath := fmt.Sprintf("%s.backup", configFilePath)
	_, err := os.Stat(configFilePath)
	if err == nil || os.IsExist(err) {
		err = os.Rename(configFilePath, backupConfigFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	b, err := yaml.Marshal(NewAzBIConfig())
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = ioutil.WriteFile(configFilePath, b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeStateFile() {
	//TODO join with ensureStateFile()
	stateFilePath := filepath.Join(SharedDirectory, stateFileName)
	backupStateFilePath := fmt.Sprintf("%s.backup", stateFilePath)
	m := make(map[interface{}]interface{})
	m["kind"] = "state"
	m[moduleShortName] = make(map[interface{}]interface{})
	m[moduleShortName].(map[interface{}]interface{})["status"] = "initialized"
	_, err := os.Stat(stateFilePath)
	if err == nil || os.IsExist(err) {
		err = os.Rename(stateFilePath, backupStateFilePath)
		if err != nil {
			log.Fatal(err)
		}
		bytes, err := ioutil.ReadFile(backupStateFilePath)
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal(bytes, &m)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Print(err)
	}
	b, err := yaml.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(stateFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = ioutil.WriteFile(stateFilePath, b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func NewAzBIConfig() *AzBIConfig {
	return &AzBIConfig{
		Kind: fmt.Sprintf("%s-config", moduleShortName),
		Params: AzBIParams{
			VmsCount:         vmsCount,
			UsePublicIP:      usePublicIPs,
			Location:         "northeurope",
			Name:             name,
			AddressSpace:     []string{"10.0.0.0/16"},
			AddressPrefixes:  []string{"10.0.1.0/24"},
			RsaPublicKeyPath: filepath.Join(SharedDirectory, fmt.Sprintf("%s.pub", vmsRsaPath)),
		},
	}
}
