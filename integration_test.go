package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/docker"
)

const (
	imageTag = "epiphanyplatform/azbi:0.0.1"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name               string
		initParams         map[string]string
		wantOutput         string
		wantConfigLocation string
		wantConfigContent  string
	}{
		{
			name:       "default init",
			initParams: nil,
			wantOutput: `#AzBI | setup | ensure required directories
#AzBI | ensure-state-file | checks if state file exists
#AzBI | template-config-file | will template config file (and backup previous if exists)
#AzBI | initialize-state-file | will initialize state file
#AzBI | display-config-file | config file content is:
kind: azbi-config
azbi:
  size: 3
  use_public_ip: true
  location: "northeurope"
  name: "epiphany"
  address_space: ["10.0.0.0/16"]
  address_prefixes: ["10.0.1.0/24"]
  rsa_pub_path: "/shared/vms_rsa.pub"`,
			wantConfigLocation: "azbi/azbi-config.yml",
			wantConfigContent: `kind: azbi-config
azbi:
  size: 3
  use_public_ip: true
  location: "northeurope"
  name: "epiphany"
  address_space: ["10.0.0.0/16"]
  address_prefixes: ["10.0.1.0/24"]
  rsa_pub_path: "/shared/vms_rsa.pub"
`,
		},
		{
			name: "init 2 machines no public ips and named rg",
			initParams: map[string]string{
				"M_VMS_COUNT":  "2",
				"M_PUBLIC_IPS": "false",
				"M_NAME":       "azbi-module-tests",
				"M_VMS_RSA":    "test_vms_rsa"},
			wantOutput: `#AzBI | setup | ensure required directories
#AzBI | ensure-state-file | checks if state file exists
#AzBI | template-config-file | will template config file (and backup previous if exists)
#AzBI | initialize-state-file | will initialize state file
#AzBI | display-config-file | config file content is:
kind: azbi-config
azbi:
  size: 2
  use_public_ip: false
  location: "northeurope"
  name: "azbi-module-tests"
  address_space: ["10.0.0.0/16"]
  address_prefixes: ["10.0.1.0/24"]
  rsa_pub_path: "/shared/test_vms_rsa.pub"`,
			wantConfigLocation: "azbi/azbi-config.yml",
			wantConfigContent: `kind: azbi-config
azbi:
  size: 2
  use_public_ip: false
  location: "northeurope"
  name: "azbi-module-tests"
  address_space: ["10.0.0.0/16"]
  address_prefixes: ["10.0.1.0/24"]
  rsa_pub_path: "/shared/test_vms_rsa.pub"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaName := "vms_rsa"
			if v, ok := tt.initParams["M_VMS_RSA"]; ok {
				rsaName = v
			}
			name := "epiphany"
			if v, ok := tt.initParams["M_NAME"]; ok {
				name = v
			}
			sharedPath, _, _, subscriptionId, _, err := setup(rsaName, name)
			if err != nil {
				t.Fatalf("setup() failed with: %v", err)
			}
			defer cleanup(sharedPath, subscriptionId, name)

			initCommand := []string{"init"}
			for k, v := range tt.initParams {
				initCommand = append(initCommand, fmt.Sprintf("%s=%s", k, v))
			}

			initOpts := &docker.RunOptions{
				Command: initCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			gotOutput := docker.Run(t, imageTag, initOpts)
			if diff := deep.Equal(gotOutput, tt.wantOutput); diff != nil {
				t.Error(diff)
			}

			expectedPath := path.Join(sharedPath, tt.wantConfigLocation)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Fatalf("missing expected file: %s", expectedPath)
			}
			gotFileContent, err := ioutil.ReadFile(expectedPath)
			if err != nil {
				t.Errorf("wasnt able to read form output file: %v", err)
			}
			if diff := deep.Equal(string(gotFileContent), tt.wantConfigContent); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestPlan(t *testing.T) {
	tests := []struct {
		name                           string
		initParams                     map[string]string
		wantPlanOutputLastLine         string
		wantStateLocation              string
		wantStateContent               string
		wantTerraformStateFileLocation string
	}{
		{
			name: "plan 2 machines no public ips and named rg",
			initParams: map[string]string{
				"M_VMS_COUNT":  "2",
				"M_PUBLIC_IPS": "false",
				"M_NAME":       "azbi-module-tests",
				"M_VMS_RSA":    "test_vms_rsa"},
			wantPlanOutputLastLine: "Plan: 7 to add, 0 to change, 0 to destroy.",
			wantStateLocation:      "state.yml",
			wantStateContent: `kind: state
azbi:
  status: initialized
`,
			wantTerraformStateFileLocation: "azbi/terraform-apply.tfplan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaName := "vms_rsa"
			if v, ok := tt.initParams["M_VMS_RSA"]; ok {
				rsaName = v
			}
			name := "epiphany-rg"
			if v, ok := tt.initParams["M_NAME"]; ok {
				name = v
			}
			sharedPath, clientId, clientSecret, subscriptionId, tenantId, err := setup(rsaName, name)
			if err != nil {
				t.Fatalf("setup() failed with: %v", err)
			}
			defer cleanup(sharedPath, subscriptionId, name)

			initCommand := []string{"init"}
			for k, v := range tt.initParams {
				initCommand = append(initCommand, fmt.Sprintf("%s=%s", k, v))
			}
			initOpts := &docker.RunOptions{
				Command: initCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			//in case of error Run function calls FailNow anyways
			docker.Run(t, imageTag, initOpts)

			planCommand := []string{"plan",
				fmt.Sprintf("M_ARM_CLIENT_ID=%s", clientId),
				fmt.Sprintf("M_ARM_CLIENT_SECRET=%s", clientSecret),
				fmt.Sprintf("M_ARM_SUBSCRIPTION_ID=%s", subscriptionId),
				fmt.Sprintf("M_ARM_TENANT_ID=%s", tenantId),
			}

			planOpts := &docker.RunOptions{
				Command: planCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			gotPlanOutput := docker.Run(t, imageTag, planOpts)
			gotPlanOutputLastLine, err := getLastLineFromMultilineSting(gotPlanOutput)
			if err != nil {
				t.Fatalf("reading last line from multiline failed with: %v", err)
			}

			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			expectedPath := path.Join(sharedPath, tt.wantStateLocation)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Fatalf("missing expected file: %s", expectedPath)
			}
			gotFileContent, err := ioutil.ReadFile(expectedPath)
			if err != nil {
				t.Errorf("wasnt able to read form output file: %v", err)
			}
			if diff := deep.Equal(string(gotFileContent), tt.wantStateContent); diff != nil {
				t.Error(diff)
			}

			expectedTerraformStatePath := path.Join(sharedPath, tt.wantTerraformStateFileLocation)
			if _, err := os.Stat(expectedTerraformStatePath); os.IsNotExist(err) {
				t.Fatalf("missing expected file: %s", expectedTerraformStatePath)
			}
		})
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name                    string
		initParams              map[string]string
		wantPlanOutputLastLine  string
		wantApplyOutputLastLine string
	}{
		{
			name: "apply 2 machines no public ips and named rg",
			initParams: map[string]string{
				"M_VMS_COUNT":  "2",
				"M_PUBLIC_IPS": "false",
				"M_NAME":       "azbi-module-tests",
				"M_VMS_RSA":    "test_vms_rsa"},
			wantPlanOutputLastLine:  "Plan: 7 to add, 0 to change, 0 to destroy.",
			wantApplyOutputLastLine: "#AzBI | terraform-output | will prepare terraform output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaName := "vms_rsa"
			if v, ok := tt.initParams["M_VMS_RSA"]; ok {
				rsaName = v
			}
			name := "epiphany-rg"
			if v, ok := tt.initParams["M_NAME"]; ok {
				name = v
			}
			sharedPath, clientId, clientSecret, subscriptionId, tenantId, err := setup(rsaName, name)
			if err != nil {
				t.Fatalf("setup() failed with: %v", err)
			}
			defer cleanup(sharedPath, subscriptionId, name)

			initCommand := []string{"init"}
			for k, v := range tt.initParams {
				initCommand = append(initCommand, fmt.Sprintf("%s=%s", k, v))
			}
			initOpts := &docker.RunOptions{
				Command: initCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			//in case of error Run function calls FailNow anyways
			docker.Run(t, imageTag, initOpts)

			planCommand := []string{"plan",
				fmt.Sprintf("M_ARM_CLIENT_ID=%s", clientId),
				fmt.Sprintf("M_ARM_CLIENT_SECRET=%s", clientSecret),
				fmt.Sprintf("M_ARM_SUBSCRIPTION_ID=%s", subscriptionId),
				fmt.Sprintf("M_ARM_TENANT_ID=%s", tenantId),
			}

			planOpts := &docker.RunOptions{
				Command: planCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			gotPlanOutput := docker.Run(t, imageTag, planOpts)
			gotPlanOutputLastLine, err := getLastLineFromMultilineSting(gotPlanOutput)
			if err != nil {
				t.Fatalf("reading last line from multiline failed with: %v", err)
			}

			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			applyCommand := []string{"apply",
				fmt.Sprintf("M_ARM_CLIENT_ID=%s", clientId),
				fmt.Sprintf("M_ARM_CLIENT_SECRET=%s", clientSecret),
				fmt.Sprintf("M_ARM_SUBSCRIPTION_ID=%s", subscriptionId),
				fmt.Sprintf("M_ARM_TENANT_ID=%s", tenantId),
			}

			applyOpts := &docker.RunOptions{
				Command: applyCommand,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			gotApplyOutput := docker.Run(t, imageTag, applyOpts)
			gotApplyOutputLastLine, err := getLastLineFromMultilineSting(gotApplyOutput)
			if err != nil {
				t.Fatalf("reading last line from multiline failed with: %v", err)
			}

			if diff := deep.Equal(gotApplyOutputLastLine, tt.wantApplyOutputLastLine); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func setup(rsaName string, rgName string) (sharedPath string, armClientId string, armClientSecret string, armSubscriptionId string, armTenantId string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	sharedPath = path.Join(wd, "tests", "shared")
	err = os.MkdirAll(sharedPath, os.ModePerm)
	if err != nil {
		return
	}
	err = generateRsaKeyPair(sharedPath, rsaName)
	if err != nil {
		return
	}
	armClientId, armClientSecret, armSubscriptionId, armTenantId, err = loadEnvironments()
	if err != nil {
		return
	}
	if isResourceGroupPresent(armSubscriptionId, rgName) {
		err = removeResourceGroup(armSubscriptionId, rgName)
		if err != nil {
			return
		}
	}
	return
}

func removeResourceGroup(subscriptionId string, name string) error {
	ctx := context.TODO()
	groupsClient := resources.NewGroupsClient(subscriptionId)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Println("error")
	}
	groupsClient.Authorizer = authorizer
	rgName := fmt.Sprintf("%s-rg", name)
	gdf, err := groupsClient.Delete(ctx, rgName)
	if err != nil {
		return err
	}
	err = gdf.Future.WaitForCompletionRef(ctx, groupsClient.BaseClient.Client)
	if err != nil {
		return err
	}
	return nil
}

func isResourceGroupPresent(subscriptionId string, name string) bool {
	groupsClient := resources.NewGroupsClient(subscriptionId)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Println("error")
	}
	groupsClient.Authorizer = authorizer
	rgName := fmt.Sprintf("%s-rg", name)
	_, err = groupsClient.Get(context.TODO(), rgName)
	if err != nil {
		return false
	} else {
		return true
	}
}

func cleanup(sharedPath string, subscriptionId string, name string) {
	_ = os.RemoveAll(sharedPath)
	_ = removeResourceGroup(subscriptionId, name)
}

func generateRsaKeyPair(directory string, name string) error {
	privateRsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRsaKey)}
	privateKeyBytes := pem.EncodeToMemory(pemBlock)

	publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
	if err != nil {
		return err
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	err = ioutil.WriteFile(path.Join(directory, name), privateKeyBytes, 0600)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(directory, fmt.Sprintf("%s.pub", name)), publicKeyBytes, 0644)
}

func getLastLineFromMultilineSting(s string) (string, error) {
	in := strings.NewReader(s)
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF {
			return string(line), nil
		}
	}
}

func loadEnvironments() (armClientId string, armClientSecret string, armSubscriptionId string, armTenantId string, err error) {
	armClientId = os.Getenv("AZURE_CLIENT_ID")
	if len(armClientId) == 0 {
		err = fmt.Errorf("expected AZURE_CLIENT_ID environment variable")
		return
	}
	armClientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	if len(armClientSecret) == 0 {
		err = fmt.Errorf("expected AZURE_CLIENT_SECRET environment variable")
		return
	}
	armSubscriptionId = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(armSubscriptionId) == 0 {
		err = fmt.Errorf("expected AZURE_SUBSCRIPTION_ID environment variable")
		return
	}
	armTenantId = os.Getenv("AZURE_TENANT_ID")
	if len(armTenantId) == 0 {
		err = fmt.Errorf("expected AZURE_TENANT_ID environment variable")
		return
	}
	return
}
