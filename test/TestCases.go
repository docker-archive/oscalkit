package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/fatih/color"
)

// SecurityControlsSubcontrolCheck is a test to verify that all controls from the catalog are being mapped correctly
func SecurityControlsSubcontrolCheck(check []catalog.Catalog, ProfileFile string) error {

	codeGeneratedControls := ProtocolsMapping(check)

	f, err := os.Open(ProfileFile)
	if err != nil {
		log.Fatal(err)
	}

	parsedProfile, err := GetProfile(f)
	if err != nil {
		log.Fatal(err)
	}

	profileControlsDetails := ProfileProcessing(parsedProfile)

	if len(codeGeneratedControls) == len(profileControlsDetails) {
		println("Perfect Count Match")
		println("Go file control, sub-control count: ", len(codeGeneratedControls))
		println("Profile control, sub-control count: ", len(profileControlsDetails))
		codeGeneratedMapping := ProtocolsMapping(check)
		mapcompareflag := AreMapsSame(profileControlsDetails, codeGeneratedMapping)
		if mapcompareflag {
			color.Green("ID, Class & Title Mapping Correct")
		} else {
			color.Red("ID, Class & Title Mapping Incorrect")
		}
	} else if len(codeGeneratedControls) > len(profileControlsDetails) {
		println("Controls in go file are greater in number then present in profile")
		println("Go file control, sub-control count: ", len(codeGeneratedControls))
		println("Profile control, sub-control count: ", len(profileControlsDetails))
		color.Red("ID, Class & Title Mapping Incorrect")
	} else if len(codeGeneratedControls) < len(profileControlsDetails) {
		println("Controls in profile are greater in number then present in go file")
		println("Go file control, sub-control count: ", len(codeGeneratedControls))
		println("Profile control, sub-control count: ", len(profileControlsDetails))
		color.Red("ID, Class & Title Mapping Incorrect")
	}
	file, _ := filepath.Glob("./oscaltesttmp*")
	if file != nil {
		for _, f := range file {
			os.RemoveAll(f)
		}
	}
	return nil
}
