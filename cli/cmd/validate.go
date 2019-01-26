package cmd

import (
	"path/filepath"

	"github.com/docker/oscalkit/validator"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var schemaFile string

// Validate ...
var Validate = cli.Command{
	Name:  "validate",
	Usage: "validate files against OSCAL XML and JSON schemas",
	Description: `Validate OSCAL-formatted XML files against a specific XML schema (.xsd)
	 or OSCAL-formatted JSON files against a specific JSON schema`,
	ArgsUsage: "[files...]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "schema, s",
			Usage:       "schema file to validate against",
			Destination: &schemaFile,
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.NewExitError("oscalkit validate requires at least one argument", 1)
		}

		if schemaFile == "" {
			return cli.NewExitError("missing schema file (-s) flag", 1)
		}

		for _, f := range c.Args() {
			if filepath.Ext(f) == ".xml" && filepath.Ext(schemaFile) != ".xsd" {
				return cli.NewExitError("Schema file should be .xsd", 1)
			}

			if filepath.Ext(f) == ".json" && filepath.Ext(schemaFile) != ".json" {
				return cli.NewExitError("Schema file should be .json", 1)
			}
		}

		return nil
	},
	Action: func(c *cli.Context) error {
		schemaValidator := validator.New(schemaFile)

		if err := schemaValidator.Validate(c.Args()...); err != nil {
			logrus.Error(err)
			return nil
		}

		logrus.Debug("Validation complete")

		return nil
	},
}
