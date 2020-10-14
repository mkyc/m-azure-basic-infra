package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
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

	tag := "epiphanyplatform/azbi:0.0.1"

	runOpts := &docker.RunOptions{
		Command: []string{"init"},
		Remove:  true,
		Volumes: []string{fmt.Sprintf("%s:/shared", sharedPath)},
	}

	docker.Run(t, tag, runOpts)
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
