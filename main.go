package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"alexejk.io/go-xmlrpc"
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

type command interface {
	Exec(uri string, processedArgs []string, args []string) error
	ValidateArgs(processedArgs []string, args []string) bool
	Help(processedArgs []string, args []string) string
}

type rect struct {
	width, height float64
}

func (r rect) Help(processedArgs []string, args []string) string {
	return ""
}

func (r rect) ValidateArgs(processedArgs []string, args []string) bool {
	return true
}

func (r rect) Exec(uri string, processedArgs []string, args []string) error {
	log.Println(uri)
	log.Println(processedArgs)
	log.Println(args)
	return nil
}

var r rect

var commands map[string]command = map[string]command{
	"aliases":      r,
	"devices":      r,
	"device-types": r,
	"events":       r,
	"identities":   i,
	"jobs":         r,
	"results":      r,
	"system":       r,
	"tags":         r,
	"utils":        r,
	"workers":      r,
}

func printHelp(processedArgs []string, args []string) {
	s := "Help:\n" + MakeHelp(commands, processedArgs, args, "[--identity IDENTITY] [--uri URI]")
	log.Fatal(s)
}

func main() {
	client, _ := xmlrpc.NewClient("https://bugzilla.mozilla.org/xmlrpc.cgi")

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
		printHelp([]string{os.Args[0]}, nil)
	}
	for k, v := range commands {
		if k == os.Args[1] {
			var c configIndentity
			var u string
			u = *uri
			if k != "identities" {
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
			}
			if !v.ValidateArgs([]string{os.Args[0], os.Args[1]}, os.Args[2:]) {
				log.Fatal(v.Help([]string{os.Args[0], os.Args[1]}, os.Args[2:]))
			}
			err := v.Exec(u, []string{os.Args[0], os.Args[1]}, os.Args[2:])
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	printHelp([]string{os.Args[0]}, nil)

	result := &struct {
		BugzillaVersion struct {
			Version string
		}
	}{}

	_ = client.Call("Bugzilla.version", nil, result)
	fmt.Printf("Version: %s\n", result.BugzillaVersion.Version)
}
