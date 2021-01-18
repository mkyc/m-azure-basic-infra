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
	type dd struct {
		id   string
		name string
		size int
	}
	type dda struct {
		lun              int
		managedDiskId    string
		virtualMachineId string
	}
	output := &azbi.Output{
		RgName:   to.StrPtr(m["rg_name"].(string)),
		VnetName: to.StrPtr(m["vnet_name"].(string)),
	}
	for _, i := range m["vm_groups"].([]interface{}) {
		vmGroup := i.(map[string]interface{})
		outputVmGroup := azbi.OutputVmGroup{
			Name: to.StrPtr(vmGroup["vm_group_name"].(string)),
		}
		dataDisks := make([]dd, 0)
		for _, j := range vmGroup["data_disks"].([]interface{}) {
			intermediateDataDisk := j.(map[string]interface{})
			dataDisks = append(dataDisks,
				dd{
					id:   intermediateDataDisk["id"].(string),
					name: intermediateDataDisk["name"].(string),
					size: int(intermediateDataDisk["size"].(float64)),
				})
		}
		logger.Debug().Msgf("Intermediate data disks struct list: %#v", dataDisks)
		dataDiskAttachments := make([]dda, 0)
		for _, j := range vmGroup["dd_attachments"].([]interface{}) {
			intermediateDataDiskAttachment := j.(map[string]interface{})
			dataDiskAttachments = append(dataDiskAttachments,
				dda{
					lun:              int(intermediateDataDiskAttachment["lun"].(float64)),
					managedDiskId:    intermediateDataDiskAttachment["managed_disk_id"].(string),
					virtualMachineId: intermediateDataDiskAttachment["virtual_machine_id"].(string),
				})
		}
		logger.Debug().Msgf("Intermediate data disk attachments struct list: %#v", dataDiskAttachments)
		for _, j := range vmGroup["vms"].([]interface{}) {
			intermediateVm := j.(map[string]interface{})
			outputVm := azbi.OutputVm{
				Name:     to.StrPtr(intermediateVm["vm_name"].(string)),
				PublicIp: to.StrPtr(intermediateVm["public_ip"].(string)),
			}
			for _, k := range intermediateVm["private_ips"].([]interface{}) {
				outputVm.PrivateIps = append(outputVm.PrivateIps, k.(string))
			}
			vmId := intermediateVm["id"].(string)
			for _, dda := range dataDiskAttachments {
				if dda.virtualMachineId == vmId {
					for _, dd := range dataDisks {
						if dd.id == dda.managedDiskId {
							outputVm.DataDisks = append(outputVm.DataDisks, azbi.OutputDataDisk{
								Size: to.IntPtr(dd.size),
								Lun:  to.IntPtr(dda.lun),
							})
						}
					}
				}
			}
			outputVmGroup.Vms = append(outputVmGroup.Vms, outputVm)
		}
		output.VmGroups = append(output.VmGroups, outputVmGroup)
	}
	return output
}
