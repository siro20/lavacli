package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type list struct {
	Name string
}

func (l list) Help() string {
	return MakeHelp(l.Name, nil, "")
}

func (l list) ValidateArgs() bool {
	if len(os.Args) != 3 {
		return false
	}
	return true
}

func (l list) Exec(uri string) error {
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

func (a add) Help() string {
	mySet := a.GetParser()
	s := ""
	mySet.VisitAll(func(f *flag.Flag) {
		s += "[--" + f.Name + " " + f.Usage + "] "
	})
	s += " <id>"
	return MakeHelp(a.Name, nil, s)
}

func (a add) ValidateArgs() bool {
	if len(os.Args) < 3 {
		return false
	}

	mySet := a.GetParser()
	mySet.Parse(os.Args[3:])

	if len(mySet.Args()) != 1 {
		return false
	}

	return true
}

func (a add) Exec(unused string) error {
	mySet := a.GetParser()
	mySet.Parse(os.Args[3:])

	configs := GetConf()
	if len(mySet.Args()) != 1 {
		log.Fatal(a.Help())
	}
	id := mySet.Args()[0]
	for k, _ := range configs {
		if k == id {
			log.Fatalf("id %s is already in config", id)
		}
	}

	uri := mySet.Lookup("uri")
	token := mySet.Lookup("token")
	username := mySet.Lookup("username")
	proxy := mySet.Lookup("proxy")

	if uri == nil {
		log.Fatal("Must specify URI")
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
	Name string
}

func (s show) Help() string {
	return MakeHelp(s.Name, nil, "<id>")
}

func (s show) ValidateArgs() bool {
	if len(os.Args) != 4 {
		return false
	}

	id := os.Args[3]

	configs := GetConf()
	for k, _ := range configs {
		if k == id {
			return true
		}
	}

	return false
}

func (s show) Exec(unused string) error {
	id := os.Args[3]

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

func (d del) Help() string {
	return MakeHelp(d.Name, nil, "<id>")
}

func (d del) ValidateArgs() bool {
	if len(os.Args) != 4 {
		return false
	}

	id := os.Args[3]

	configs := GetConf()
	for k, _ := range configs {
		if k == id {
			return true
		}
	}

	return false
}

func (d del) Exec(unused string) error {
	id := os.Args[3]

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

func (i identities) Help() string {
	return MakeHelp(i.Name, i.Commands, "")
}

func (i identities) ValidateArgs() bool {
	if len(os.Args) < 3 {
		return false
	}
	found := false
	for k, _ := range i.Commands {
		if os.Args[2] == k {
			found = true
		}
	}
	return found
}

func (i identities) Exec(unused string) error {
	for k, v := range i.Commands {
		if os.Args[2] == k {
			if !v.ValidateArgs() {
				log.Fatal(v.Help())
			}
			return v.Exec(unused)
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
		"show": show{
			"identities show",
		},
		"list": list{
			"identities list",
		},
	},
}
