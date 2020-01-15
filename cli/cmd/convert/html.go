package convert

import (
	"fmt"
	"github.com/docker/oscalkit/pkg/oscal_source"
	"github.com/urfave/cli"
	"os"
)

// ConvertHTML ...
var ConvertHTML = cli.Command{
	Name:        "html",
	Usage:       "convert OSCAL file to human readable HTML",
	Description: `The command accepts source file and generates HTML representation of given file`,
	ArgsUsage:   "[source-file]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "output-path, o",
			Usage:       "Output path for converted file(s). Defaults to current working directory",
			Destination: &outputPath,
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() != 1 {
			// Check for stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				return nil
			}

			return cli.NewExitError("oscalkit convert html requires at one argument", 1)
		}

		return nil
	},
	Action: func(c *cli.Context) error {
		for _, sourcePath := range c.Args() {
			source, err := oscal_source.Open(sourcePath)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("could not load input file: %s", err), 1)
			}
			defer source.Close()

			buffer, err := source.HTML()
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("could convert to HTML: %s", err), 1)
			}
			if outputPath == "" {
				fmt.Println(buffer.String())
				return nil
			}

			f, err := os.Create(outputPath)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("could write to file: %s", err), 1)
			}
			_, err = f.WriteString(buffer.String())
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("could write to file: %s", err), 1)
			}
			err = f.Close()
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("failed to close file: %s", err), 1)
			}
		}
		return nil
	},
}
