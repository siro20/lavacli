// SPDX-License-Identifier: BSD-3-Clause

package lava

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
