// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"

	"github.com/siro20/lavacli/pkg/lava"
)

type listIdentityCmd struct {
}

func (c *listIdentityCmd) Run(ctx *context) error {
	ids, err := lava.IdentitiesList()
	if err != nil {
		return err
	}
	fmt.Println("Identities:")
	for _, v := range ids {
		fmt.Printf("* %s\n", v.Name)
	}
	return nil
}

type addIdentityCmd struct {
	URI      string `arg:"" required:"" help:"The URI of the RPC XML interface."`
	Token    string `arg:"" required:"" help:"The authentication token of the user."`
	Username string `arg:"" required:"" help:"The user to authenticate with."`
	Proxy    string `arg:"" optional:"" help:"The proxy URI."`
}

func (c *addIdentityCmd) Run(ctx *context) error {

	var i lava.Indentity
	i.Name = c.Username
	i.URI = c.URI
	i.Token = c.Token
	i.Username = c.Username
	i.Proxy = c.Proxy

	return lava.IdentitiesAdd(i)
}

type showIdentityCmd struct {
	ID string `arg:"" required:"" help:"The identity to show."`
}

func (c *showIdentityCmd) Run(ctx *context) error {
	v, err := lava.IdentitiesShow(c.ID)
	if err != nil {
		return err
	}
	if v.Proxy != "" {
		fmt.Printf("proxy: %s\n", v.Proxy)
	}
	if v.Token != "" {
		fmt.Printf("token: %s\n", v.Token)
	}
	if v.URI != "" {
		fmt.Printf("uri: %s\n", v.URI)
	}
	if v.Username != "" {
		fmt.Printf("username: %s\n", v.Username)
	}
	return nil
}

type deleteIdentityCmd struct {
	ID string `arg:"" required:"" help:"The identity to delete."`
}

func (c *deleteIdentityCmd) Run(ctx *context) error {
	configs, err := lava.GetConf()

	if err != nil {
		return err
	}
	found := false
	for k := range configs {
		if k == c.ID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Identity not found")
	}
	return lava.IdentitiesDelete(c.ID)
}

type identityCmd struct {
	List   listIdentityCmd   `cmd:"" help:"Lists identities"`
	Add    addIdentityCmd    `cmd:"" help:"Add an identitiy"`
	Show   showIdentityCmd   `cmd:"" help:"Show an identitiy"`
	Delete deleteIdentityCmd `cmd:"" help:"Delete an identitiy"`
}
