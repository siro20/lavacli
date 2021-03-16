// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type listJobsCmd struct {
	YAML   bool   `flag:"" optional:"" help:"Print as YAML" default:"false"`
	JSON   bool   `flag:"" optional:"" help:"Print as JSON" default:"false"`
	State  string `flag:"" optional:"" help:"[SUBMITTED, SCHEDULING, SCHEDULED, RUNNING, CANCELING, FINISHED]"`
	Health string `flag:"" optional:"" help:"[UNKNOWN, COMPLETE, INCOMPLETE, CANCELED]"`
	Start  int    `flag:"" optional:"" help:"Start at offset" default:"0"`
	Limit  int    `flag:"" optional:"" help:"Limit to #count jobs" default:"25"`
}

func (c *listJobsCmd) Run(ctx *context) error {

	ret, err := ctx.LavaCon.JobsList(c.State,
		c.Health, c.Start, c.Limit)
	if err != nil {
		return err
	}

	if c.YAML {
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
		fmt.Printf("Jobs (from %d to %d):\n", c.Start+1, c.Limit)
		for _, v := range ret {
			fmt.Printf("* %d: %s,%s [%s] (%s) - %s\n", v.ID, v.State, v.Health, v.Submitter, v.Description, v.DeviceType)
		}
	}
	return nil
}

type showJobCmd struct {
	YAML bool `flag:"" optional:"" help:"Print as YAML" default:"false"`
	JSON bool `flag:"" optional:"" help:"Print as JSON" default:"false"`
	ID   int  `arg:"" required:"" help:"Job ID"`
}

func (c *showJobCmd) Run(ctx *context) error {

	ret, err := ctx.LavaCon.JobsShow(c.ID)
	if err != nil {
		return err
	}

	if c.YAML {
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
		fmt.Printf("id          : %d\n", ret.ID)
		fmt.Printf("description : %s\n", ret.Description)
		fmt.Printf("submitter   : %s\n", ret.Submitter)
		fmt.Printf("device-type : %s\n", ret.DeviceType)
		fmt.Printf("device      : %v\n", ret.Device)
		fmt.Printf("health-check: %v\n", ret.HealthCheck)
		fmt.Printf("state       : %v\n", ret.State)
		fmt.Printf("health      : %s\n", ret.Health)
		fmt.Printf("pipeline    : %v\n", ret.Pipeline)
		fmt.Printf("tags        : %v\n", ret.Tags)
		fmt.Printf("visibility  : %v\n", ret.Visibility)
		fmt.Printf("submit time : %s\n", ret.SubmitTime)
		fmt.Printf("start time  : %s\n", ret.StartTime)
		fmt.Printf("end time    : %s\n", ret.EndTime)
	}

	return nil
}

type definitionJobCmd struct {
	ID int `arg:"" required:"" help:"Job ID"`
}

func (c *definitionJobCmd) Run(ctx *context) error {

	ret, err := ctx.LavaCon.JobsDefinition(c.ID)
	if err != nil {
		return err
	}

	fmt.Println(ret)
	return nil
}

type validateJobCmd struct {
	Filename string `arg:"" required:"" help:"File path to local job definition file"`
	Strict   bool   `flag:"" optional:"" help:"Strict mode"`
}

func (c *validateJobCmd) Run(ctx *context) error {
	path, err := filepath.Abs(c.Filename)
	if err != nil {
		return fmt.Errorf("Failed to resolv path: #%v ", err)
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read file: #%v ", err)
	}

	ret, err := ctx.LavaCon.JobsValidate(string(yamlFile), c.Strict)
	if err != nil {
		return err
	}

	for k, v := range ret {
		fmt.Printf("%s: %v\n", k, v)
	}

	return nil
}

type submitJobCmd struct {
	Filename string `arg:"" required:"" help:"File path to local job definition file"`
}

func (c *submitJobCmd) Run(ctx *context) error {
	path, err := filepath.Abs(c.Filename)
	if err != nil {
		return fmt.Errorf("Failed to resolv path: #%v ", err)
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read file: #%v ", err)
	}

	ret, err := ctx.LavaCon.JobsSubmitString(string(yamlFile))
	if err != nil {
		return err
	}

	fmt.Println(ret)
	return nil
}

type cancelJobCmd struct {
	ID int `arg:"" required:"" help:"Job ID"`
}

func (c *cancelJobCmd) Run(ctx *context) error {
	return ctx.LavaCon.JobsCancel(c.ID)
}

type logsJobCmd struct {
	ID  int  `arg:"" required:"" help:"Job ID"`
	Raw bool `flag:"" optional:"" help:"Print log in raw mode"`
}

func (c *logsJobCmd) Run(ctx *context) error {
	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Green = "\033[32m"
	var Yellow = "\033[33m"
	var Blue = "\033[34m"
	var Gray = "\033[37m"

	ret, err := ctx.LavaCon.JobsLogs(c.ID, c.Raw)
	if err != nil {
		return err
	}
	if c.Raw {
		fmt.Printf("%s\n", ret.Data)
	} else {
		for i := range ret.Decoded {
			lvl := ret.Decoded[i].Level
			colorStart := ""
			colorStop := ""
			if lvl == "feedback" {
				lvl = ""
				colorStart = Yellow
				colorStop = Reset
			} else if lvl == "target" {
				lvl = ""
				colorStart = Green
				colorStop = Reset
			} else if lvl == "results" {
				lvl = ""
				colorStart = Blue
				colorStop = Reset
			} else if lvl == "error" {
				lvl = ""
				colorStart = Red
				colorStop = Reset
			} else if lvl == "debug" {
				lvl = ""
				colorStart = Gray
				colorStop = Reset
			}
			fmt.Printf("%s: %s %s %s\n", ret.Decoded[i].DateTime,
				colorStart,
				ret.Decoded[i].Message,
				colorStop)
		}
	}

	//for k, v := range ret {
	//	fmt.Printf("%s: %v\n", k, v)
	//}

	return nil
}

type jobsCmd struct {
	List       listJobsCmd      `cmd:"" help:"Lists jobs"`
	Show       showJobCmd       `cmd:"" help:"Show job details"`
	Definition definitionJobCmd `cmd:"" help:"Handle job definition"`
	Validate   validateJobCmd   `cmd:"" help:"Validate job definition"`
	Submit     submitJobCmd     `cmd:"" help:"Submit new job"`
	Cancel     cancelJobCmd     `cmd:"" help:"Cancel running job"`
	Logs       logsJobCmd       `cmd:"" help:"Show job log"`
}
