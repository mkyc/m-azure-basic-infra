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
	"github.com/epiphany-platform/m-azure-basic-infrastructure/cmd"
	"gopkg.in/yaml.v2"
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
	imageTag = "epiphanyplatform/azbi"
)

func TestMetadata(t *testing.T) {
	tests := []struct {
		name               string
		wantOutputTemplate string
	}{
		{
			name: "default metadata",
			wantOutputTemplate: `#AzBI | metadata | should print component metadata
labels:
  version: %s
  name: Azure Basic Infrastructure
  short: azbi
  kind: infrastructure
  provider: azure
  provides-vms: true
  provides-pubips: true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := dockerRun(t, "metadata", nil, "")
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
			wantOutput: `kind: azbi-config
azbi:
  size: 3
  use_public_ip: true
  location: northeurope
  name: epiphany
  address_space: [10.0.0.0/16]
  address_prefixes: [10.0.1.0/24]
  rsa_pub_path: /shared/vms_rsa.pub`,
			wantConfigLocation: "azbi/azbi-config.yml",
			wantConfigContent: `kind: azbi-config
azbi:
  size: 3
  use_public_ip: true
  location: northeurope
  name: epiphany
  address_space: [10.0.0.0/16]
  address_prefixes: [10.0.1.0/24]
  rsa_pub_path: /shared/vms_rsa.pub
`,
		},
		{
			name: "init 2 machines no public ips and named rg",
			initParams: map[string]string{
				"M_VMS_COUNT":  "2",
				"M_PUBLIC_IPS": "false",
				"M_NAME":       "azbi-module-tests",
				"M_VMS_RSA":    "test_vms_rsa"},
			wantOutput: `kind: azbi-config
azbi:
  size: 2
  use_public_ip: false
  location: northeurope
  name: azbi-module-tests
  address_space: [10.0.0.0/16]
  address_prefixes: [10.0.1.0/24]
  rsa_pub_path: /shared/test_vms_rsa.pub`,
			wantConfigLocation: "azbi/azbi-config.yml",
			wantConfigContent: `kind: azbi-config
azbi:
  size: 2
  use_public_ip: false
  location: northeurope
  name: azbi-module-tests
  address_space: [10.0.0.0/16]
  address_prefixes: [10.0.1.0/24]
  rsa_pub_path: /shared/test_vms_rsa.pub
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, remoteSharedPath, localSharedPath, environments, _ := setup(t, tt.initParams)
			defer cleanup(t, localSharedPath, environments["M_ARM_SUBSCRIPTION_ID"], name)

			gotOutput := dockerRun(t, "init", tt.initParams, remoteSharedPath)
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
			name, remoteSharedPath, localSharedPath, environments, _ := setup(t, tt.initParams)
			defer cleanup(t, localSharedPath, environments["M_ARM_SUBSCRIPTION_ID"], name)

			dockerRun(t, "init", tt.initParams, remoteSharedPath)

			gotPlanOutputLastLine := getLastLineFromMultilineSting(t, dockerRun(t, "plan", environments, remoteSharedPath))
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
			name, remoteSharedPath, localSharedPath, environments, privateKey := setup(t, tt.initParams)
			defer cleanup(t, localSharedPath, environments["M_ARM_SUBSCRIPTION_ID"], name)

			dockerRun(t, "init", tt.initParams, remoteSharedPath)

			gotPlanOutputLastLine := getLastLineFromMultilineSting(t, dockerRun(t, "plan", environments, remoteSharedPath))
			if diff := deep.Equal(gotPlanOutputLastLine, tt.wantPlanOutputLastLine); diff != nil {
				t.Error(diff)
			}

			gotApplyOutputLastLine := getLastLineFromMultilineSting(t, dockerRun(t, "apply", environments, remoteSharedPath))

			if diff := deep.Equal(gotApplyOutputLastLine, tt.wantApplyOutputLastLine); diff != nil {
				t.Error(diff)
			}

			if v, ok := tt.initParams["M_PUBLIC_IPS"]; ok && v == "true" {
				data, err := ioutil.ReadFile(path.Join(localSharedPath, "state.yml"))
				if err != nil {
					t.Fatal(err)
				}
				m := make(map[interface{}]interface{})
				err = yaml.Unmarshal(data, &m)
				if err != nil {
					t.Fatal(err)
				}
				publicIPs := m["azbi"].(map[interface{}]interface{})["output"].(map[interface{}]interface{})["public_ips.value"].([]interface{})
				for _, p := range publicIPs {
					s := p.(string)
					validateSshConnectivity(t, privateKey, s)
				}
			}
		})
	}
}

// dockerRun function wraps docker run operation and returns `docker run` output.
func dockerRun(t *testing.T, command string, params map[string]string, sharedPath string) string {
	c := []string{command}
	for k, v := range params {
		c = append(c, fmt.Sprintf("%s=%s", k, v))
	}

	var opts *docker.RunOptions
	if sharedPath != "" {
		opts = &docker.RunOptions{
			Command: c,
			Remove:  true,
			Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
		}
	} else {
		opts = &docker.RunOptions{
			Command: c,
			Remove:  true,
		}
	}

	//in case of error Run function calls FailNow anyways
	return docker.Run(t, fmt.Sprintf("%s:%s", imageTag, cmd.Version), opts)
}

// setup function ensures that all prerequisites for tests are in place.
func setup(t *testing.T, initParams map[string]string) (string, string, string, map[string]string, ssh.Signer) {
	rsaName := "vms_rsa"
	if v, ok := initParams["M_VMS_RSA"]; ok {
		rsaName = v
	}
	name := "epiphany-rg"
	if v, ok := initParams["M_NAME"]; ok {
		name = v
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
	if isResourceGroupPresent(t, environments["M_ARM_SUBSCRIPTION_ID"], name) {
		removeResourceGroup(t, environments["M_ARM_SUBSCRIPTION_ID"], name)
	}
	return name, remoteSharedPath, localSharedPath, environments, privateKey
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
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ipString), sshConfig)
	if err != nil {
		t.Error(err)
	}
	session, err := connection.NewSession()
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
