package generate

import (
	"github.com/urfave/cli"
)

// Generate Cli command to generate go code for controls
var Generate = cli.Command{
	Name:  "generate",
	Usage: "generates catalogs code/xml/json against provided profile",
	Subcommands: []cli.Command{
		Catalog,
		Code,
		Implementation,
	},
}
