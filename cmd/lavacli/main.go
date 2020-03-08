// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/siro20/lavacli/pkg/lava"
)

func CheckHelp(args []string) bool {
	for i := range args {
		if args[i] == "--help" {
			return true
		}
	}
	return false
}

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

func (g group) Exec(con *lava.LavaConnection, processedArgs []string, args []string) error {
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
	Exec(con *lava.LavaConnection, processedArgs []string, args []string) error
	ValidateArgs(processedArgs []string, args []string) bool
	Help(processedArgs []string, args []string) string
}

var lavacli group = group{
	map[string]command{
		"identities":   i,
		"devices":      d,
		"jobs":         j,
		"results":      r,
		"device-types": t,
	},
}

func main() {
	identity := flag.String("identity", "default", "identity stored in the configuration")
	uri := flag.String("uri", "", "URI of the lava-server RPC endpoint")
	proxy := flag.String("proxy", "", "Proxy to use when connecting")

	flag.Parse()

	if identity == nil {
		fmt.Fprintf(os.Stderr, "No identity specified\n")
		os.Exit(1)
	}
	if uri == nil {
		fmt.Fprintf(os.Stderr, "No URI specified\n")
		os.Exit(1)
	}
	if proxy == nil {
		fmt.Fprintf(os.Stderr, "No proxy specified\n")
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, lavacli.Help([]string{os.Args[0]}, nil)+"\n")
		os.Exit(1)
	}

	var c *lava.LavaConnection
	if os.Args[1] != "identities" {
		var err error
		if *uri != "" {
			c, err = lava.LavaConnectByUri(*uri, *proxy)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		} else {
			c, err = lava.LavaConnectByConfigID(*identity)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		}
	}

	err := lavacli.Exec(c, []string{os.Args[0]}, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
