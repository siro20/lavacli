package main

import (
	"flag"
	"fmt"
)

type list struct {
	Name string
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

func (l list) Exec(unused string, processedArgs []string, args []string) error {
	configs := GetConf()
	fmt.Println("Identities:")
	for k, _ := range configs {
		fmt.Printf("* %s\n", k)
	}
	return nil
}

type add struct {
	Name string
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
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	s += " <id>"
	return MakeHelp(nil, processedArgs, args, s)
}

func (a add) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) < 4 {
		return false
	}

	mySet := a.GetParser()
	mySet.Parse(args)

	if len(mySet.Args()) != 1 {
		return false
	}

	return true
}

func (a add) Exec(unused string, processedArgs []string, args []string) error {
	mySet := a.GetParser()
	mySet.Parse(args)

	configs := GetConf()
	if len(mySet.Args()) != 1 {
		return fmt.Errorf("%s", a.Help(processedArgs, args))
	}
	id := mySet.Args()[0]
	for k, _ := range configs {
		if k == id {
			return fmt.Errorf("id %s is already in config", id)
		}
	}

	uri := mySet.Lookup("uri")
	token := mySet.Lookup("token")
	username := mySet.Lookup("username")
	proxy := mySet.Lookup("proxy")

	if uri == nil {
		return fmt.Errorf("Must specify URI")
	}

	var c configIndentity
	c.Uri = uri.Value.String()
	if token != nil {
		c.Token = token.Value.String()
	}
	if username != nil {
		c.Username = username.Value.String()
	}
	if proxy != nil {
		c.Proxy = proxy.Value.String()
	}
	configs[id] = c

	SetConf(configs)

	return nil
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

	configs := GetConf()
	for k, _ := range configs {
		if k == id {
			return true
		}
	}

	return false
}

func (s show) Exec(unused string, processedArgs []string, args []string) error {
	id := args[0]

	configs := GetConf()

	for k, v := range configs {
		if k == id {
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
	}

	return fmt.Errorf("Internal error, identity not found")
}

type del struct {
	Name string
}

func (d del) Help(processedArgs []string, args []string) string {
	return MakeHelp(nil, processedArgs, args, "<id>")
}

func (d del) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) != 1 {
		return false
	}

	id := args[0]

	configs := GetConf()
	for k, _ := range configs {
		if k == id {
			return true
		}
	}

	return false
}

func (d del) Exec(unused string, processedArgs []string, args []string) error {
	id := args[0]

	configs := GetConf()

	new := map[string]configIndentity{}
	for k, v := range configs {
		if k != id {
			new[k] = v
		}
	}

	SetConf(new)

	return nil
}

type identities struct {
	Name     string
	Commands map[string]command
}

func (i identities) Help(processedArgs []string, args []string) string {
	return MakeHelp(i.Commands, processedArgs, args, "")
}

func (i identities) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) < 1 {
		return false
	}
	found := false
	for k, _ := range i.Commands {
		if args[0] == k {
			found = true
		}
	}
	return found
}

func (i identities) Exec(unused string, processedArgs []string, args []string) error {
	for k, v := range i.Commands {
		if args[0] == k {
			if !v.ValidateArgs(append(processedArgs, args[0]), args[1:]) {
				return fmt.Errorf("%s", v.Help(append(processedArgs, args[0]), args[1:]))
			}
			return v.Exec(unused, append(processedArgs, args[0]), args[1:])
		}
	}

	return fmt.Errorf("Internal error: Command not found")
}

var i identities = identities{
	"identities",
	map[string]command{
		"add": add{
			"identities add",
		},
		"delete": del{
			"identities delete",
		},
		"show": show{},
		"list": list{
			"identities list",
		},
	},
}
