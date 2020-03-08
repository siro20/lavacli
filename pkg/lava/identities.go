// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"fmt"
)

type LavaIndentity struct {
	Name     string
	Token    string
	Uri      string
	Username string
	Proxy    string
}

func LavaIdentitiesList() ([]LavaIndentity, error) {
	var ret []LavaIndentity
	configs, err := LavaGetConf()
	if err != nil {
		return nil, err
	}
	for k, v := range configs {
		ret = append(ret, LavaIndentity{
			k,
			v.Token,
			v.Uri,
			v.Username,
			v.Proxy,
		})
	}

	return ret, nil
}

func LavaIdentitiesAdd(id LavaIndentity) error {
	configs, err := LavaGetConf()
	if err != nil {
		return err
	}
	for k, _ := range configs {
		if k == id.Name {
			return fmt.Errorf("id %s is already in config", id.Name)
		}
	}
	if id.Uri == "" {
		return fmt.Errorf("Must specify URI in identity")
	}
	var c LavaConfigIndentity
	c.Uri = id.Uri
	c.Token = id.Token
	c.Username = id.Username
	c.Proxy = id.Proxy

	configs[id.Name] = c

	return LavaSetConf(configs)
}

func LavaIdentitiesShow(name string) (*LavaIndentity, error) {
	var ret LavaIndentity
	configs, err := LavaGetConf()
	if err != nil {
		return nil, err
	}
	for k, v := range configs {
		if k == name {
			ret = LavaIndentity{
				k,
				v.Token,
				v.Uri,
				v.Username,
				v.Proxy,
			}
			return &ret, nil
		}
	}

	return nil, fmt.Errorf("id %s not found in config", name)
}

func LavaIdentitiesDelete(name string) error {
	configs, err := LavaGetConf()
	if err != nil {
		return err
	}
	for k, _ := range configs {
		if k == name {

			new := map[string]LavaConfigIndentity{}
			for k2, v2 := range configs {
				if k2 != name {
					new[k2] = v2
				}
			}

			return LavaSetConf(new)
		}
	}

	return fmt.Errorf("id %s not found in config", name)
}
