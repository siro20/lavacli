// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//ConfigIndentity represent a field within the lavacli.yaml
type ConfigIndentity struct {
	Token    string `yaml:"token,omitempty"`
	URI      string `yaml:"uri,omitempty"`
	Username string `yaml:"username,omitempty"`
	Proxy    string `yaml:"proxy,omitempty"`
}

// GetConf loads the lavacli.yaml
func GetConf() (map[string]ConfigIndentity, error) {
	var c map[string]ConfigIndentity

	path := os.Getenv("XDG_CONFIG_HOME")
	if path == "" {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		path = usr.HomeDir + "/.config"
	}

	path += "/lavacli.yaml"
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// SetConf saves the config to lavacli.yaml
func SetConf(c map[string]ConfigIndentity) error {
	path := os.Getenv("XDG_CONFIG_HOME")
	if path == "" {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		path = usr.HomeDir + "/.config"
	}

	path += "/lavacli.yaml"
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	d, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		return err
	}

	return nil
}
