package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/opencontrol/oscalkit"
)

func main() {

	var check = oscalkit.ApplicableControls
	profile := flag.String("p", "", "Path of the profile")
	flag.Parse()
	file, _ := filepath.Glob("./oscaltesttmp*")
	if file != nil {
		for _, f := range file {
			os.RemoveAll(f)
		}
	}
	SecurityControlsSubcontrolCheck(check, *profile)
}
