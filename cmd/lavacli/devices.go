// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

type listDevicesCmd struct {
	Yaml bool `flag:"" optional:"" help:"Output as YAML" default:"false"`
	JSON bool `flag:"" optional:"" help:"Output as JSON" default:"false"`
}

func (c *listDevicesCmd) Run(ctx *context) error {
	ret, err := ctx.Con.LavaDevicesList()
	if err != nil {
		return err
	}

	if c.Yaml {
		d, err := yaml.Marshal(&ret)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
	} else if c.JSON {
		d, err := json.Marshal(&ret)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
	} else {
		fmt.Printf("Devices:\n")
		for _, v := range ret {
			fmt.Printf("* %s (%s): %s,%s\n", v.Hostname, v.Type, v.State, v.Health)
		}
	}
	return nil
}

type showDevicesCmd struct {
	DeviceName string `arg:"" required:"" help:"The device name to show."`
	Yaml       bool   `flag:"" optional:"" help:"Output as YAML" default:"false"`
	JSON       bool   `flag:"" optional:"" help:"Output as JSON" default:"false"`
}

func (c *showDevicesCmd) Run(ctx *context) error {

	ret, err := ctx.Con.LavaDevicesShow(c.DeviceName)
	if err != nil {
		return err
	}

	if c.Yaml {
		d, err := yaml.Marshal(&ret)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
	} else if c.JSON {
		d, err := json.Marshal(&ret)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
	} else {
		fmt.Printf("name       : %s\n", ret.Hostname)
		fmt.Printf("device-type: %s\n", ret.DeviceType)
		fmt.Printf("state      : %s\n", ret.State)
		fmt.Printf("health     : %s\n", ret.Health)
		fmt.Printf("health job : %v\n", ret.HealthJob)
		fmt.Printf("description: %s\n", ret.Description)
		fmt.Printf("deivce-dict: %v\n", ret.HasDeviceDict)
		fmt.Printf("worker     : %s\n", ret.Worker)
		fmt.Printf("current-job: %d\n", ret.CurrentJob)
		fmt.Printf("tags       : %v\n", ret.Tags)
	}
	return nil
}

type listDevicesTagCmd struct {
	DeviceName string `arg:"" required:"" help:"The device name to show."`
	Yaml       bool   `flag:"" optional:"" help:"Output as YAML" default:"false"`
	JSON       bool   `flag:"" optional:"" help:"Output as JSON" default:"false"`
}

func (c *listDevicesTagCmd) Run(ctx *context) error {
	ret, err := ctx.Con.LavaDevicesTagsList(c.DeviceName)
	if err != nil {
		return err
	}

	if c.Yaml {
		d, err := yaml.Marshal(&ret)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
	} else if c.JSON {
		d, err := json.Marshal(&ret)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
	} else {
		fmt.Printf("Tags:\n")
		for i := range ret {
			fmt.Printf("* %s\n", ret[i])
		}
	}
	return nil
}

type addDevicesTagCmd struct {
	DeviceName string `arg:"" required:"" help:"The device name to show."`
	Name       string `arg:"" required:"" help:"The tag to add to."`
}

func (c *addDevicesTagCmd) Run(ctx *context) error {
	return ctx.Con.LavaDevicesTagsAdd(c.DeviceName, c.Name)
}

type deleteDeviceTagCmd struct {
	DeviceName string `arg:"" required:"" help:"The device name to show."`
	Name       string `arg:"" required:"" help:"The tag to delete to."`
}

func (c *deleteDeviceTagCmd) Run(ctx *context) error {
	return ctx.Con.LavaDevicesTagsDelete(c.DeviceName, c.Name)
}

type devicesTagsCmd struct {
	List   listDevicesTagCmd  `cmd:"" help:"Lists tags"`
	Add    addDevicesTagCmd   `cmd:"" help:"Add a tag"`
	Delete deleteDeviceTagCmd `cmd:"" help:"Delete a tag"`
}

type devicesCmd struct {
	List listDevicesCmd `cmd:"" help:"Lists devices"`
	Tags devicesTagsCmd `cmd:"" help:"Handle device tags"`
	Show showDevicesCmd `cmd:"" help:"Show device properties"`
}
