// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/siro20/lavacli/pkg/lava"
)

func connect(ctx *context) (c *lava.Connection, err error) {

	if ctx.URI != "" {
		c, err = lava.ConnectByURI(ctx.URI, ctx.Proxy, lava.DefaultOptions)
		if err != nil {
			err = fmt.Errorf("failed to connect by using URI %s: %v", ctx.URI, err)
			return
		}
	} else {
		c, err = lava.ConnectByConfigID(ctx.Profile, lava.DefaultOptions)
		if err != nil {
			err = fmt.Errorf("failed to connect by using identity %s: %v", ctx.Profile, err)
			return
		}
	}

	return
}

type context struct {
	Profile string
	URI     string
	Proxy   string
	LavaCon *lava.Connection
}

var cli struct {
	Profile string `help:"identity stored in the configuration." default:"default"`
	URI     string `help:"URI of the lava-server RPC endpoint. Default:Read from config."`
	Proxy   string `help:"Proxy to use when connecting. Default:Read from config."`

	Identities  identityCmd    `cmd:"" help:"Deals with identities in lavacli.yaml"`
	Devices     devicesCmd     `cmd:"" help:"Configure devices on the LAVA server."`
	Jobs        jobsCmd        `cmd:"" help:"Configure jobs on the LAVA server."`
	Results     resultsCmd     `cmd:"" help:"Get results on the LAVA server."`
	DeviceTypes deviceTypesCmd `cmd:"" help:"Configure device types on the LAVA server."`
}

func main() {
	var err error
	ctx := kong.Parse(&cli,
		kong.UsageOnError(),
		kong.Description("CLI for LAVA XMLRPC interface."),
		kong.ConfigureHelp(kong.HelpOptions{
			//Compact: true,
			Summary: true,
			Tree:    true,
		}))

	myCtx := context{Profile: cli.Profile,
		URI:   cli.URI,
		Proxy: cli.Proxy}

	if ctx.Command() != "identities" {
		myCtx.LavaCon, err = connect(&myCtx)
		ctx.FatalIfErrorf(err)
	}
	// Call the Run() method of the selected parsed command.
	err = ctx.Run(&myCtx)

	ctx.FatalIfErrorf(err)
}
