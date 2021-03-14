// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"fmt"
)

type Indentity struct {
	Name     string
	Token    string
	URI      string
	Username string
	Proxy    string
}

func IdentitiesList() ([]Indentity, error) {
	var ret []Indentity
	configs, err := GetConf()
	if err != nil {
		return nil, err
	}
	for k, v := range configs {
		ret = append(ret, Indentity{
			k,
			v.Token,
			v.URI,
			v.Username,
			v.Proxy,
		})
	}

	return ret, nil
}

func IdentitiesAdd(id Indentity) error {
	configs, err := GetConf()
	if err != nil {
		return err
	}
	for k, _ := range configs {
		if k == id.Name {
			return fmt.Errorf("id %s is already in config", id.Name)
		}
	}
	if id.URI == "" {
		return fmt.Errorf("Must specify URI in identity")
	}
	var c ConfigIndentity
	c.URI = id.URI
	c.Token = id.Token
	c.Username = id.Username
	c.Proxy = id.Proxy

	configs[id.Name] = c

	return SetConf(configs)
}

func IdentitiesShow(name string) (*Indentity, error) {
	var ret Indentity
	configs, err := GetConf()
	if err != nil {
		return nil, err
	}
	for k, v := range configs {
		if k == name {
			ret = Indentity{
				k,
				v.Token,
				v.URI,
				v.Username,
				v.Proxy,
			}
			return &ret, nil
		}
	}

	return nil, fmt.Errorf("id %s not found in config", name)
}

func IdentitiesDelete(name string) error {
	configs, err := GetConf()
	if err != nil {
		return err
	}
	for k, _ := range configs {
		if k == name {

			new := map[string]ConfigIndentity{}
			for k2, v2 := range configs {
				if k2 != name {
					new[k2] = v2
				}
			}

			return SetConf(new)
		}
	}

	return fmt.Errorf("id %s not found in config", name)
}
