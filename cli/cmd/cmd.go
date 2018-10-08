// oscalkit - OSCAL conversion utility
// Written in 2017 by Andrew Weiss <andrew.weiss@docker.com>

// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.

// You should have received a copy of the CC0 Public Domain Dedication along with this software.
// If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/opencontrol/oscalkit/cli/cmd/convert"
	"github.com/opencontrol/oscalkit/cli/version"
	"github.com/urfave/cli"
)

// Execute ...
func Execute() error {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("oscal-proxy version %s, build %s, date %s\n", c.App.Version, version.Version, version.Date)
	}

	app := cli.NewApp()
	app.Name = "oscalkit"
	app.Version = version.Version
	app.Usage = "OSCAL toolkit"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable debug command output",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}
	app.Commands = []cli.Command{
		convert.Convert,
		Validate,
		Sign,
	}

	return app.Run(os.Args)
}
