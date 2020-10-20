package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
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
			remoteSharedPath, localSharedPath, environments := setup(t, rsaName, name)
			defer cleanup(t, localSharedPath, environments["M_ARM_SUBSCRIPTION_ID"], name)

			gotOutput := dockerInit(t, tt.initParams, remoteSharedPath)
			if diff := deep.Equal(gotOutput, tt.wantOutput); diff != nil {
				t.Error(diff)
			}

			expectedPath := path.Join(localSharedPath, tt.wantConfigLocation)
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
			remoteSharedPath, localSharedPath, environments := setup(t, rsaName, name)
			defer cleanup(t, localSharedPath, environments["M_ARM_SUBSCRIPTION_ID"], name)

			dockerInit(t, tt.initParams, remoteSharedPath)

			gotPlanOutputLastLine := dockerPlan(t, environments, remoteSharedPath)
			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			expectedPath := path.Join(localSharedPath, tt.wantStateLocation)
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

			expectedTerraformStatePath := path.Join(localSharedPath, tt.wantTerraformStateFileLocation)
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
		{
			name: "apply 2 machines with public ips and named rg",
			initParams: map[string]string{
				"M_VMS_COUNT":  "2",
				"M_PUBLIC_IPS": "true",
				"M_NAME":       "azbi-module-tests",
				"M_VMS_RSA":    "test_vms_rsa"},
			wantPlanOutputLastLine:  "Plan: 12 to add, 0 to change, 0 to destroy.",
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
			remoteSharedPath, localSharedPath, environments := setup(t, rsaName, name)
			defer cleanup(t, localSharedPath, environments["M_ARM_SUBSCRIPTION_ID"], name)

			dockerInit(t, tt.initParams, remoteSharedPath)

			gotPlanOutputLastLine := dockerPlan(t, environments, remoteSharedPath)
			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			gotApplyOutputLastLine := dockerApply(t, environments, remoteSharedPath)

			if diff := deep.Equal(gotApplyOutputLastLine, tt.wantApplyOutputLastLine); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func dockerRun(t *testing.T, command string, params map[string]string, sharedPath string) string {

	c := []string{command}
	for k, v := range params {
		c = append(c, fmt.Sprintf("%s=%s", k, v))
	}

	opts := &docker.RunOptions{
		Command: c,
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	//in case of error Run function calls FailNow anyways
	return docker.Run(t, imageTag, opts)
}

func dockerInit(t *testing.T, initParams map[string]string, sharedPath string) string {
	return dockerRun(t, "init", initParams, sharedPath)
}

func dockerPlan(t *testing.T, planParams map[string]string, sharedPath string) string {
	return getLastLineFromMultilineSting(t, dockerRun(t, "plan", planParams, sharedPath))
}

func dockerApply(t *testing.T, applyParams map[string]string, sharedPath string) string {
	return getLastLineFromMultilineSting(t, dockerRun(t, "apply", applyParams, sharedPath))
}

func setup(t *testing.T, rsaName string, name string) (string, string, map[string]string) {
	environments := loadEnvironments(t)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var remoteSharedPath string
	if v, ok := environments["K8S_HOST_PATH"]; ok && v != "" {
		remoteSharedPath = v
	} else {
		remoteSharedPath = path.Join(wd, "tests", "shared")
	}
	var localSharedPath string
	if v, ok := environments["K8S_VOL_PATH"]; ok && v != "" {
		localSharedPath = v
	} else {
		localSharedPath = path.Join(wd, "tests", "shared")
	}
	err = os.MkdirAll(remoteSharedPath, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	generateRsaKeyPair(t, remoteSharedPath, rsaName)
	if isResourceGroupPresent(t, environments["M_ARM_SUBSCRIPTION_ID"], name) {
		removeResourceGroup(t, environments["M_ARM_SUBSCRIPTION_ID"], name)
	}
	return remoteSharedPath, localSharedPath, environments
}

func removeResourceGroup(t *testing.T, subscriptionId string, name string) {
	t.Logf("Will prepare new az groups client")
	ctx := context.TODO()
	groupsClient := resources.NewGroupsClient(subscriptionId)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	groupsClient.Authorizer = authorizer
	rgName := fmt.Sprintf("%s-rg", name)
	t.Logf("Will perform delete RG operation")
	gdf, err := groupsClient.Delete(ctx, rgName)
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	now := time.Now()

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		defer close(done)
		t.Logf("Will start waiting for RG deletion finish.")
		err = gdf.Future.WaitForCompletionRef(ctx, groupsClient.BaseClient.Client)
		t.Logf("Finished RG deletion.")
		if err != nil {
			t.Fatal(err)
		}
	}()

	for {
		select {
		case <-ticker.C:
			t.Logf("Waiting for deletion to complete: %v", time.Since(now).Round(time.Second))
		case <-done:
			t.Logf("Finished waiting for RG deletion.")
			ticker.Stop()
			return
		}
	}
}

func isResourceGroupPresent(t *testing.T, subscriptionId string, name string) bool {
	groupsClient := resources.NewGroupsClient(subscriptionId)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		t.Error(err)
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

func cleanup(t *testing.T, sharedPath string, subscriptionId string, name string) {
	t.Logf("cleanup()")
	_ = os.RemoveAll(sharedPath)
	if isResourceGroupPresent(t, subscriptionId, name) {
		removeResourceGroup(t, subscriptionId, name)
	}
}

func generateRsaKeyPair(t *testing.T, directory string, name string) {
	privateRsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err)
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRsaKey)}
	privateKeyBytes := pem.EncodeToMemory(pemBlock)

	publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	err = ioutil.WriteFile(path.Join(directory, name), privateKeyBytes, 0600)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join(directory, fmt.Sprintf("%s.pub", name)), publicKeyBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func getLastLineFromMultilineSting(t *testing.T, s string) string {
	in := strings.NewReader(s)
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if err == io.EOF {
			return string(line)
		}
	}
}

func loadEnvironments(t *testing.T) map[string]string {
	result := make(map[string]string)
	result["M_ARM_CLIENT_ID"] = os.Getenv("AZURE_CLIENT_ID")
	if len(result["M_ARM_CLIENT_ID"]) == 0 {
		t.Fatalf("expected AZURE_CLIENT_ID environment variable")
	}
	result["M_ARM_CLIENT_SECRET"] = os.Getenv("AZURE_CLIENT_SECRET")
	if len(result["M_ARM_CLIENT_SECRET"]) == 0 {
		t.Fatalf("expected AZURE_CLIENT_SECRET environment variable")
	}
	result["M_ARM_SUBSCRIPTION_ID"] = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(result["M_ARM_SUBSCRIPTION_ID"]) == 0 {
		t.Fatalf("expected AZURE_SUBSCRIPTION_ID environment variable")
	}
	result["M_ARM_TENANT_ID"] = os.Getenv("AZURE_TENANT_ID")
	if len(result["M_ARM_TENANT_ID"]) == 0 {
		t.Fatalf("expected AZURE_TENANT_ID environment variable")
	}
	result["K8S_VOL_PATH"] = os.Getenv("K8S_VOL_PATH")
	result["K8S_HOST_PATH"] = os.Getenv("K8S_HOST_PATH")
	return result
}
