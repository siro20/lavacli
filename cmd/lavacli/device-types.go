// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/kolo/xmlrpc"
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

func (l devicesTypesList) Exec(con *xmlrpc.Client, processedArgs []string, args []string) error {

	mySet := l.GetParser()
	mySet.Parse(args)

	isYaml := mySet.Lookup("yaml")
	isJson := mySet.Lookup("json")
	showAll := mySet.Lookup("all")

	ret, err := lava.LavaDevicesTypesList(con, showAll != nil && showAll.Value.String() == "true")
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

var t group = group{
	map[string]command{
		"list": devicesTypesList{},
	},
}
