// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type configIndentity struct {
	Token    string `yaml:"token,omitempty"`
	Uri      string `yaml:"uri,omitempty"`
	Username string `yaml:"username,omitempty"`
	Proxy    string `yaml:"proxy,omitempty"`
}

func GetConf() map[string]configIndentity {
	var c map[string]configIndentity

	path := os.Getenv("XDG_CONFIG_HOME")
	if path == "" {
		usr, _ := user.Current()
		path = usr.HomeDir + "/.config"
	}

	path += "/lavacli.yaml"
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to resolv path: #%v ", err)
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}

	return c
}

func SetConf(c map[string]configIndentity) {
	path := os.Getenv("XDG_CONFIG_HOME")
	if path == "" {
		usr, _ := user.Current()
		path = usr.HomeDir + "/.config"
	}

	path += "/lavacli.yaml"
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to resolv path: #%v ", err)
	}

	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
