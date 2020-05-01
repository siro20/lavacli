// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/siro20/lavacli/pkg/lava"

	"gopkg.in/yaml.v2"
)

type jobsList struct {
}

func (l jobsList) GetParser() *flag.FlagSet {
	var yaml bool
	var json bool
	var state string
	var health string
	var start int
	var limit int

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.BoolVar(&yaml, "yaml", false, "print as yaml")
	mySet.BoolVar(&json, "json", false, "print as json")
	mySet.StringVar(&state, "state", "", "[SUBMITTED, SCHEDULING, SCHEDULED, RUNNING, CANCELING, FINISHED]")
	mySet.StringVar(&health, "health", "", "[UNKNOWN, COMPLETE, INCOMPLETE, CANCELED]")
	mySet.IntVar(&start, "start", 0, "Start at offset")
	mySet.IntVar(&limit, "limit", 25, "Limit to #count jobs")

	return mySet
}

func (l jobsList) Help(processedArgs []string, args []string) string {
	mySet := l.GetParser()
	s := ""
	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	return MakeHelp(nil, processedArgs, args, s)
}

func (l jobsList) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) > 5 {
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

func (l jobsList) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := l.GetParser()
	mySet.Parse(args)

	state := mySet.Lookup("state")
	health := mySet.Lookup("health")

	start, err := strconv.Atoi(mySet.Lookup("start").Value.String())
	if err != nil {
		return err
	}

	limit, err := strconv.Atoi(mySet.Lookup("limit").Value.String())
	if err != nil {
		return err
	}

	ret, err := con.LavaJobsList(state.Value.String(),
		health.Value.String(), start, limit)
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
		fmt.Printf("Jobs (from %d to %d):\n", start+1, limit)
		for _, v := range ret {
			fmt.Printf("* %d: %s,%s [%s] (%s) - %s\n", v.ID, v.State, v.Health, v.Submitter, v.Description, v.DeviceType)
		}
	}

	return nil
}

type jobsShow struct {
}

func (j jobsShow) GetParser() *flag.FlagSet {
	var yaml bool
	var json bool

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.BoolVar(&yaml, "yaml", false, "print as yaml")
	mySet.BoolVar(&json, "json", false, "print as json")

	return mySet
}

func (j jobsShow) Help(processedArgs []string, args []string) string {
	mySet := j.GetParser()
	s := ""
	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	s += "<devicename> "
	return MakeHelp(nil, processedArgs, args, s)
}

func (j jobsShow) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) > 2 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := j.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return false
	}
	return true
}

func (j jobsShow) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := j.GetParser()
	mySet.Parse(args)

	id, err := strconv.Atoi(mySet.Args()[0])
	if err != nil {
		return err
	}
	ret, err := con.LavaJobsShow(id)
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

type jobsDefinition struct {
}

func (j jobsDefinition) Help(processedArgs []string, args []string) string {
	return MakeHelp(nil, processedArgs, args, "<id>")
}

func (j jobsDefinition) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 1 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	return true
}

func (j jobsDefinition) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	ret, err := con.LavaJobsDefinition(id)
	if err != nil {
		return err
	}

	fmt.Println(ret)

	return nil
}

type jobsValidate struct {
}

func (j jobsValidate) GetParser() *flag.FlagSet {
	var strict bool

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.BoolVar(&strict, "strict", false, "strict mode")

	return mySet
}

func (j jobsValidate) Help(processedArgs []string, args []string) string {
	mySet := j.GetParser()
	s := ""
	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	s += "<definition file> "
	return MakeHelp(nil, processedArgs, args, s)
}

func (j jobsValidate) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 1 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := j.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return false
	}
	return true
}

func (j jobsValidate) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := j.GetParser()
	mySet.Parse(args)

	isStrict := mySet.Lookup("strict")

	path, err := filepath.Abs(mySet.Args()[0])
	if err != nil {
		return fmt.Errorf("Failed to resolv path: #%v ", err)
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read file: #%v ", err)
	}

	ret, err := con.LavaJobsValidate(string(yamlFile),
		isStrict != nil && isStrict.Value.String() == "true")
	if err != nil {
		return err
	}

	for k, v := range ret {
		fmt.Printf("%s: %v\n", k, v)
	}

	return nil
}

type jobsSubmit struct {
}

func (j jobsSubmit) Help(processedArgs []string, args []string) string {
	return MakeHelp(nil, processedArgs, args, "<definition>")
}

func (j jobsSubmit) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 1 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	return true
}

func (j jobsSubmit) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	path, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("Failed to resolv path: #%v ", err)
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read file: #%v ", err)
	}

	ret, err := con.LavaJobsSubmit(string(yamlFile))
	if err != nil {
		return err
	}

	fmt.Println(ret)

	return nil
}

type jobsCancel struct {
}

func (j jobsCancel) Help(processedArgs []string, args []string) string {
	return MakeHelp(nil, processedArgs, args, "<id>")
}

func (j jobsCancel) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 1 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	return true
}

func (j jobsCancel) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	err = con.LavaJobsCancel(id)

	return err
}

var j group = group{
	map[string]command{
		"list":       jobsList{},
		"show":       jobsShow{},
		"definition": jobsDefinition{},
		"validate":   jobsValidate{},
		"submit":     jobsSubmit{},
		"cancel":     jobsCancel{},
	},
}
