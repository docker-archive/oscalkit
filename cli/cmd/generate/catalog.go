package generate

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/docker/oscalkit/generator"
	"github.com/urfave/cli"
)

var isJSON bool

// Catalog generates json/xml catalogs
var Catalog = cli.Command{
	Name:  "catalogs",
	Usage: "generates json/xml catalogs provided profile",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "profile, p",
			Usage:       "profile to intersect against",
			Destination: &profilePath,
		},
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "output filename",
			Destination: &outputFileName,
			Value:       "output",
		},
		cli.BoolFlag{

			Name:        "json, j",
			Usage:       "flag for generating catalogs in json",
			Destination: &isJSON,
		},
	},
	Before: func(c *cli.Context) error {
		if profilePath == "" {
			return cli.NewExitError("oscalkit generate is missing the --profile flag", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {

		profilePath, err := generator.GetAbsolutePath(profilePath)
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

		profile, err := generator.ReadProfile(f)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		catalogs, err := generator.CreateCatalogsFromProfile(profile)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("cannot create catalogs from profile, err: %v", err), 1)
		}

		var bytes []byte
		if !isJSON {
			bytes, err = xml.MarshalIndent(catalogs, "", "  ")
			if err != nil {
				return err
			}
			return ioutil.WriteFile(outputFileName+".xml", bytes, 0644)
		}
		bytes, err = json.MarshalIndent(catalogs, "", "  ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile(outputFileName+".json", bytes, 0644)

	},
	After: func(c *cli.Context) error {
		logrus.Info("catalog file generated")
		return nil
	},
}
