// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"flag"
	"fmt"
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
	return MakeHelp(g.Commands, processedArgs, args, "[--identity ID] [--uri URI]")
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
		fmt.Fprintf(os.Stderr, "No identity specified\n")
		os.Exit(1)
	}
	if uri == nil {
		fmt.Fprintf(os.Stderr, "No URI specified\n")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, lavacli.Help([]string{os.Args[0]}, nil)+"\n")
		os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "Identity %s not found in lavacli.yaml\n", *identity)
			os.Exit(1)
		}

		if c.Uri == "" && *uri == "" {
			fmt.Fprintf(os.Stderr, "No URI specified\n")
			os.Exit(1)
		}

		if c.Username != "" && c.Token != "" {
			url, err := url.Parse(c.Uri)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse URI: %v\n", err)
				os.Exit(1)
			}
			u = fmt.Sprintf("%s://%s:%s@%s%s", url.Scheme, c.Username, c.Token, url.Host, url.Path)
		} else if c.Uri != "" {
			u = c.Uri
		}
		con, err = getXMLRPCClient(u, c.Proxy)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	err := lavacli.Exec(con, []string{os.Args[0]}, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
