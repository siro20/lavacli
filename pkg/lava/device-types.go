// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"github.com/kolo/xmlrpc"
)

type LavaDeviceTypesListing struct {
	Devices   int    `xmlrpc:"devices"`
	Installed bool   `xmlrpc:"installed"`
	Name      string `xmlrpc:"name"`
	Template  bool   `xmlrpc:"template"`
}

func LavaDevicesTypesList(con *xmlrpc.Client, showAll bool) ([]LavaDeviceTypesListing, error) {
	var ret []LavaDeviceTypesListing

	err := con.Call("scheduler.device_types.list", showAll, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
