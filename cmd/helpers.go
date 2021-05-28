package cmd

import (
	"errors"
	"io/ioutil"
	"os"

	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/load"
	"github.com/epiphany-platform/e-structures/utils/to"
)

func ensureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func checkAndLoad(stateFilePath string, configFilePath string) (*azbi.Config, *st.State, error) {
	logger.Debug().Msgf("checkAndLoad(%s, %s)", stateFilePath, configFilePath)
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("state file does not exist, please run init first")
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("config file does not exist, please run init first")
	}

	state, err := load.State(stateFilePath)
	if err != nil {
		return nil, nil, err
	}

	config, err := load.AzBIConfig(configFilePath)
	if err != nil {
		return nil, nil, err
	}

	return config, state, nil
}

func backupFile(path string) error {
	logger.Debug().Msgf("backupFile(%s)", path)
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

	// two internal intermediate data structures to hold extracted map values
	type intermediateDataDisk struct {
		id   string
		name string
		size int
	}
	type intermediateDataDiskAttachment struct {
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
		intermediateDataDisks := make([]intermediateDataDisk, 0)
		for _, j := range vmGroup["data_disks"].([]interface{}) {
			tempDataDisk := j.(map[string]interface{})
			intermediateDataDisks = append(intermediateDataDisks,
				intermediateDataDisk{
					id:   tempDataDisk["id"].(string),
					name: tempDataDisk["name"].(string),
					size: int(tempDataDisk["size"].(float64)),
				})
		}
		logger.Debug().Msgf("Intermediate data disks struct list: %#v", intermediateDataDisks)
		intermediateDataDiskAttachments := make([]intermediateDataDiskAttachment, 0)
		for _, j := range vmGroup["dd_attachments"].([]interface{}) {
			tempDataDiskAttachment := j.(map[string]interface{})
			intermediateDataDiskAttachments = append(intermediateDataDiskAttachments,
				intermediateDataDiskAttachment{
					lun:              int(tempDataDiskAttachment["lun"].(float64)),
					managedDiskId:    tempDataDiskAttachment["managed_disk_id"].(string),
					virtualMachineId: tempDataDiskAttachment["virtual_machine_id"].(string),
				})
		}
		logger.Debug().Msgf("Intermediate data disk attachments struct list: %#v", intermediateDataDiskAttachments)
		for _, j := range vmGroup["vms"].([]interface{}) {
			tempVm := j.(map[string]interface{})
			outputVm := azbi.OutputVm{
				Name:     to.StrPtr(tempVm["vm_name"].(string)),
				PublicIp: to.StrPtr(tempVm["public_ip"].(string)),
			}
			for _, k := range tempVm["private_ips"].([]interface{}) {
				outputVm.PrivateIps = append(outputVm.PrivateIps, k.(string))
			}
			vmId := tempVm["id"].(string)
			for _, dda := range intermediateDataDiskAttachments {
				if dda.virtualMachineId == vmId {
					for _, dd := range intermediateDataDisks {
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
