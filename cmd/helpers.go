package cmd

import (
	"errors"
	"io/ioutil"
	"os"

	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
)

func ensureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func loadState(path string) (*st.State, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return st.NewState(), nil
	} else {
		state := &st.State{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = state.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return state, nil
	}
}

func saveState(path string, state *st.State) error {
	bytes, err := state.Marshal()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(path string) (*azbi.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return azbi.NewConfig(), nil
	} else {
		config := &azbi.Config{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = config.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}

func saveConfig(path string, config *azbi.Config) error {
	bytes, err := config.Marshal()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func checkAndLoad(stateFilePath string, configFilePath string) (*azbi.Config, *st.State, error) {
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("state file does not exist, please run init first")
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("config file does not exist, please run init first")
	}

	state, err := loadState(stateFilePath)
	if err != nil {
		return nil, nil, err
	}

	config, err := loadConfig(configFilePath)
	if err != nil {
		return nil, nil, err
	}

	return config, state, nil
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
	logger.Debug().Msgf("Received output map: %#v", m)
	output := &azbi.Output{
		RgName:   to.StrPtr(m["rg_name"].(string)),
		VnetName: to.StrPtr(m["vnet_name"].(string)),
	}
	for _, i := range m["vm_groups"].([]interface{}) {
		vmGroup := i.(map[string]interface{})
		outputVmGroup := azbi.OutputVmGroup{
			Name: vmGroup["vm_group_name"].(string),
		}
		for _, j := range vmGroup["vms"].([]interface{}) {
			vm := j.(map[string]interface{})
			outputVm := azbi.OutputVm{
				Name:     vm["vm_name"].(string),
				PublicIp: vm["public_ip"].(string),
			}
			for _, k := range vm["private_ips"].([]interface{}) {
				outputVm.PrivateIps = append(outputVm.PrivateIps, k.(string))
			}
			outputVmGroup.Vms = append(outputVmGroup.Vms, outputVm)
		}
		output.VmGroups = append(output.VmGroups, outputVmGroup)
	}
	return output
}
