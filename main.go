package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kolo/xmlrpc"
)

func MakeHelp(cmds map[string]command, processedArgs []string, args []string, optarg string) string {
	s := "usage: "
	for _, k := range processedArgs {
		s += k + " "
	}
	if cmds != nil {
		s += ": {"
	}
	for k, _ := range cmds {
		s += k + ","
	}
	if strings.HasSuffix(s, ",") {
		s = s[:len(s)-1]
	}
	if cmds != nil {
		s += "} ..."
	}
	s += " " + optarg + " "
	s += "\n"
	return s
}

type group struct {
	Commands map[string]command
}

func (g group) Help(processedArgs []string, args []string) string {
	return MakeHelp(g.Commands, processedArgs, args, "")
}

func (g group) ValidateArgs(processedArgs []string, args []string) bool {
	if len(args) < 1 {
		return false
	}
	found := false
	for k, _ := range g.Commands {
		if args[0] == k {
			found = true
		}
	}
	return found
}

func (g group) Exec(con *xmlrpc.Client, processedArgs []string, args []string) error {
	for k, v := range g.Commands {
		if args[0] == k {
			if !v.ValidateArgs(append(processedArgs, args[0]), args[1:]) {
				return fmt.Errorf("%s", v.Help(append(processedArgs, args[0]), args[1:]))
			}
			return v.Exec(con, append(processedArgs, args[0]), args[1:])
		}
	}

	return fmt.Errorf("Internal error: Command not found\n%s", g.Help(append(processedArgs, args[0]), args[1:]))
}

type command interface {
	Exec(con *xmlrpc.Client, processedArgs []string, args []string) error
	ValidateArgs(processedArgs []string, args []string) bool
	Help(processedArgs []string, args []string) string
}

var lavacli group = group{
	map[string]command{
		"identities": i,
		"devices":    d,
		"jobs":       j,
	},
}

func getXMLRPCClient(uri string, proxy string) (*xmlrpc.Client, error) {
	tr := http.Transport{}

	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		tr.Proxy = http.ProxyURL(u)
	}

	client, err := xmlrpc.NewClient(uri, &tr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	identity := flag.String("identity", "default", "identity stored in the configuration")
	uri := flag.String("uri", "", "URI of the lava-server RPC endpoint")

	flag.Parse()

	if identity == nil {
		log.Fatal("No identity specified")
	}
	if uri == nil {
		log.Fatal("No URI specified")
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage:\n" + lavacli.Help([]string{os.Args[0]}, nil))
	}

	var c configIndentity
	var u string
	var con *xmlrpc.Client
	u = *uri
	if os.Args[1] != "identities" {
		var err error
		configs := GetConf()
		found := false
		for k, v := range configs {
			if k == *identity {
				c = v
				found = true
			}
		}
		if !found {
			log.Fatal("Identity not found in config")
		}

		if c.Uri == "" && *uri == "" {
			log.Fatal("No URI specified")
		}

		if c.Username != "" && c.Token != "" {
			url, err := url.Parse(c.Uri)
			if err != nil {
				log.Fatalf("Failed to parse URI: %v", err)
			}
			u = fmt.Sprintf("%s://%s:%s@%s%s", url.Scheme, c.Username, c.Token, url.Host, url.Path)
		} else if c.Uri != "" {
			u = c.Uri
		}
		con, err = getXMLRPCClient(u, c.Proxy)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := lavacli.Exec(con, []string{os.Args[0]}, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
