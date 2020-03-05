package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"alexejk.io/go-xmlrpc"
	"gopkg.in/yaml.v2"
)

func MakeHelp(name string, cmds map[string]command, optarg string) string {
	s := "usage: " + os.Args[0] + " " + name
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
	Exec(string) error
	ValidateArgs() bool
	Help() string
}

type rect struct {
	width, height float64
}

func (r rect) Help() string {
	return ""
}

func (r rect) ValidateArgs() bool {
	return true
}
func (r rect) Exec(uri string) error {
	log.Println(uri)
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

func printHelp() {
	s := "Help:\n" + MakeHelp("", commands, "")
	log.Fatal(s)
}

type configIndentity struct {
	Token    string `yaml:"token,omitempty"`
	Uri      string `yaml:"uri,omitempty"`
	Username string `yaml:"username,omitempty"`
	Proxy    string `yaml:"proxy,omitempty"`
}

func GetConf() map[string]configIndentity {
	var c map[string]configIndentity

	path := os.Getenv("XDG_CONFIG_HOME")
	if path == "" {
		usr, _ := user.Current()
		path = usr.HomeDir + "/.config"
	}

	path += "/lavacli.yaml"
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to resolv path: #%v ", err)
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}

	return c
}

func SetConf(c map[string]configIndentity) {
	path := os.Getenv("XDG_CONFIG_HOME")
	if path == "" {
		usr, _ := user.Current()
		path = usr.HomeDir + "/.config"
	}

	path += "/lavacli.yaml"
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to resolv path: #%v ", err)
	}

	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
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
		printHelp()
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
			if !v.ValidateArgs() {
				log.Fatal(v.Help())
			}
			v.Exec(u)
			return
		}
	}
	printHelp()

	result := &struct {
		BugzillaVersion struct {
			Version string
		}
	}{}

	_ = client.Call("Bugzilla.version", nil, result)
	fmt.Printf("Version: %s\n", result.BugzillaVersion.Version)
}
