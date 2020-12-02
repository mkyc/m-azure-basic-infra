package cmd

import (
	"errors"
	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	state "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
	"io/ioutil"
	"os"
)

func ensureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func loadState(path string) (*state.State, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return state.NewState(), nil
	} else {
		s := &state.State{}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = s.Unmarshall(b)
		if err != nil {
			return nil, err
		}
		return s, nil
	}
}

func saveState(path string, s *state.State) error {
	b, err := s.Marshall()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(path string) (*azbi.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return azbi.NewConfig(), nil
	} else {
		c := &azbi.Config{}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = c.Unmarshall(b)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

func saveConfig(path string, c *azbi.Config) error {
	b, err := c.Marshall()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func checkAndLoad(stateFilePath string, configFilePath string) (*azbi.Config, *state.State, error) {
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("state file does not exist, please run init first")
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("config file does not exist, please run init first")
	}

	s, err := loadState(stateFilePath)
	if err != nil {
		return nil, nil, err
	}

	c, err := loadConfig(configFilePath)
	if err != nil {
		return nil, nil, err
	}

	return c, s, nil
}

func backupFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	} else {
		backupPath := path + ".backup"

		input, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(backupPath, input, 0644)
		if err != nil {
			return err
		}
		return nil
	}
}

func produceOutput(m map[string]interface{}) *azbi.Output {
	o := &azbi.Output{
		RgName:   to.StrPtr(m["rg_name"].(string)),
		VnetName: to.StrPtr(m["vnet_name"].(string)),
	}
	for _, i := range m["private_ips"].([]interface{}) {
		o.PrivateIps = append(o.PrivateIps, i.(string))
	}
	for _, i := range m["public_ips"].([]interface{}) {
		o.PublicIps = append(o.PublicIps, i.(string))
	}
	for _, i := range m["vm_names"].([]interface{}) {
		o.VmNames = append(o.VmNames, i.(string))
	}
	return o
}
