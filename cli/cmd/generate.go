package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	generator "github.com/opencontrol/oscalkit/generator"
	"github.com/opencontrol/oscalkit/templates"
	"github.com/opencontrol/oscalkit/types/oscal/catalog"
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

		bytes, err := ioutil.ReadFile(profilePath)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("cannot read profile. path: %s, err: %v ", profilePath, err), 1)
		}
		profile, err := generator.ReadProfile(bytes)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		newFile, err := os.Create("catalogs.go")
		if err != nil {
			return cli.NewExitError("cannot create file for catalogs", 1)
		}
		catalogs := generator.CreateCatalogsFromProfile(profile)
		t, err := templates.GetCatalogTemplate()
		if err != nil {
			return cli.NewExitError("cannot fetch template", 1)
		}
		err = t.Execute(newFile, struct {
			Catalogs []*catalog.Catalog
		}{catalogs})
		if err != nil {
			return cli.NewExitError("cannot write file for catalogs", 1)
		}
		logrus.Info("catalogs.go file created.")
		return nil

	},
}
