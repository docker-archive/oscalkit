package cmd

import (
	"fmt"
	"os"

	"github.com/docker/oscalkit/generator"
	"github.com/docker/oscalkit/types/oscal"
	"github.com/urfave/cli"
)

// Catalog generates json/xml catalogs
var Info = cli.Command{
	Name:      "info",
	Usage:     "Provides information about particular OSCAL resource",
	ArgsUsage: "[file]",
	Action: func(c *cli.Context) error {
		for _, filePath := range c.Args() {
			profilePath, err := generator.GetAbsolutePath(filePath)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("cannot get absolute path, err: %v", err), 1)
			}

			_, err = os.Stat(profilePath)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("cannot fetch file, err %v", err), 1)
			}
			f, err := os.Open(profilePath)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			defer f.Close()

			o, err := oscal.New(f)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			if o.Profile != nil {
				fmt.Println("OSCAL Profile")
				return nil
			}
			if o.Catalog != nil {
				fmt.Println("OSCAL Catalog")
				return nil
			}
			return cli.NewExitError("Unrecognized OSCAL resource", 1)
		}
		return cli.NewExitError("No file provided", 1)
	},
}
