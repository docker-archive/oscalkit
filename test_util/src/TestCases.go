package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/fatih/color"
)

// SecurityControlsControlCheck is a test to verify that all controls from the catalog are being mapped correctly
func SecurityControlsControlCheck(check []catalog.Catalog, ProfileFile string) error {

	codeGeneratedControls := ProtocolsMapping(check)

	f, err := os.Open(ProfileFile)
	if err != nil {
		log.Fatal(err)
	}

	parsedProfile, err := GetProfile(f)
	if err != nil {
		log.Fatal(err)
	}

	listParentControls := ParentControls(parsedProfile)

	profileControlsDetails := ProfileProcessing(parsedProfile, listParentControls)

	if Count(codeGeneratedControls, "controls") == Count(profileControlsDetails, "controls") {
		color.Green("Controls & SubControls Count Matched")
		println("Go file control & sub-control count: ", Count(codeGeneratedControls, "controls"))
		println("Profile control & sub-control count: ", Count(profileControlsDetails, "controls"))
	} else if Count(codeGeneratedControls, "controls") > Count(profileControlsDetails, "controls") {
		color.Red("Controls & Controls in go file are greater in number then present in profile")
		println("Go file control & sub-control count: ", Count(codeGeneratedControls, "controls"))
		println("Profile control & sub-control count: ", Count(profileControlsDetails, "controls"))
	} else if Count(codeGeneratedControls, "controls") < Count(profileControlsDetails, "controls") {
		color.Red("Controls & Controls in profile are greater in number then present in go file")
		println("Go file control & sub-control count: ", Count(codeGeneratedControls, "controls"))
		println("Profile control & sub-control count: ", Count(profileControlsDetails, "controls"))
	}

	controlMapCompareFlag := AreMapsSame(profileControlsDetails, codeGeneratedControls, "controls")
	if controlMapCompareFlag {
		color.Green("ID, Class & Title Mapping Of All Controls & SubControls Correct")
	} else {
		color.Red("ID, Class & Title Mapping Of All Controls & SubControls Incorrect")
	}

	if Count(codeGeneratedControls, "parts") == Count(profileControlsDetails, "parts") {
		color.Green("Parts Count Matched")
		println("Go file parts count: ", Count(codeGeneratedControls, "parts"))
		println("Profile parts count: ", Count(profileControlsDetails, "parts"))
	} else if Count(codeGeneratedControls, "parts") > Count(profileControlsDetails, "parts") {
		color.Red("Parts in go file are greater in number then present in profile")
		println("Go file parts count: ", Count(codeGeneratedControls, "parts"))
		println("Profile parts count: ", Count(profileControlsDetails, "parts"))
	} else if Count(codeGeneratedControls, "parts") < Count(profileControlsDetails, "parts") {
		color.Red("Parts in profile are greater in number then present in go file")
		println("Go file parts count: ", Count(codeGeneratedControls, "parts"))
		println("Profile parts count: ", Count(profileControlsDetails, "parts"))
	}

	partsMapCompareFlag := AreMapsSame(profileControlsDetails, codeGeneratedControls, "parts")
	if partsMapCompareFlag {
		color.Green("ID, Class & Title Mapping Of Parts Correct")
	} else {
		color.Red("ID, Class & Title Mapping Of All Parts Incorrect")
	}

	file, _ := filepath.Glob("./oscaltesttmp*")
	if file != nil {
		for _, f := range file {
			os.RemoveAll(f)
		}
	}
	return nil
}
