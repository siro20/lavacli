// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/kolo/xmlrpc"
	"gopkg.in/yaml.v2"
)

type LavaJobsListing struct {
	Description string `xmlrpm:"description"`
	DeviceType  string `xmlrpc:"device_type"`
	Health      string `xmlrpc:"health"`
	ID          int    `xmlrpc:"id"`
	State       string `xmlrpc:"state"`
	Submitter   string `xmlrpc:"submitter"`
}

func LavaJobsList(con *xmlrpc.Client) ([]LavaJobsListing, error) {
	var ret []LavaJobsListing

	err := con.Call("scheduler.jobs.list", nil, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type LavaJobState struct {
	Description    string    `xmlrpc:"description"`
	DeviceType     string    `xmlrpc:"device_type"`
	Device         string    `xmlrpc:"device"`
	State          string    `xmlrpc:"state"`
	ID             int       `xmlrpc:"id"`
	EndTime        time.Time `xmlrpc:"end_time"`
	SubmitTime     time.Time `xmlrpc:"submit_time"`
	FailureComment string    `xmlrpc:"failure_comment"`
	Status         int       `xmlrpc:"status"`
	HealthCheck    bool      `xmlrpc:"health_check"`
	Pipeline       bool      `xmlrpc:"pipeline"`
	Tags           []string  `xmlrpc:"tags"`
	Visibility     string    `xmlrpc:"visibility"`
	Submitter      string    `xmlrpc:"submitter"`
	StartTime      time.Time `xmlrpc:"start_time"`
	Health         string    `xmlrpc:"health"`
}

func LavaJobsShow(con *xmlrpc.Client, id int) (*LavaJobState, error) {
	var ret LavaJobState

	err := con.Call("scheduler.jobs.show", id, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type LavaJobDefintion string

func LavaJobsDefinition(con *xmlrpc.Client, id int) (LavaJobDefintion, error) {
	var ret LavaJobDefintion

	err := con.Call("scheduler.jobs.definition", id, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// ******************

type jobsList struct {
}

func (l jobsList) GetParser() *flag.FlagSet {
	var yaml bool
	var json bool

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.BoolVar(&yaml, "yaml", false, "print as yaml")
	mySet.BoolVar(&json, "json", false, "print as json")

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
	if len(args) > 1 {
		return false
	}
	mySet := l.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) > 0 {
		return false
	}
	return true
}

func (l jobsList) Exec(con *xmlrpc.Client, processedArgs []string, args []string) error {

	ret, err := LavaJobsList(con)
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
		fmt.Printf("jobs:\n")
		for _, v := range ret {
			fmt.Printf("* %d %s,%s [%s] (%s) - %s\n", v.ID, v.State, v.Health, v.Submitter, v.Description, v.DeviceType)
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
	mySet := j.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return false
	}
	return true
}

func (j jobsShow) Exec(con *xmlrpc.Client, processedArgs []string, args []string) error {

	mySet := j.GetParser()
	mySet.Parse(args)

	id, err := strconv.Atoi(mySet.Args()[0])
	if err != nil {
		return err
	}
	ret, err := LavaJobsShow(con, id)
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

	return true
}

func (j jobsDefinition) Exec(con *xmlrpc.Client, processedArgs []string, args []string) error {

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	ret, err := LavaJobsDefinition(con, id)
	if err != nil {
		return err
	}

	fmt.Println(ret)

	return nil
}

var j group = group{
	map[string]command{
		"list":       jobsList{},
		"show":       jobsShow{},
		"definition": jobsDefinition{},
	},
}
