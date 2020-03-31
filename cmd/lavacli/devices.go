// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/siro20/lavacli/pkg/lava"

	"gopkg.in/yaml.v2"
)

type devicesList struct {
}

func (l devicesList) GetParser() *flag.FlagSet {
	var yaml bool
	var json bool

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.BoolVar(&yaml, "yaml", false, "print as yaml")
	mySet.BoolVar(&json, "json", false, "print as json")

	return mySet
}

func (l devicesList) Help(processedArgs []string, args []string) string {
	mySet := l.GetParser()
	s := ""
	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	return MakeHelp(nil, processedArgs, args, s)
}

func (l devicesList) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) > 1 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := l.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) > 0 {
		return false
	}
	return true
}

func (l devicesList) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	ret, err := con.LavaDevicesList()
	if err != nil {
		return err
	}

	mySet := l.GetParser()
	mySet.Parse(args)

	isYaml := mySet.Lookup("yaml")
	isJson := mySet.Lookup("json")

	if isYaml != nil && isYaml.Value.String() == "true" {
		d, err := yaml.Marshal(&ret)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
	} else if isJson != nil && isJson.Value.String() == "true" {
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

type devicesShow struct {
}

func (d devicesShow) GetParser() *flag.FlagSet {
	var yaml bool
	var json bool

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.BoolVar(&yaml, "yaml", false, "print as yaml")
	mySet.BoolVar(&json, "json", false, "print as json")

	return mySet
}

func (d devicesShow) Help(processedArgs []string, args []string) string {
	mySet := d.GetParser()
	s := ""
	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	s += "<devicename> "
	return MakeHelp(nil, processedArgs, args, s)
}

func (d devicesShow) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) > 2 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := d.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return false
	}
	return true
}

func (d devicesShow) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := d.GetParser()
	mySet.Parse(args)

	name := mySet.Args()[0]

	ret, err := con.LavaDevicesShow(name)
	if err != nil {
		return err
	}

	isYaml := mySet.Lookup("yaml")
	isJson := mySet.Lookup("json")

	if isYaml != nil && isYaml.Value.String() == "true" {
		d, err := yaml.Marshal(&ret)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
	} else if isJson != nil && isJson.Value.String() == "true" {
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

type devicesTagsList struct {
}

func (l devicesTagsList) GetParser() *flag.FlagSet {
	var yaml bool
	var json bool

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.BoolVar(&yaml, "yaml", false, "print as yaml")
	mySet.BoolVar(&json, "json", false, "print as json")

	return mySet
}

func (l devicesTagsList) Help(processedArgs []string, args []string) string {
	mySet := l.GetParser()
	s := "hostname "
	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	return MakeHelp(nil, processedArgs, args, s)
}

func (l devicesTagsList) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) > 2 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := l.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) == 0 {
		return false
	}
	return true
}

func (l devicesTagsList) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := l.GetParser()
	mySet.Parse(args)

	name := mySet.Args()[0]

	ret, err := con.LavaDevicesTagsList(name)
	if err != nil {
		return err
	}

	isYaml := mySet.Lookup("yaml")
	isJson := mySet.Lookup("json")

	if isYaml != nil && isYaml.Value.String() == "true" {
		d, err := yaml.Marshal(&ret)
		if err != nil {
			return err
		}
		fmt.Println(string(d))
	} else if isJson != nil && isJson.Value.String() == "true" {
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

var tags group = group{
	map[string]command{
		"list": devicesTagsList{},
	},
}

var d group = group{
	map[string]command{
		"list": devicesList{},
		"show": devicesShow{},
		"tags": tags,
	},
}
