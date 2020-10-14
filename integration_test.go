package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/docker"
)

const (
	imageTag = "epiphanyplatform/azbi:0.0.1"
)

var (
	sharedPath = ""
)

func TestMain(m *testing.M) {
	p, err := setup()
	if err != nil {
		fmt.Printf("setup() failed with: %v\n", err)
		os.Exit(1)
	}
	sharedPath = p
	code := m.Run()
	_ = cleanup()
	os.Exit(code)
}

func TestInitDefaultConfig(t *testing.T) {
	tests := []struct {
		name            string
		initParams      []string
		wantOutput      string
		wantFile        string
		wantFileContent string
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
			wantFile: "azbi/azbi-config.yml",
			wantFileContent: `kind: azbi-config
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := []string{"init"}
			command = append(command, tt.initParams...)
			runOpts := &docker.RunOptions{
				Command: command,
				Remove:  true,
				Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
			}

			output := docker.Run(t, imageTag, runOpts)
			if diff := deep.Equal(output, tt.wantOutput); diff != nil {
				t.Error(diff)
			}

			expectedPath := path.Join(sharedPath, tt.wantFile)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Fatalf("missing expected file: %s", expectedPath)
			}
			gotFileContent, err := ioutil.ReadFile(expectedPath)
			if err != nil {
				t.Errorf("wasnt able to read form output file: %v", err)
			}
			if diff := deep.Equal(string(gotFileContent), tt.wantFileContent); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func setup() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	p := path.Join(wd, "tests", "shared")
	return p, os.MkdirAll(p, os.ModePerm)
}

func cleanup() error {
	return os.RemoveAll(sharedPath)
}
