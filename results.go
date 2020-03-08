// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"

	"github.com/kolo/xmlrpc"
	"gopkg.in/yaml.v2"
)

type LavaResult []struct {
	Name         string `yaml:"name,omitempty" json:"name,omitempty"`
	Result       string `yaml:"result,omitempty" json:"result,omitempty"`
	Id           string `yaml:"id,omitempty" json:"id,omitempty"`
	Job          string `yaml:"job,omitempty" json:"job,omitempty"`
	Level        string `yaml:"level,omitempty" json:"level,omitempty"`
	LogLineEnd   string `yaml:"log_end_line,omitempty" json:"log_end_line,omitempty"`
	LogLineStart string `yaml:"log_start_line,omitempty" json:"log_start_line,omitempty"`
	Suite        string `yaml:"suite,omitempty" json:"suite,omitempty"`
	URL          string `yaml:"url,omitempty" json:"url,omitempty"`
	Measurement  string `yaml:"measurement,omitempty" json:"measurement,omitempty"`
	Logged       string `yaml:"logged,omitempty" json:"logged,omitempty"`
	Metadata     struct {
		Case       string `yaml:"case,omitempty" json:"case,omitempty"`
		Definition string `yaml:"definition,omitempty" json:"definition,omitempty"`
		Duration   string `yaml:"duration,omitempty" json:"duration,omitempty"`
		// This is broken as it should contain a slice of maps with a single key,
		// but it can be anything...
		//Extra      []map[string]interface{} `yaml:"extra,omitempty"`
		Level      string `yaml:"level,omitempty" json:"level,omitempty"`
		Namespace  string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
		Result     string `yaml:"result,omitempty" json:"result,omitempty"`
		UUID       string `yaml:"uuid,omitempty" json:"uuid,omitempty"`
		Revision   string `yaml:"revision,omitempty" json:"revision,omitempty"`
		Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
		CommitID   string `yaml:"commit_id,omitempty" json:"commit_id,omitempty"`
		Path       string `yaml:"path,omitempty" json:"path,omitempty"`
	} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

func LavaResultsAsYAML(con *xmlrpc.Client, id int) (string, error) {
	var ret string

	err := con.Call("results.get_testjob_results_yaml", id, &ret)

	return ret, err
}

func LavaResults(con *xmlrpc.Client, id int) (LavaResult, error) {
	var ret LavaResult
	yamlStr, err := LavaResultsAsYAML(con, id)

	err = yaml.Unmarshal([]byte(yamlStr), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func LavaResultsAsJSON(con *xmlrpc.Client, id int) (string, error) {
	results, err := LavaResults(con, id)

	d, err := json.Marshal(&results)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

type resultsShow struct {
}

func (r resultsShow) GetParser() *flag.FlagSet {
	var yaml bool
	var json bool

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.BoolVar(&yaml, "yaml", false, "print as yaml")
	mySet.BoolVar(&json, "json", false, "print as json")

	return mySet
}

func (r resultsShow) Help(processedArgs []string, args []string) string {
	mySet := r.GetParser()
	s := ""
	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	s += "<job_id> "
	return MakeHelp(nil, processedArgs, args, s)
}

func (r resultsShow) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) > 2 {
		return false
	}
	if CheckHelp(args) {
		return false
	}
	mySet := r.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return false
	}
	return true
}

func (r resultsShow) Exec(con *xmlrpc.Client, processedArgs []string, args []string) error {

	mySet := r.GetParser()
	mySet.Parse(args)

	id, err := strconv.Atoi(mySet.Args()[0])
	if err != nil {
		return err
	}

	isYaml := mySet.Lookup("yaml")
	isJson := mySet.Lookup("json")

	if isYaml != nil && isYaml.Value.String() == "true" {
		ret, err := LavaResultsAsYAML(con, id)
		if err != nil {
			return err
		}
		fmt.Printf(ret)
	} else if isJson != nil && isJson.Value.String() == "true" {
		ret, err := LavaResultsAsJSON(con, id)
		if err != nil {
			return err
		}
		fmt.Printf(ret)
	} else {
		ret, err := LavaResults(con, id)
		if err != nil {
			return err
		}
		for i := range ret {
			if len(ret[i].Name) == 0 {
				continue
			}
			if len(ret[i].Result) == 0 {
				continue
			}
			fmt.Printf("* %s [%s]\n", ret[i].Name, ret[i].Result)
		}

	}

	return nil
}

var r resultsShow
