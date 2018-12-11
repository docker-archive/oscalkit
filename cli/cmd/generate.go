package cmd

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/opencontrol/oscalkit/generator"
	"github.com/opencontrol/oscalkit/templates"
	"github.com/opencontrol/oscalkit/types/oscal/catalog"

	"github.com/urfave/cli"
)

var profilePath string

//Generate Cli command to generate go code for controls
var Generate = cli.Command{
	Name:  "generate",
	Usage: "generates go code against provided profile",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "profile, p",
			Usage:       "profile to intersect against",
			Destination: &profilePath,
		},
	},
	Before: func(c *cli.Context) error {
		if profilePath == "" {
			return cli.NewExitError("oscalkit generate is missing the --profile flag", 1)
		}

		return nil
	},
	Action: func(c *cli.Context) error {

		outputFileName := "catalogs.go"
		profilePath, err := generator.GetAbsolutePath(profilePath)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("cannot get absolute path, err: %v", err), 1)
		}

		fileStat, err := os.Stat(profilePath)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("cannot fetch file, err %v", err), 1)
		}
		f, err := os.Open(fileStat.Name())
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		defer f.Close()

		profile, err := generator.ReadProfile(f)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		newFile, err := os.Create(outputFileName)
		if err != nil {
			return cli.NewExitError("cannot create file for catalogs", 1)
		}
		defer newFile.Close()

		catalogs := generator.CreateCatalogsFromProfile(profile)
		t, err := templates.GetCatalogTemplate()
		if err != nil {
			return cli.NewExitError("cannot fetch template", 1)
		}
		err = t.Execute(newFile, struct {
			Catalogs []*catalog.Catalog
		}{catalogs})

		//TODO: discuss better approach for formatting generate code file.
		if err != nil {
			return cli.NewExitError("cannot write file for catalogs", 1)
		}
		b, err := ioutil.ReadFile(outputFileName)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("cannot open %s file", outputFileName), 1)
		}
		b, err = format.Source(b)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("cannot format %s file", outputFileName), 1)
		}
		newFile.Write(b)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("cannot write formmated "), 1)
		}
		logrus.Info("catalogs.go file created.")
		return nil

	},
}
