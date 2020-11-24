package cmd

import (
	"encoding/json"
	"fmt"
	terra "github.com/mkyc/go-terraform"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	moduleShortName = "azbi" //TODO move to main.consts file
	configFileName  = "azbi-config.yml"
	stateFileName   = "state.yml"
	terraformDir    = "terraform"
	tfVarsFile      = "vars.tfvars.json"
	tfStateFile     = "terraform.tfstate"
	applyTfPlanFile = "terraform-apply.tfplan"
)

var (
	cfgFile string

	Version string

	SharedDirectory    string
	ResourcesDirectory string

	//init variables
	vmsCount     int
	usePublicIPs bool
	name         string
	vmsRsaPath   string

	//plan variables
	clientId       string
	clientSecret   string
	subscriptionId string
	tenantId       string
)

type AzBIParams struct {
	VmsCount         int      `yaml:"size" json:"size"`
	UsePublicIP      bool     `yaml:"use_public_ip" json:"use_public_ip"`
	Location         string   `yaml:"location" json:"location"`
	Name             string   `yaml:"name" json:"name"`
	AddressSpace     []string `yaml:"address_space,flow" json:"address_space"`
	AddressPrefixes  []string `yaml:"address_prefixes,flow" json:"address_prefixes"`
	RsaPublicKeyPath string   `yaml:"rsa_pub_path" json:"rsa_pub_path"`
}

type AzBIConfig struct {
	Kind   string     `yaml:"kind"`
	Params AzBIParams `yaml:"azbi"`
}

func printMetadata() string {
	//TODO change to debug log
	log.Println("printMetadata")
	return fmt.Sprintf(`labels:
  version: %s
  name: Azure Basic Infrastructure
  short: %s
  kind: infrastructure
  provider: azure
  provides-vms: true
  provides-pubips: true`, Version, moduleShortName)
}

func ensureSharedDir() {
	//TODO change to debug log
	log.Println("ensureSharedDir")
	err := os.MkdirAll(filepath.Join(SharedDirectory, moduleShortName), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func ensureStateFile() {
	//TODO change to debug log
	log.Println("ensureStateFile")
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
	//TODO change to debug log
	log.Println("initializeConfigFile")
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

//TODO make State a receiver
func initializeStateFile() {
	//TODO change to debug log
	log.Println("initializeStateFile")
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

//TODO Config implement Stringer
func displayCurrentConfigFile() {
	//TODO change to debug log
	log.Println("displayCurrentConfigFile")
	configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
	bytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(bytes))
}

//TODO possibly remove because it's part of shared library
func NewAzBIConfig() *AzBIConfig {
	//TODO change to debug log
	log.Println("NewAzBIConfig")
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

//TODO possibly remove because it's part of shared library
func validateConfig() {
	//TODO change to debug log
	log.Println("validateConfig")
}

//TODO possibly remove because it's part of shared library
func validateState() {
	//TODO change to debug log
	log.Println("validateState")
}

//TODO possibly remove because it's part of shared library
func loadConfig() (*AzBIConfig, error) {
	//TODO change to debug log
	log.Println("loadConfig")
	configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var result AzBIConfig
	err = yaml.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//TODO make config a receiver
func marshalConfigParams(config *AzBIConfig, tfVarsPath string) error {
	//TODO change to debug log
	log.Println("marshalConfigParams")
	params := config.Params
	b, err := json.Marshal(&params)
	if err != nil {
		return err
	}
	log.Println(string(b))
	return ioutil.WriteFile(tfVarsPath, b, 0644)
}

//TODO make Params a receiver
func templateTfVars() {
	//TODO change to debug log
	log.Println("templateTfVars")
	c, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	tfVarsFile := filepath.Join(ResourcesDirectory, terraformDir, tfVarsFile)
	err = marshalConfigParams(c, tfVarsFile)
	if err != nil {
		log.Fatal(err)
	}
}

//TODO make State a receiver
func showModulePlan() {
	log.Println("showModulePlan")
	//#AzBI | module-plan | will perform module plan
	//@yq m -x $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) > $(M_SHARED)/$(M_MODULE_SHORT)/azbi-future-state.tmp
	//@yq w -i $(M_SHARED)/$(M_MODULE_SHORT)/azbi-future-state.tmp kind state
	//@- yq compare $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/azbi-future-state.tmp
	//@rm $(M_SHARED)/$(M_MODULE_SHORT)/azbi-future-state.tmp
}

//TODO make State a receiver
func terraformPlan() {
	log.Println("terraformPlan")

	options, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: filepath.Join(ResourcesDirectory, terraformDir),
		VarFiles:     []string{filepath.Join(ResourcesDirectory, terraformDir, tfVarsFile)},
		EnvVars: map[string]string{
			"TF_IN_AUTOMATION":    "true",
			"ARM_CLIENT_ID":       clientId,
			"ARM_CLIENT_SECRET":   clientSecret,
			"ARM_SUBSCRIPTION_ID": subscriptionId,
			"ARM_TENANT_ID":       tenantId,
		},
		PlanFilePath:  filepath.Join(SharedDirectory, moduleShortName, applyTfPlanFile),
		StateFilePath: filepath.Join(SharedDirectory, moduleShortName, tfStateFile),
		NoColor:       true,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = terra.Plan(options)
	if err != nil {
		log.Fatal(err)
	}
}

func terraformApply() {
	log.Println("terraformApply")

	options, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: filepath.Join(ResourcesDirectory, terraformDir),
		EnvVars: map[string]string{
			"TF_IN_AUTOMATION":    "true",
			"ARM_CLIENT_ID":       clientId,
			"ARM_CLIENT_SECRET":   clientSecret,
			"ARM_SUBSCRIPTION_ID": subscriptionId,
			"ARM_TENANT_ID":       tenantId,
		},
		PlanFilePath:  filepath.Join(SharedDirectory, moduleShortName, applyTfPlanFile),
		StateFilePath: filepath.Join(SharedDirectory, moduleShortName, tfStateFile),
		NoColor:       true,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = terra.Apply(options)
	if err != nil {
		log.Fatal(err)
	}
}
