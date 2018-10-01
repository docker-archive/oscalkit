// oscalkit - OSCAL conversion utility
// Written in 2017 by Andrew Weiss <andrew.weiss@docker.com>

// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.

// You should have received a copy of the CC0 Public Domain Dedication along with this software.
// If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

package convert

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/opencontrol/oscalkit/types/oscal"
	"github.com/urfave/cli"
)

var outputPath string
var outputFile string

// ConvertOSCAL ...
var ConvertOSCAL = cli.Command{
	Name:  "oscal",
	Usage: "convert between one or more OSCAL file formats",
	Description: `Convert between OSCAL-formatted XML and JSON files. The command accepts
   one or more source file paths and can also be used with source file contents
	 piped/redirected from STDIN.`,
	ArgsUsage: "[source-files...]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "output-path, o",
			Usage:       "Output path for converted file(s). Defaults to current working directory",
			Destination: &outputPath,
		},
		cli.StringFlag{
			Name:        "output-file, f",
			Usage:       `File name for converted output from STDIN. Defaults to "stdin.<json|xml|yaml>"`,
			Destination: &outputFile,
		},
		cli.BoolFlag{
			Name:        "yaml",
			Usage:       "If source file format is XML or JSON, also generate equivalent YAML output",
			Destination: &yaml,
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() < 1 {
			// Check for stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				return nil
			}

			return cli.NewExitError("oscalkit convert requires at least one argument", 1)
		}

		if c.NArg() > 1 {
			for _, arg := range c.Args() {
				// Prevent the use of both stdin and specific source files
				if arg == "-" {
					return cli.NewExitError("Cannot use both file path and '-' (STDIN) in args", 1)
				}
			}
		}

		if c.Args().First() != "-" && outputFile != "" {
			return cli.NewExitError("--output-file (-f) is only used when converting from STDIN (-)", 1)
		}

		return nil
	},
	Action: func(c *cli.Context) error {
		// Parse stdin via pipe or redirection
		if c.NArg() <= 0 || c.Args().First() == "-" {
			outputFormat, err := validateStdin(os.Stdin)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Error parsing from STDIN: %s", err), 1)
			}

			destFile, err := os.Create(outputFile)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Error opening output file %s: %s", outputFile, err), 1)
			}
			defer destFile.Close()

			return convert(os.Stdin, destFile, outputFormat)
		}

		// Convert each source file
		for _, sourcePath := range c.Args() {
			// if isValidURL(sourcePath) {
			// 	resp, err := http.Get(sourcePath)
			// 	if err != nil {
			// 		return cli.NewExitError(fmt.Sprintf("Error convert to OSCAL from URL %s: %s", sourcePath, err), 1)
			// 	}
			// 	defer resp.Body.Close()

			// 	rawResp, _ := ioutil.ReadAll(resp.Body)
			// 	logrus.Info(string(rawResp))
			// }

			matches, _ := filepath.Glob(sourcePath)

			for _, match := range matches {
				srcFile, err := os.Open(match)
				if err != nil {
					return err
				}
				defer srcFile.Close()

				destPath, outputFormat := createOutputPath(sourcePath)

				destFile, err := os.Create(destPath)
				if err != nil {
					return err
				}
				defer destFile.Close()

				if err := convert(srcFile, destFile, outputFormat); err != nil {
					return cli.NewExitError(fmt.Sprintf("Error converting to OSCAL from file %s: %s", match, err), 1)
				}
			}
		}

		return nil
	},
}

func validateStdin(stdin *os.File) (string, error) {
	rawSource, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	var xmls string
	var jsons string

	if err = xml.Unmarshal(rawSource, &xmls); err == nil {
		if outputFile == "" {
			outputFile = "stdin.json"
		}

		return "json", nil
	} else if err = json.Unmarshal(rawSource, &jsons); err == nil {
		if outputFile == "" {
			outputFile = "stdin.xml"
		}

		return "xml", nil
	}

	return "", errors.New("File content from STDIN is neither XML nor JSON")

}

// Not yet parsing rawSource arg for STDIN
func convert(src io.Reader, dest io.Writer, outputFormat string) error {
	switch outputFormat {
	case "json":
		logrus.Debug("Converting to JSON")

		oscal, err := oscal.New(src)
		if err != nil {
			return err
		}

		oscalJSON, err := oscal.RawJSON(true)
		if err != nil {
			return err
		}

		if _, err := dest.Write(oscalJSON); err != nil {
			return err
		}

		logrus.Info("Successfully converted to JSON")

		return nil

	case "xml":
		logrus.Debug("Converting to XML")

		oscal, err := oscal.New(src)
		if err != nil {
			return err
		}

		oscalXML, err := oscal.RawXML(true)
		if err != nil {
			return err
		}

		if _, err := dest.Write(oscalXML); err != nil {
			return err
		}

		logrus.Info("Successfully converted to XML")

		return nil

	case "yaml":
		logrus.Debug("Converting to YAML")

		oscal, err := oscal.New(src)
		if err != nil {
			return err
		}

		oscalYAML, err := oscal.RawYAML()
		if err != nil {
			return err
		}

		if _, err := dest.Write(oscalYAML); err != nil {
			return err
		}

		logrus.Info("Successfully converted to YAML")

		return nil
	}

	return fmt.Errorf("Output format %s is not supported", outputFormat)
}

// func isValidURL(urlStr string) bool {
// 	_, err := url.ParseRequestURI(urlStr)
// 	if err != nil {
// 		return false
// 	}

// 	return true
// }

func createOutputPath(srcPath string) (string, string) {
	if srcPath == "" {
		return "", ""
	}

	sourceExt := filepath.Ext(srcPath)[1:]

	var outputFormat string
	if sourceExt == "xml" {
		outputFormat = "json"
	} else {
		outputFormat = "xml"
	}

	filePath := fmt.Sprintf("%s.%s", strings.Split(path.Base(srcPath), ".")[0], outputFormat)

	if outputPath != "" {
		filePath = path.Join(outputPath, filePath)
	}

	return filePath, outputFormat
}
