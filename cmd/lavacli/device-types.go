// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type listDeviceTypesCmd struct {
	Yaml bool `flag:"" optional:"" help:"Output as YAML" default:"false"`
	JSON bool `flag:"" optional:"" help:"Output as JSON" default:"false"`
	All  bool `flag:"" optional:"" help:"Show all types" default:"false"`
}

func (c *listDeviceTypesCmd) Run(ctx *context) error {
	ret, err := ctx.LavaCon.DevicesTypesList(c.All)
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
		fmt.Printf("Device-Types:\n")
		for i := range ret {
			fmt.Printf("* %s (%d)\n", ret[i].Name, ret[i].Devices)
		}
	}

	return nil
}

type getTemplateCmd struct {
	Name string `arg:"" required:"" help:"Name of the template."`
}

func (c *getTemplateCmd) Run(ctx *context) error {

	ret, err := ctx.LavaCon.DevicesTypesTemplateGet(c.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s", ret)

	return nil
}

type setTemplateCmd struct {
	Name     string `arg:"" required:"" help:"Name of the template."`
	Filename string `arg:"" required:"" help:"Local filename of the template."`
}

func (c *setTemplateCmd) Run(ctx *context) error {
	templateFile, err := ioutil.ReadFile(c.Filename)
	if err != nil && !strings.HasSuffix(c.Filename, ".yaml") {
		templateFile, err = ioutil.ReadFile(c.Filename + ".yaml")
		if err != nil {
			return fmt.Errorf("Failed to read file: #%v ", err)
		}
	}

	err = ctx.LavaCon.DevicesTypesTemplateSet(c.Name, string(templateFile))

	return err
}

type getHealthCheckCmd struct {
	Name string `arg:"" required:"" help:"Name of the health-check."`
}

func (c *getHealthCheckCmd) Run(ctx *context) error {

	ret, err := ctx.LavaCon.DevicesTypesHealthCheckGet(c.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s", ret)

	return nil
}

type setHealthCheckCmd struct {
	Name     string `arg:"" required:"" help:"Name of the health-check."`
	Filename string `arg:"" required:"" help:"Local filename of the health-check."`
}

func (c *setHealthCheckCmd) Run(ctx *context) error {
	healthFile, err := ioutil.ReadFile(c.Filename)
	if err != nil && !strings.HasSuffix(c.Filename, ".yaml") {
		healthFile, err = ioutil.ReadFile(c.Filename + ".yaml")
		if err != nil {
			return fmt.Errorf("Failed to read file: #%v ", err)
		}
	}

	err = ctx.LavaCon.DevicesTypesHealthCheckSet(c.Name, string(healthFile))

	return err
}

type templateCmd struct {
	Get getTemplateCmd `cmd:"" help:"Get (download) a template from the server."`
	Set setTemplateCmd `cmd:"" help:"Set (upload) a template to the server."`
}

type healthCheckCmd struct {
	Get getHealthCheckCmd `cmd:"" help:"Get (download) a health-check from the server."`
	Set setHealthCheckCmd `cmd:"" help:"Set (upload) a health-check to the server."`
}

type deviceTypesCmd struct {
	List        listDeviceTypesCmd `cmd:"" help:"Lists device types"`
	Template    templateCmd        `cmd:"" help:"Handle device templates"`
	HealthCheck healthCheckCmd     `cmd:"" help:"Handle health checks"`
}
