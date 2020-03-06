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

	return fmt.Errorf("Internal error: Command not found")
}

type command interface {
	Exec(con *xmlrpc.Client, processedArgs []string, args []string) error
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

func (r rect) Exec(con *xmlrpc.Client, processedArgs []string, args []string) error {
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

func printHelp(processedArgs []string, args []string) {
	s := "Help:\n" + MakeHelp(commands, processedArgs, args, "[--identity IDENTITY] [--uri URI]")
	log.Fatal(s)
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
		printHelp([]string{os.Args[0]}, nil)
	}
	for k, v := range commands {
		if k == os.Args[1] {
			var c configIndentity
			var u string
			var con *xmlrpc.Client
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
				con, err := getXMLRPCClient(u, c.Proxy)
				if err != nil {
					log.Fatal(err)
				}

				type LavaJobState struct {
					Description    string `xmlrpc:"description"`
					State          string `xmlrpc:"state"`
					ID             int    `xmlrpc:"id"`
					EndTime        string `xmlrpc:"end_time"`
					SubmitTime     string `xmlrpc:"submit_time"`
					FailureComment string `xmlrpc:"failure_comment"`
					Status         int    `xmlrpc:"status"`
					HealthCheck    string `xmlrpc:"health_check"`
				}
				result := struct {
					Jobs []LavaJobState `xmlrpc:""`
				}{}
				var result2 interface{}

				err = con.Call("scheduler.all_devices", nil, &result2)
				log.Printf("%v\n", err)
				log.Printf("%v\n", result2)

				log.Printf("%v\n", result)
			}
			if !v.ValidateArgs([]string{os.Args[0], os.Args[1]}, os.Args[2:]) {
				log.Fatal(v.Help([]string{os.Args[0], os.Args[1]}, os.Args[2:]))
			}
			err := v.Exec(con, []string{os.Args[0], os.Args[1]}, os.Args[2:])
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	printHelp([]string{os.Args[0]}, nil)

}
