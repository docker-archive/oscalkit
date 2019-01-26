package convert

var includeXML bool

// ConvertOpenControl ...
// var ConvertOpenControl = cli.Command{
// 	Name:  "opencontrol",
// 	Usage: `convert from OpenControl format to OSCAL "implementation" format`,
// 	Description: `Convert OpenControl-formatted "component" and "opencontrol" YAML into
// 	 OSCAL-formatted "implementation" layer JSON`,
// 	ArgsUsage: "[opencontrol.yaml-filepath] [opencontrols-dir-path]",
// 	Flags: []cli.Flag{
// 		cli.BoolFlag{
// 			Name:        "yaml, y",
// 			Usage:       "Generate YAML in addition to JSON",
// 			Destination: &yaml,
// 		},
// 		cli.BoolFlag{
// 			Name:        "xml, x",
// 			Usage:       "Generate XML in addition to JSON",
// 			Destination: &includeXML,
// 		},
// 	},
// 	Before: func(c *cli.Context) error {
// 		if c.NArg() != 2 {
// 			return cli.NewExitError("Missing opencontrol.yaml file and path to opencontrols/ directory", 1)
// 		}

// 		return nil
// 	},
// 	Action: func(c *cli.Context) error {
// 		ocOSCAL, err := oscal.NewFromOC(oscal.OpenControlOptions{
// 			OpenControlYAMLFilepath: c.Args().First(),
// 			OpenControlsDir:         c.Args()[1],
// 		})
// 		if err != nil {
// 			return cli.NewExitError(err, 1)
// 		}

// 		if includeXML {
// 			rawXMLOCOSCAL, err := ocOSCAL.XML(true)
// 			if err != nil {
// 				return cli.NewExitError(fmt.Sprintf("Error producing raw XML: %s", err), 1)
// 			}
// 			if err := ioutil.WriteFile("opencontrol-oscal.xml", rawXMLOCOSCAL, 0644); err != nil {
// 				return cli.NewExitError(err, 1)
// 			}
// 		}

// 		if yaml {
// 			rawYAMLOCOSCAL, err := ocOSCAL.YAML()
// 			if err != nil {
// 				return cli.NewExitError(err, 1)
// 			}
// 			if err := ioutil.WriteFile("opencontrol-oscal.yaml", rawYAMLOCOSCAL, 0644); err != nil {
// 				return cli.NewExitError(err, 1)
// 			}
// 		}

// 		rawOCOSCAL, err := ocOSCAL.JSON(true)
// 		if err != nil {
// 			return cli.NewExitError(err, 1)
// 		}

// 		if err := ioutil.WriteFile("opencontrol-oscal.json", rawOCOSCAL, 0644); err != nil {
// 			return cli.NewExitError(err, 1)
// 		}

// 		return nil
// 	},
// }
