package cmd

import (
	"errors"
	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	state "github.com/epiphany-platform/e-structures/state/v0"
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
		return os.Rename(path, backupPath)
	}
}
