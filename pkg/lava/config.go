// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type LavaConfigIndentity struct {
	Token    string `yaml:"token,omitempty"`
	Uri      string `yaml:"uri,omitempty"`
	Username string `yaml:"username,omitempty"`
	Proxy    string `yaml:"proxy,omitempty"`
}

func LavaGetConf() (map[string]LavaConfigIndentity, error) {
	var c map[string]LavaConfigIndentity

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

func LavaSetConf(c map[string]LavaConfigIndentity) error {
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
