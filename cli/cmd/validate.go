package cmd

import (
	"fmt"
	"github.com/docker/oscalkit/pkg/oscal_source"
	"github.com/urfave/cli"
)

var schemaFile string

// Validate ...
var Validate = cli.Command{
	Name:        "validate",
	Usage:       "validate files against OSCAL XML and JSON schemas",
	Description: `Validate OSCAL-formatted files against a specific OSCAL schema`,
	ArgsUsage:   "[files...]",
	Before: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.NewExitError("oscalkit validate requires at least one argument", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		for _, filePath := range c.Args() {
			os, err := oscal_source.Open(filePath)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Could not open oscal file: %v", err), 1)
			}
			defer os.Close()

			err = os.Validate()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
		}
		return nil
	},
}
