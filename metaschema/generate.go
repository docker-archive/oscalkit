// +build ignore

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/docker/oscalkit/metaschema"
)

const (
	oscalRepo         = "https://github.com/usnistgov/OSCAL.git"
	metaschemaBaseDir = "OSCAL/schema/metaschema/%s"
)

var (
	pkgName = map[string]string{
		"catalog":              "catalog",
		"profile":              "profile",
		"system-security-plan": "ssp",
	}
)

func main() {
	rmCmd := exec.Command("rm", "-rf", "OSCAL/")
	if err := rmCmd.Run(); err != nil {
		log.Fatal(err)
	}

	cloneCmd := exec.Command("git", "clone", oscalRepo)
	if err := cloneCmd.Run(); err != nil {
		log.Fatal(err)
	}

	metaschemaPaths := map[string]string{
		"catalog": "oscal-catalog-metaschema.xml",
		"profile": "oscal-profile-metaschema.xml",
		"ssp":     "oscal-ssp-metaschema.xml",
	}

	for pkg, metaschemaPath := range metaschemaPaths {
		f, err := os.Open(fmt.Sprintf(metaschemaBaseDir, metaschemaPath))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		meta, err := decode(f)
		if err != nil {
			log.Fatal(err)
		}
		meta.Root = pkg

		if err := metaschema.GenerateTypes(meta); err != nil {
			log.Fatalf("Error generating go types for metaschema: %s", err)
		}
	}

	rmCmd = exec.Command("rm", "-rf", "OSCAL/")
	if err := rmCmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func decode(r io.Reader) (*metaschema.Metaschema, error) {
	var meta metaschema.Metaschema

	d := xml.NewDecoder(r)

	if err := d.Decode(&meta); err != nil {
		return nil, fmt.Errorf("Error decoding metaschema: %s", err)
	}

	if meta.Import != nil && meta.Import.Href != nil && meta.Import.Href.URL != nil {
		imf, err := os.Open(fmt.Sprintf(metaschemaBaseDir, meta.Import.Href.URL.String()))
		if err != nil {
			return nil, err
		}
		defer imf.Close()

		importedMeta, err := decode(imf)
		if err != nil {
			return nil, err
		}

		meta.ImportedMetaschema = importedMeta
	}

	return &meta, nil
}
