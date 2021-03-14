// SPDX-License-Identifier: BSD-3-Clause

package lava

import "encoding/base64"

type DeviceTypesListing struct {
	Devices   int    `xmlrpc:"devices"`
	Installed bool   `xmlrpc:"installed"`
	Name      string `xmlrpc:"name"`
	Template  bool   `xmlrpc:"template"`
}

func (c Connection) DevicesTypesList(showAll bool) ([]DeviceTypesListing, error) {
	var ret []DeviceTypesListing

	err := c.con.Call("scheduler.device_types.list", showAll, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c Connection) DevicesTypesTemplateSet(name string, template string) error {
	var ret []DeviceTypesListing
	var args []interface{}
	args = append(args, name)
	args = append(args, template)

	err := c.con.Call("scheduler.device_types.set_template", args, &ret)
	if err != nil {
		return err
	}

	return nil
}

func (c Connection) DevicesTypesTemplateGet(name string) (string, error) {
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

func (c Connection) DevicesTypesHealthCheckSet(name string, template string) error {
	var ret []DeviceTypesListing
	var args []interface{}
	args = append(args, name)
	args = append(args, template)

	err := c.con.Call("scheduler.device_types.set_health_check", args, &ret)
	if err != nil {
		return err
	}

	return nil
}

func (c Connection) DevicesTypesHealthCheckGet(name string) (string, error) {
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
