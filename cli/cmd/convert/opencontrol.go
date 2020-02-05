package convert

import (
	"github.com/docker/oscalkit/pkg/oc2oscal"
	"github.com/urfave/cli"
)

// ConvertOpenControl ...
var ConvertOpenControl = cli.Command{
	Name:        "opencontrol",
	Usage:       `convert from OpenControl format to OSCAL "implementation" format`,
	Description: `Convert OpenControl masonry repository into OSCAL directory`,
	ArgsUsage:   "[masonry-repository] [output-directory]",
	Before: func(c *cli.Context) error {
		if c.NArg() != 2 {
			return cli.NewExitError("Missing masonry repository or output directory", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		err := oc2oscal.Convert(c.Args()[0], c.Args()[1])
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	},
}
