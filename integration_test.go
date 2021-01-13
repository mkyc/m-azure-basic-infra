package main

import (
	"bufio"
	"bytes"
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
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/m-azure-basic-infrastructure/cmd"
	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/docker"
	"golang.org/x/crypto/ssh"
)

func TestMetadata(t *testing.T) {
	tests := []struct {
		name               string
		wantOutputTemplate string
	}{
		{
			name: "default metadata",
			wantOutputTemplate: `labels:
  kind: infrastructure
  name: Azure Basic Infrastructure
  provider: azure
  provides-pubips: true
  provides-vms: true
  short: azbi
  version: %s
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := dockerRun(t, "metadata", nil, nil, "")
			if diff := deep.Equal(gotOutput, fmt.Sprintf(tt.wantOutputTemplate, cmd.Version)); diff != nil {
				t.Error(diff)
			}
		})
	}
}

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
			wantOutput: `Initialized config: 
{
	"kind": "azbi",
	"version": "v0.1.0",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "main",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			}
		],
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 1,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": true,
			"subnet_names": ["main"],
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}`,
			wantConfigLocation: "azbi/azbi-config.json",
			wantConfigContent: `{
	"kind": "azbi",
	"version": "v0.1.0",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "main",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			}
		],
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 1,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": true,
			"subnet_names": ["main"],
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}`,
		},
		{
			name: "pass name and vms_rsa cli arguments",
			initParams: map[string]string{
				"--name":    "azbi-module-tests",
				"--vms_rsa": "test_vms_rsa"},
			wantOutput: `Initialized config: 
{
	"kind": "azbi",
	"version": "v0.0.2",
	"params": {
		"name": "azbi-module-tests",
		"location": "northeurope",
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "main",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			}
		],
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 1,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": true,
			"subnet_names": ["main"],
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}],
		"rsa_pub_path": "/shared/test_vms_rsa.pub"
	}
}`,
			wantConfigLocation: "azbi/azbi-config.json",
			wantConfigContent: `{
	"kind": "azbi",
	"version": "v0.1.0",
	"params": {
		"name": "azbi-module-tests",
		"location": "northeurope",
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "main",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			}
		],
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 1,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": true,
			"subnet_names": ["main"],
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}],
		"rsa_pub_path": "/shared/test_vms_rsa.pub"
	}
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, remoteSharedPath, localSharedPath, environments, _ := setup(t, tt.initParams)
			defer cleanup(t, localSharedPath, environments["SUBSCRIPTION_ID"], name)

			gotOutput := dockerRun(t, "init", tt.initParams, nil, remoteSharedPath)
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
			name: "default plan",
			initParams: map[string]string{
				"--name":    "azbi-module-tests",
				"--vms_rsa": "test_vms_rsa"},
			wantPlanOutputLastLine: "\tAdd: 8, Change: 0, Destroy: 0",
			wantStateLocation:      "state.json",
			wantStateContent: `{
	"kind": "state",
	"version": "v0.0.2",
	"azbi": {
		"status": "initialized",
		"config": null,
		"output": null
	}
}`,
			wantTerraformStateFileLocation: "azbi/terraform-apply.tfplan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, remoteSharedPath, localSharedPath, environments, _ := setup(t, tt.initParams)
			defer cleanup(t, localSharedPath, environments["SUBSCRIPTION_ID"], name)

			dockerRun(t, "init", tt.initParams, nil, remoteSharedPath)

			gotPlanOutputLastLine := getLastLineFromMultilineSting(t, dockerRun(t, "plan", nil, environments, remoteSharedPath))
			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			expectedPath := path.Join(localSharedPath, tt.wantStateLocation)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Fatalf("missing expected file: %s", expectedPath)
			}

			gotStateContent, err := ioutil.ReadFile(expectedPath)
			if err != nil {
				t.Errorf("wasnt able to read form output file: %v", err)
			}
			if diff := deep.Equal(string(gotStateContent), tt.wantStateContent); diff != nil {
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
			name: "default apply",
			initParams: map[string]string{
				"--name":    "azbi-module-tests",
				"--vms_rsa": "test_vms_rsa"},
			wantPlanOutputLastLine:  "\tAdd: 8, Change: 0, Destroy: 0",
			wantApplyOutputLastLine: "\tAdd: 8, Change: 0, Destroy: 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, remoteSharedPath, localSharedPath, environments, privateKey := setup(t, tt.initParams)
			defer cleanup(t, localSharedPath, environments["SUBSCRIPTION_ID"], name)

			dockerRun(t, "init", tt.initParams, nil, remoteSharedPath)

			gotPlanOutputLastLine := getLastLineFromMultilineSting(t, dockerRun(t, "plan", nil, environments, remoteSharedPath))
			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			gotApplyOutputLastLine := getLastLineFromMultilineSting(t, dockerRun(t, "apply", nil, environments, remoteSharedPath))

			if diff := deep.Equal(gotApplyOutputLastLine, tt.wantApplyOutputLastLine); diff != nil {
				t.Error(diff)
			}

			// public IPs are enabled by default
			data, err := ioutil.ReadFile(path.Join(localSharedPath, "state.json"))
			if err != nil {
				t.Fatal(err)
			}
			state := &st.State{}
			err = state.Unmarshal(data)
			if err != nil {
				t.Fatal(err)
			}
			// check connectivity to all VMs from default VM group
			Vms := state.AzBI.Output.VmGroups[0].Vms
			for _, vm := range Vms {
				validateSshConnectivity(t, privateKey, *vm.PublicIp)
			}
		})
	}
}

// dockerRun function wraps docker run operation and returns `docker run` output.
func dockerRun(t *testing.T, command string, parameters map[string]string, environments map[string]string, sharedPath string) string {
	commandWithParameters := []string{command}
	for k, v := range parameters {
		commandWithParameters = append(commandWithParameters, fmt.Sprintf("%s=%s", k, v))
	}

	var opts *docker.RunOptions
	if sharedPath != "" {
		opts = &docker.RunOptions{
			Command: commandWithParameters,
			Remove:  true,
			Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
		}
	} else {
		opts = &docker.RunOptions{
			Command: commandWithParameters,
			Remove:  true,
		}
	}
	var envs []string
	for k, v := range environments {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	opts.EnvironmentVariables = envs

	//in case of error Run function calls FailNow anyways
	return docker.Run(t, fmt.Sprintf("%s:%s", prepareImageTag(t), cmd.Version), opts)
}

// setup function ensures that all prerequisites for tests are in place.
func setup(t *testing.T, initParams map[string]string) (string, string, string, map[string]string, ssh.Signer) {
	rsaName := "vms_rsa"
	if value, ok := initParams["--vms_rsa"]; ok {
		rsaName = value
	}
	name := "epiphany-rg"
	if value, ok := initParams["--name"]; ok {
		name = value
	}

	environments := loadEnvironmentVariables(t)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	var remoteSharedPath string
	if v, ok := environments["K8S_VOL_PATH"]; ok && v != "" {
		remoteSharedPath = v
	} else {
		remoteSharedPath = path.Join(wd, "shared")
	}
	var localSharedPath string
	if v, ok := environments["K8S_HOST_PATH"]; ok && v != "" {
		localSharedPath = v
	} else {
		localSharedPath = path.Join(wd, "shared")
	}
	err = os.MkdirAll(localSharedPath, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	privateKey := generateRsaKeyPair(t, localSharedPath, rsaName)
	if isResourceGroupPresent(t, environments["SUBSCRIPTION_ID"], name) {
		removeResourceGroup(t, environments["SUBSCRIPTION_ID"], name)
	}
	return name, remoteSharedPath, localSharedPath, environments, privateKey
}

// prepareImageTag returns IMAGE_REPOSITORY environment variable
func prepareImageTag(t *testing.T) string {
	imageRepository := os.Getenv("IMAGE_REPOSITORY")
	if len(imageRepository) == 0 {
		t.Fatal("expected IMAGE_REPOSITORY environment variable")
	}
	return imageRepository
}

// cleanup function removes directories created during test and ensures that resource
// group gets removed if it was created.
func cleanup(t *testing.T, sharedPath string, subscriptionId string, name string) {
	t.Logf("cleanup()")
	_ = os.RemoveAll(sharedPath)
	if isResourceGroupPresent(t, subscriptionId, name) {
		removeResourceGroup(t, subscriptionId, name)
	}
}

// isResourceGroupPresent function checks if resource group with given name exists.
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

// removeResourceGroup function invokes Delete operation on provided resource
// group name and waits for operation completion.
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

// generateRsaKeyPair function generates RSA public and private keys and returns
// ssh.Signer that can create signatures that verify against a public key.
func generateRsaKeyPair(t *testing.T, directory string, name string) ssh.Signer {
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
	signer, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		t.Fatal(err)
	}
	return signer
}

// validateSshConnectivity function checks possibility to connect to provided IP
// address and run `uname -a` command. In case of failed connection attempt or error
// while running command test will fail.
func validateSshConnectivity(t *testing.T, signer ssh.Signer, ipString string) {
	sshConfig := &ssh.ClientConfig{
		User:            "operations",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var (
		retries   = 30
		client    *ssh.Client
		connected = false
		err       error
	)

	for i := 0; i < retries; i++ {
		client, err = ssh.Dial("tcp", fmt.Sprintf("%s:22", ipString), sshConfig)
		if err != nil {
			t.Log(err)
			time.Sleep(1 * time.Second)
		} else {
			connected = true
		}
	}
	if !connected {
		t.Error(err)
	}
	session, err := client.NewSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run("uname -a")
	if err != nil {
		t.Error()
	}
	t.Logf("ssh connectivity test result: %s", b.String())
}

// getLastLineFromMultilineSting is helper function to extract just last line
// from multiline string.
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

// loadEnvironmentVariables obtains 6 variables from environment.
// Two of them (K8S_VOL_PATH and K8S_HOST_PATH) are optional and
// are not checked but another four (AZURE_CLIENT_ID, AZURE_CLIENT_SECRET
// AZURE_SUBSCRIPTION_ID and AZURE_TENANT_ID) are required and if missing
// will cause test to fail.
func loadEnvironmentVariables(t *testing.T) map[string]string {
	result := make(map[string]string)
	result["CLIENT_ID"] = os.Getenv("AZURE_CLIENT_ID")
	if len(result["CLIENT_ID"]) == 0 {
		t.Fatalf("expected AZURE_CLIENT_ID environment variable")
	}
	result["CLIENT_SECRET"] = os.Getenv("AZURE_CLIENT_SECRET")
	if len(result["CLIENT_SECRET"]) == 0 {
		t.Fatalf("expected AZURE_CLIENT_SECRET environment variable")
	}
	result["SUBSCRIPTION_ID"] = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(result["SUBSCRIPTION_ID"]) == 0 {
		t.Fatalf("expected AZURE_SUBSCRIPTION_ID environment variable")
	}
	result["TENANT_ID"] = os.Getenv("AZURE_TENANT_ID")
	if len(result["TENANT_ID"]) == 0 {
		t.Fatalf("expected AZURE_TENANT_ID environment variable")
	}
	result["K8S_VOL_PATH"] = os.Getenv("K8S_VOL_PATH")
	result["K8S_HOST_PATH"] = os.Getenv("K8S_HOST_PATH")
	return result
}
