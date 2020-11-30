package cmd

import (
	"encoding/json"
	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	"io/ioutil"
	"log"
	"path/filepath"

	terra "github.com/mkyc/go-terraform"
)

const (
	moduleShortName   = "azbi" //TODO move to main.consts file
	configFileName    = "azbi-config.yml"
	stateFileName     = "state.yml"
	terraformDir      = "terraform"
	tfVarsFile        = "vars.tfvars.json"
	tfStateFile       = "terraform.tfstate"
	applyTfPlanFile   = "terraform-apply.tfplan"
	destroyTfPlanFile = "terraform-destroy.tfplan"
)

//TODO consider moving those variables nearer functions
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
	destroy        bool

	//if output should be in json
	outputInJson bool
)

//TODO make config a receiver
func marshalConfigParams(config *azbi.Config, tfVarsPath string) error {
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
func templateTfVars(c *azbi.Config) {
	//TODO change to debug log
	log.Println("templateTfVars")
	tfVarsFile := filepath.Join(ResourcesDirectory, terraformDir, tfVarsFile)
	err := marshalConfigParams(c, tfVarsFile)
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

//TODO make State a receiver
func terraformPlanDestroy() {
	log.Println("terraformPlanDestroy")

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
		PlanFilePath:  filepath.Join(SharedDirectory, moduleShortName, destroyTfPlanFile),
		StateFilePath: filepath.Join(SharedDirectory, moduleShortName, tfStateFile),
		NoColor:       true,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = terra.PlanDestroy(options)
	if err != nil {
		log.Fatal(err)
	}
}

//TODO make State a receiver
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

//TODO make State a receiver
func updateStateAfterApply() {
	log.Println("updateStateAfterApply")
	//#AzBI | update-state-after-apply | will update state file after apply
	//@cp $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/azbi-config.tmp.yml
	//@yq d -i $(M_SHARED)/$(M_MODULE_SHORT)/azbi-config.tmp.yml kind
	//@yq m -x -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/azbi-config.tmp.yml
	//@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status applied
	//@rm $(M_SHARED)/$(M_MODULE_SHORT)/azbi-config.tmp.yml
}

//TODO make State a receiver
func terraformOutput() {
	log.Println("terraformOutput")
	options, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: filepath.Join(ResourcesDirectory, terraformDir),
		EnvVars: map[string]string{
			"TF_IN_AUTOMATION": "true",
		},
		StateFilePath: filepath.Join(SharedDirectory, moduleShortName, tfStateFile),
	})
	if err != nil {
		log.Fatal(err)
	}
	m, err := terra.OutputAll(options)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v\n", m)
	//TODO add m to state output
}

//TODO make State a receiver
func terraformDestroy() {
	log.Println("terraformDestroy")

	options, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: filepath.Join(ResourcesDirectory, terraformDir),
		EnvVars: map[string]string{
			"TF_IN_AUTOMATION":      "true",
			"TF_WARN_OUTPUT_ERRORS": "1",
			"ARM_CLIENT_ID":         clientId,
			"ARM_CLIENT_SECRET":     clientSecret,
			"ARM_SUBSCRIPTION_ID":   subscriptionId,
			"ARM_TENANT_ID":         tenantId,
		},
		PlanFilePath:  filepath.Join(SharedDirectory, moduleShortName, destroyTfPlanFile),
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

func updateStateAfterDestroy() {
	log.Println("updateStateAfterDestroy")
	//#AzBI | update-state-after-destroy | will clean state file after destroy
	//@yq d -i $(M_SHARED)/$(M_STATE_FILE_NAME) '$(M_MODULE_SHORT)'
	//@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status destroyed
}
