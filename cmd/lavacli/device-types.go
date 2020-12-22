// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/siro20/lavacli/pkg/lava"
	"gopkg.in/yaml.v2"
)

type devicesTypesList struct {
}

func (l devicesTypesList) GetParser() *flag.FlagSet {
	var yaml bool
	var json bool
	var all bool

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.BoolVar(&yaml, "yaml", false, "print as yaml")
	mySet.BoolVar(&json, "json", false, "print as json")
	mySet.BoolVar(&all, "all", false, "Show all types")

	return mySet
}

func (l devicesTypesList) Help(processedArgs []string, args []string) string {
	mySet := l.GetParser()
	s := ""
	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	return MakeHelp(nil, processedArgs, args, s)
}

func (l devicesTypesList) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) > 2 {
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

func (l devicesTypesList) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := l.GetParser()
	mySet.Parse(args)

	isYaml := mySet.Lookup("yaml")
	isJson := mySet.Lookup("json")
	showAll := mySet.Lookup("all")

	ret, err := con.LavaDevicesTypesList(showAll != nil && showAll.Value.String() == "true")
	if err != nil {
		return err
	}

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
		fmt.Printf("Device-Types:\n")
		for i := range ret {
			fmt.Printf("* %s (%d)\n", ret[i].Name, ret[i].Devices)
		}
	}

	return nil
}

type devicesTypesGet struct {
}

func (l devicesTypesGet) GetParser() *flag.FlagSet {

	mySet := flag.NewFlagSet("", flag.ExitOnError)

	return mySet
}

func (l devicesTypesGet) Help(processedArgs []string, args []string) string {
	mySet := l.GetParser()
	s := "<name> "

	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	return MakeHelp(nil, processedArgs, args, s)
}

func (l devicesTypesGet) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 1 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := l.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return false
	}
	return true
}

func (l devicesTypesGet) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := l.GetParser()
	mySet.Parse(args)

	name := mySet.Args()[0]

	ret, err := con.LavaDevicesTypesTemplateGet(name)
	if err != nil {
		return err
	}

	fmt.Printf("%s", ret)

	return nil
}

type devicesTypesSet struct {
}

func (l devicesTypesSet) GetParser() *flag.FlagSet {

	mySet := flag.NewFlagSet("", flag.ExitOnError)

	return mySet
}

func (l devicesTypesSet) Help(processedArgs []string, args []string) string {
	mySet := l.GetParser()
	s := "<name> <template>"

	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	return MakeHelp(nil, processedArgs, args, s)
}

func (l devicesTypesSet) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 2 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := l.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 2 {
		return false
	}
	return true
}

func (l devicesTypesSet) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := l.GetParser()
	mySet.Parse(args)

	name := mySet.Args()[0]
	template := mySet.Args()[1]

	templateFile, err := ioutil.ReadFile(template)
	if err != nil {
		return fmt.Errorf("Failed to read file: #%v ", err)
	}

	err = con.LavaDevicesTypesTemplateSet(name, string(templateFile))

	return err
}

//

type healthCheckGet struct {
}

func (l healthCheckGet) GetParser() *flag.FlagSet {

	mySet := flag.NewFlagSet("", flag.ExitOnError)

	return mySet
}

func (l healthCheckGet) Help(processedArgs []string, args []string) string {
	mySet := l.GetParser()
	s := "<name> "

	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	return MakeHelp(nil, processedArgs, args, s)
}

func (l healthCheckGet) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 1 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := l.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return false
	}
	return true
}

func (l healthCheckGet) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := l.GetParser()
	mySet.Parse(args)

	name := mySet.Args()[0]

	ret, err := con.LavaDevicesTypesHealthCheckGet(name)
	if err != nil {
		return err
	}

	fmt.Printf("%s", ret)

	return nil
}

type healthCheckSet struct {
}

func (l healthCheckSet) GetParser() *flag.FlagSet {

	mySet := flag.NewFlagSet("", flag.ExitOnError)

	return mySet
}

func (l healthCheckSet) Help(processedArgs []string, args []string) string {
	mySet := l.GetParser()
	s := "<name> <healthcheck>"

	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	return MakeHelp(nil, processedArgs, args, s)
}

func (l healthCheckSet) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 2 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := l.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 2 {
		return false
	}
	return true
}

func (l healthCheckSet) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := l.GetParser()
	mySet.Parse(args)

	name := mySet.Args()[0]
	template := mySet.Args()[1]

	templateFile, err := ioutil.ReadFile(template)
	if err != nil {
		return fmt.Errorf("Failed to read file: #%v ", err)
	}

	err = con.LavaDevicesTypesHealthCheckSet(name, string(templateFile))

	return err
}

var templates group = group{
	map[string]command{
		"get": devicesTypesGet{},
		"set": devicesTypesSet{},
	},
}

var healthCheck group = group{
	map[string]command{
		"get": healthCheckGet{},
		"set": healthCheckSet{},
	},
}

var t group = group{
	map[string]command{
		"list":         devicesTypesList{},
		"template":     templates,
		"health-check": healthCheck,
	},
}
