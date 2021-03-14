// SPDX-License-Identifier: BSD-3-Clause

package lava

type DeviceList struct {
	Hostname   string `xmlrpc:"hostname"`
	Type       string `xmlrpc:"type"`
	State      string `xmlrpc:"state"`
	Health     string `xmlrpc:"health"`
	CurrentJob int    `xmlrpc:"current_job"`
	Pipeline   bool   `xmlrpc:"pipeline"`
}

func (c Connection) DevicesList() ([]DeviceList, error) {
	var ret []DeviceList

	err := c.con.Call("scheduler.devices.list", nil, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type Device struct {
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

func (c Connection) DevicesShow(hostname string) (*Device, error) {
	var ret Device

	err := c.con.Call("scheduler.devices.show", hostname, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c Connection) DevicesTagsList(hostname string) ([]string, error) {
	var ret []string

	err := c.con.Call("scheduler.devices.tags.list", hostname, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c Connection) DevicesTagsDelete(hostname string, name string) error {
	var args []interface{}
	args = append(args, hostname)
	args = append(args, name)

	return c.con.Call("scheduler.devices.tags.delete", args, nil)
}

func (c Connection) DevicesTagsAdd(hostname string, name string) error {
	var args []interface{}
	args = append(args, hostname)
	args = append(args, name)

	return c.con.Call("scheduler.devices.tags.add", args, nil)
}
