// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"flag"
	"fmt"

	"github.com/kolo/xmlrpc"
	"github.com/siro20/lavacli/pkg/lava"
)

type list struct {
}

func (l list) Help(processedArgs []string, args []string) string {
	return MakeHelp(nil, processedArgs, args, "")
}

func (l list) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 0 {
		return false
	}
	return true
}

func (l list) Exec(unused *xmlrpc.Client, processedArgs []string, args []string) error {
	ids, err := lava.LavaIdentitiesList()
	if err != nil {
		return err
	}
	fmt.Println("Identities:")
	for _, v := range ids {
		fmt.Printf("* %s\n", v.Name)
	}
	return nil
}

type add struct {
}

func (a add) GetParser() *flag.FlagSet {
	var uri string
	var token string
	var username string
	var proxy string

	mySet := flag.NewFlagSet("", flag.ExitOnError)
	mySet.StringVar(&uri, "uri", "", "URI")
	mySet.StringVar(&token, "token", "", "TOKEN")
	mySet.StringVar(&username, "username", "", "USERNAME")
	mySet.StringVar(&proxy, "proxy", "", "PROXY")

	return mySet
}

func (a add) Help(processedArgs []string, args []string) string {
	mySet := a.GetParser()
	s := ""
	mySet.VisitAll(func(f *flag.Flag) {
		if f.Name != "uri" {
			s += "[--" + f.Name + " " + f.Usage + "] "
		} else {
			s += "--" + f.Name + " " + f.Usage + " "
		}
	})
	s += " <id>"
	return MakeHelp(nil, processedArgs, args, s)
}

func (a add) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) < 3 {
		return false
	}
	mySet := a.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return false
	}

	return true
}

func (a add) Exec(unused *xmlrpc.Client, processedArgs []string, args []string) error {
	mySet := a.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return fmt.Errorf("%s", a.Help(processedArgs, args))
	}
	id := mySet.Args()[0]

	uri := mySet.Lookup("uri")
	token := mySet.Lookup("token")
	username := mySet.Lookup("username")
	proxy := mySet.Lookup("proxy")

	if uri == nil {
		return fmt.Errorf("Must specify URI")
	}

	var i lava.LavaIndentity
	i.Name = id
	i.Uri = uri.Value.String()
	if token != nil {
		i.Token = token.Value.String()
	}
	if username != nil {
		i.Username = username.Value.String()
	}
	if proxy != nil {
		i.Proxy = proxy.Value.String()
	}

	return lava.LavaIdentitiesAdd(i)
}

type show struct {
}

func (s show) Help(processedArgs []string, args []string) string {
	return MakeHelp(nil, processedArgs, args, "<id>")
}

func (s show) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 1 {
		return false
	}

	id := args[0]

	configs, err := lava.LavaGetConf()
	if err != nil {
		for k, _ := range configs {
			if k == id {
				return true
			}
		}
	}

	return false
}

func (s show) Exec(unused *xmlrpc.Client, processedArgs []string, args []string) error {
	id := args[0]

	v, err := lava.LavaIdentitiesShow(id)
	if err != nil {
		return err
	}
	if v.Proxy != "" {
		fmt.Printf("proxy: %s\n", v.Proxy)
	}
	if v.Token != "" {
		fmt.Printf("token: %s\n", v.Token)
	}
	if v.Uri != "" {
		fmt.Printf("uri: %s\n", v.Uri)
	}
	if v.Username != "" {
		fmt.Printf("username: %s\n", v.Username)
	}
	return nil
}

type del struct {
}

func (d del) Help(processedArgs []string, args []string) string {
	return MakeHelp(nil, processedArgs, args, "<id>")
}

func (d del) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 1 {
		return false
	}

	id := args[0]

	configs, err := lava.LavaGetConf()

	if err != nil {
		for k, _ := range configs {
			if k == id {
				return true
			}
		}
	}

	return false
}

func (d del) Exec(unused *xmlrpc.Client, processedArgs []string, args []string) error {
	id := args[0]

	return lava.LavaIdentitiesDelete(id)
}

var i group = group{
	map[string]command{
		"add":    add{},
		"delete": del{},
		"show":   show{},
		"list":   list{},
	},
}
