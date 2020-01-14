package cmd

import (
	"fmt"
	"os"

	"github.com/docker/oscalkit/generator"
	"github.com/docker/oscalkit/types/oscal"
	"github.com/docker/oscalkit/types/oscal/catalog"
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
				fmt.Println("OSCAL Profile (represents subset of controls from OSCAL catalog(s))")
				fmt.Println("ID:\t", o.Profile.Id)
				printMetadata(o.Profile.Metadata)
				return nil
			}
			if o.Catalog != nil {
				fmt.Println("OSCAL Catalog (represents library of control assessment objectives and activities)")
				fmt.Println("ID:\t", o.Catalog.Id)
				printMetadata(o.Catalog.Metadata)
				return nil
			}
			return cli.NewExitError("Unrecognized OSCAL resource", 1)
		}
		return cli.NewExitError("No file provided", 1)
	},
}

func printMetadata(m *catalog.Metadata) {
	if m == nil {
		return
	}
	fmt.Println("Metadata:")
	fmt.Println("\tTitle:\t\t\t", m.Title)
	if m.Published != "" {
		fmt.Println("\tPublished:\t\t", m.Published)
	}
	if m.LastModified != "" {
		fmt.Println("\tLast Modified:\t\t", m.LastModified)
	}
	if m.Version != "" {
		fmt.Println("\tDocument Version:\t", m.Version)
	}
	if m.OscalVersion != "" {
		fmt.Println("\tOSCAL Version:\t\t", m.OscalVersion)
	}
}
