// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/siro20/lavacli/pkg/lava"
)

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

func (r resultsShow) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {

	mySet := r.GetParser()
	mySet.Parse(args)

	id, err := strconv.Atoi(mySet.Args()[0])
	if err != nil {
		return err
	}

	isYaml := mySet.Lookup("yaml")
	isJson := mySet.Lookup("json")

	if isYaml != nil && isYaml.Value.String() == "true" {
		ret, err := con.LavaResultsAsYAML(id)
		if err != nil {
			return err
		}
		fmt.Printf(ret)
	} else if isJson != nil && isJson.Value.String() == "true" {
		ret, err := con.LavaResultsAsJSON(id)
		if err != nil {
			return err
		}
		fmt.Printf(ret)
	} else {
		ret, err := con.LavaResults(id)
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
