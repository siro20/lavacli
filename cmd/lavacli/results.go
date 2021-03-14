// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
)

type resultsShowCmd struct {
	ID   int  `arg:"" required:"" help:"Job ID"`
	Yaml bool `flag:"" optional:"" help:"Output as YAML" default:"false"`
	JSON bool `flag:"" optional:"" help:"Output as JSON" default:"false"`
}

func (c *resultsShowCmd) Run(ctx *context) error {

	if c.Yaml {
		ret, err := ctx.Con.LavaResultsAsYAML(c.ID)
		if err != nil {
			return err
		}
		fmt.Printf(ret)
	} else if c.JSON {
		ret, err := ctx.Con.LavaResultsAsJSON(c.ID)
		if err != nil {
			return err
		}
		fmt.Printf(ret)
	} else {
		ret, err := ctx.Con.LavaResults(c.ID)
		if err != nil {
			return err
		}
		for i := range ret {
			if len(ret[i].Name) == 0 {
				continue
			}
			if len(ret[i].Result) == 0 {
				continue
			}
			fmt.Printf("* %s [%s]\n", ret[i].Name, ret[i].Result)
		}
	}

	return nil
}

type resultsCmd struct {
	Show resultsShowCmd `cmd:"" help:"Show results of a finished job"`
}
