// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"github.com/kolo/xmlrpc"
)

type LavaDeviceList struct {
	Hostname   string `xmlrpc:"hostname"`
	Type       string `xmlrpc:"type"`
	State      string `xmlrpc:"state"`
	Health     string `xmlrpc:"health"`
	CurrentJob int    `xmlrpc:"current_job"`
	Pipeline   bool   `xmlrpc:"pipeline"`
}

func LavaDevicesList(con *xmlrpc.Client) ([]LavaDeviceList, error) {
	var ret []LavaDeviceList

	err := con.Call("scheduler.devices.list", nil, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type LavaDevice struct {
	Description   string   `xmlrpc:"description"`
	HasDeviceDict bool     `xmlrpc:"has_device_dict"`
	HealthJob     bool     `xmlrpc:"health_job"`
	Worker        string   `xmlrpc:"worker"`
	Tags          []string `xmlrpc:"tags"`
	Hostname      string   `xmlrpc:"hostname"`
	DeviceType    string   `xmlrpc:"device_type"`
	State         string   `xmlrpc:"state"`
	Health        string   `xmlrpc:"health"`
	CurrentJob    int      `xmlrpc:"current_job"`
	Pipeline      bool     `xmlrpc:"pipeline"`
}

func LavaDevicesShow(con *xmlrpc.Client, hostname string) (*LavaDevice, error) {
	var ret LavaDevice

	err := con.Call("scheduler.devices.show", hostname, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
