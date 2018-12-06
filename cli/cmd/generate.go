package cmd

import (
	"fmt"
	"os"

	generator "github.com/opencontrol/oscalkit/generator"
	"github.com/urfave/cli"
)

var profilePath string

//Generate Cli command to generate go code for controls
var Generate = cli.Command{
	Name:  "generate",
	Usage: "generates go code against provided catalogs and profile",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "profile, p",
			Usage:       "profile to intersect against",
			Destination: &profilePath,
		},
	},
	Before: func(c *cli.Context) error {
		if profilePath == "" {
			return cli.NewExitError("oscalkit sign is missing the --profile flag", 1)
		}

		return nil
	},
	Action: func(c *cli.Context) error {
		f, err := os.Open(profilePath)
		if err != nil {
			s := fmt.Sprintf("cannot open profile. path: %s ", err)
			return cli.NewExitError(s, 1)
		}
		profile, err := generator.ReadProfile(f)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		newFile, err := os.Create("catalogs.go")
		if err != nil {
			return cli.NewExitError("cannot create file for catalogs", 1)
		}
		catalogs := generator.IntersectProfile(profile)
		err = generator.GenerateCatalogs(newFile, catalogs)
		if err != nil {
			return cli.NewExitError("cannot write file for catalogs", 1)
		}
		return nil

	},
}
