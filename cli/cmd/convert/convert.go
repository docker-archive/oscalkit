package convert

import (
	"github.com/urfave/cli"
)

var yaml bool

// Convert ...
var Convert = cli.Command{
	Name:  "convert",
	Usage: "convert between one or more OSCAL file formats and from OpenControl format",
	Subcommands: []cli.Command{
		ConvertOSCAL,
		ConvertHTML,
		ConvertOpenControl,
	},
}
