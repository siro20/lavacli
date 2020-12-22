// SPDX-License-Identifier: BSD-3-Clause

package lava

import "encoding/base64"

type LavaDeviceTypesListing struct {
	Devices   int    `xmlrpc:"devices"`
	Installed bool   `xmlrpc:"installed"`
	Name      string `xmlrpc:"name"`
	Template  bool   `xmlrpc:"template"`
}

func (c LavaConnection) LavaDevicesTypesList(showAll bool) ([]LavaDeviceTypesListing, error) {
	var ret []LavaDeviceTypesListing

	err := c.con.Call("scheduler.device_types.list", showAll, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c LavaConnection) LavaDevicesTypesTemplateSet(name string, template string) error {
	var ret []LavaDeviceTypesListing
	var args []interface{}
	args = append(args, name)
	args = append(args, template)

	err := c.con.Call("scheduler.device_types.set_template", args, &ret)
	if err != nil {
		return err
	}

	return nil
}

func (c LavaConnection) LavaDevicesTypesTemplateGet(name string) (string, error) {
	var data string

	err := c.con.Call("scheduler.device_types.get_template", name, &data)
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func (c LavaConnection) LavaDevicesTypesHealthCheckSet(name string, template string) error {
	var ret []LavaDeviceTypesListing
	var args []interface{}
	args = append(args, name)
	args = append(args, template)

	err := c.con.Call("scheduler.device_types.set_health_check", args, &ret)
	if err != nil {
		return err
	}

	return nil
}

func (c LavaConnection) LavaDevicesTypesHealthCheckGet(name string) (string, error) {
	var data string

	err := c.con.Call("scheduler.device_types.get_health_check", name, &data)
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
